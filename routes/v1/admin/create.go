package admin

import (
	"net/mail"

	"github.com/gin-gonic/gin"
	"github.com/manthan307/corebase/db"
	"github.com/manthan307/corebase/db/schema"
	"github.com/manthan307/corebase/utils/crypto"
	"go.uber.org/zap"
)

type AdminCreateInputBody struct {
	Username string `json:"username" binding:"required,min=3,max=32,alphanum"`
	Password string `json:"password" binding:"required,min=8"`
	Email    string `json:"email" binding:"required,email"`
}

func CreateAdmin(log *zap.Logger, client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			c.JSON(400, gin.H{
				"message": "token is required",
			})
			return
		}

		role, err := ValidateToken(c.Request.Context(), token, client)
		if err != nil {
			c.JSON(400, gin.H{
				"message": err.Error(),
			})
			return
		}

		var body AdminCreateInputBody
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(400, gin.H{
				"message": err.Error(),
			})
			return
		}

		if _, err := mail.ParseAddress(body.Email); err != nil {
			c.JSON(400, gin.H{
				"message": "Invalid email address",
			})
			return
		}

		isExists, WhichOneExists, err := client.Admin.Exists(c.Request.Context(), body.Email, body.Username)
		if err != nil {
			c.JSON(500, gin.H{
				"message": err.Error(),
			})
			return
		}
		if isExists {
			c.JSON(409, gin.H{
				"message": WhichOneExists + " already exists",
			})
			return
		}

		hash := crypto.HashPassword(body.Password)

		admin, err := client.Admin.Create(c.Request.Context(), schema.AdminParams{
			Username: body.Username,
			Password: hash,
			Email:    body.Email,
			Role:     role,
		})
		if err != nil {
			c.JSON(500, gin.H{
				"message": err.Error(),
			})
			return
		}

		if err := DeleteToken(c.Request.Context(), token, client); err != nil {
			log.Error("Error deleting token", zap.Error(err))
			c.JSON(500, gin.H{
				"message": "Error deleting token",
			})
			return
		}

		// c.SetCookie()

		c.JSON(200, gin.H{
			"id":        admin.ID,
			"username":  admin.Username,
			"email":     admin.Email,
			"role":      admin.Role,
			"createdAt": admin.CreatedAt,
			"updatedAt": admin.UpdatedAt,
		})

	}
}

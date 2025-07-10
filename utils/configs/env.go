package configs

import "github.com/joho/godotenv"

func LoadEnv(path string) {
	err := godotenv.Load(path)
	if err != nil {
		panic(err)
	}
}

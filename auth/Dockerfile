# Base image
FROM node:20-alpine

# Create app directory
WORKDIR /app

# Copy package files separately for caching
COPY package*.json ./

# Install dependencies
RUN npm install --production

# Copy the rest of the code
COPY . .

# Expose the port your app runs on
EXPOSE 3000

# Start the app
CMD ["node", "index.js"]

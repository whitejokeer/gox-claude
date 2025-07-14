module blog-app/app

go 1.21

require (
	blog-app v0.0.0-00010101000000-000000000000
	github.com/joho/godotenv v1.5.1
	github.com/gin-gonic/gin v1.9.1
	github.com/golang-jwt/jwt/v5 v5.2.0
)

replace blog-app => ../

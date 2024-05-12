package middleware

import (
	"github.com/DaffaJatmiko/go-task-manager/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Auth() gin.HandlerFunc {
	return gin.HandlerFunc(func(ctx *gin.Context) {
		// TODO: answer here
		sessionToken, err := ctx.Cookie("session_token")
		if err != nil {
			if ctx.Request.Header.Get("Content-Type") == "application/json" {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				ctx.Abort()
				return 
			}
			ctx.Redirect(http.StatusSeeOther, "/login")
			return 
		}

		claims := &model.Claims{}
		token, err := jwt.ParseWithClaims(sessionToken, claims, func(token *jwt.Token) (interface{}, error){
			return model.JwtKey, nil
		})

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}

		if !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		ctx.Set("email", claims.Email)
		ctx.Next()
	})
}

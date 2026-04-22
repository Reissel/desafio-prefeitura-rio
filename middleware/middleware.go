package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// TODO: Alterar secret do JWT
var jwtSecret = []byte("abobrinhaabobrinhaabobrinhaabobrinhaabobrinha")

// TODO: Alterar secret do header X-Signature-256
const SecretKey = "abobrinha"

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header must be Bearer token"})
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			cpf, found := claims["preferred_username"].(string)
			if !found || cpf == "" {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "preferred_username claim missing"})
				return
			}

			c.Set("cpf", cpf)
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		c.Next()
	}
}

func SignatureMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		headerSignature := c.GetHeader("X-Signature-256")
		if headerSignature == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing signature"})
			return
		}

		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Could not read body"})
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		h := hmac.New(sha256.New, []byte(SecretKey))
		h.Write(bodyBytes)
		expectedSignature := hex.EncodeToString(h.Sum(nil))

		if !hmac.Equal([]byte(headerSignature), []byte(expectedSignature)) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
			return
		}

		c.Next()
	}
}

package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

//Generate JWT
func GenerateJWT(userID string) (string, error){
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})

	secret := os.Getenv("JWT_SECRET")

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

//Verify JWT
func VerifyJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil

}

//Auth Middleware
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}
		// "Bearer eyjjbcjbdcjbdkj.cbkjBCJKBCKJEBFF.6G7347ED87GHE8DFU"

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		_, err := VerifyJWT(tokenString)
		if err != nil {
			http.Error(w, "invalid token", http.StatusInternalServerError)
			return
		}
		next.ServeHTTP(w, r)

	})
}

//Get UserID from Token
func GetUserIDFromToken(r *http.Request) (uint, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0, fmt.Errorf("authorization header required")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := VerifyJWT(tokenString)
	if err != nil {
		return 0, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("could not parse claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return 0, fmt.Errorf("user id not found in token")
	}

	userIDInt, _ := strconv.Atoi(userID)

	return uint(userIDInt), nil
}


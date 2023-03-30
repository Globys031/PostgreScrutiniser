package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Globys031/PostgreScrutiniser/backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JwtWrapper struct {
	SecretKey       string
	Issuer          string
	ExpirationHours int64
}

type JwtClaims struct {
	jwt.RegisteredClaims
	Name string // postgrescrutiniser main system user name
}

// Generate JSON WEB TOKEN that will be saved in client's local storage.
func (wrapper *JwtWrapper) GenerateToken(user LoginRequest, logger *utils.Logger) (signedToken string, err error) {
	claims := &JwtClaims{
		Name: *user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Hour * time.Duration(wrapper.ExpirationHours))),
			Issuer:    wrapper.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err = token.SignedString([]byte(wrapper.SecretKey))
	if err != nil {
		err = fmt.Errorf("failed to generate token: %v", err)
		logger.LogError(err)
		return "", err
	}

	return signedToken, nil
}

func (wrapper *JwtWrapper) ValidateToken(signedToken string, logger *utils.Logger) (claims *JwtClaims, err error) {
	// 1. Validate token
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JwtClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(wrapper.SecretKey), nil
		},
	)
	if err != nil {
		logger.LogError(fmt.Errorf("failed to validate token: %v", err))
		return
	}

	claims, ok := token.Claims.(*JwtClaims)
	if !ok {
		logger.LogError(fmt.Errorf("failed to validate token: %v", err))
		return
	}

	// 2. If valid, check if token has expired
	if claims.ExpiresAt.Time.Before(time.Now().Local()) {
		return nil, fmt.Errorf("JWT is expired")
	}

	return claims, nil
}

// Used as a middleware function for anything that requires authentification
func (wrapper *JwtWrapper) ValidateTokenMiddleware(logger *utils.Logger) MiddlewareFunc {
	return func(c *gin.Context) {
		// Get the token from the Authorization header
		var token string
		if authorizationHeader := c.GetHeader("Authorization"); authorizationHeader != "" {
			token = strings.Split(authorizationHeader, "Bearer ")[1]
		}

		// Validate the token using the JwtWrapper.ValidateToken function
		claims, err := wrapper.ValidateToken(token, logger)
		if err != nil {
			errorMsg := &ErrorMessage{
				ErrorMessage: err.Error(),
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorMsg)
			return
		}

		// If the token is valid, set the claims in the request context and continue
		// c.Set("jwtClaims", claims)
		c.Set("bearerAuth.Scopes", claims)
		c.Next()
	}
}

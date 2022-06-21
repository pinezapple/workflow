// Package middlewares provides methods to preprocessing the requese, before pass it to controllers
package middlewares

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/uc-cdis/go-authutils/authutils"
	"github.com/vfluxus/heimdall/core"
)

var (
	jwtApp   *authutils.JWTApplication
	expected *authutils.Expected
	logger   *core.Logger = core.GetLogger()
)

// InitAuthN init JWT validator
func InitAuthN(config *core.ServerConfig) {
	jwtApp = authutils.NewJWTApplication(config.JWKSUrl)
	expected = &authutils.Expected{
		Audiences: config.Audiences,
		Issuers:   config.Issuers,
		Purpose:   &config.Purpose,
	}
}

// ValidateToken validate authentication token
func ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		message, err := validateToken(c)
		if err != nil {
			c.Header("WWW-Authenticate", message)
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}

// DecodeToken extracts user info in JWT token
func DecodeToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		config := core.GetConfig().Server
		if config.Environment == "production" {
			message, err := validateToken(c)
			if err != nil {
				c.Header("WWW-Authenticate", message)
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
		}

		authHeader := c.GetHeader("Authorization")
		if len(authHeader) < 6 {
			logger.Errorf("Invalid Authorization length", authHeader)
			c.Header("WWW-Authenticate", "token error")
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		jwtToken := authHeader[len("Bearer"):]
		claims, err := jwtApp.Decode(jwtToken)
		if err != nil {
			logger.Errorf("Decode jwt token error: ", err)
			c.Header("WWW-Authenticate", "token error")
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		cl := *claims

		claimContext, ok := cl["context"]
		if !ok {
			c.AbortWithStatus(http.StatusForbidden)
			return
		} else {
			mapClaim := claimContext.(map[string]interface{})
			user, ok := mapClaim["user"]
			if ok {
				name := user.(map[string]interface{})
				username, ok := name["name"]
				if ok {
					c.Set("UserName", username)
					c.Next()
					return
				}
			}

			c.AbortWithStatus(http.StatusForbidden)
		}
	}
}

func validateToken(c *gin.Context) (string, error) {
	realm := "Authentication required"
	realm = "Authen realm=" + strconv.Quote(realm)

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		logger.Errorf("Authorization header: ", authHeader)
		return realm, errors.New("Unauthorized")
	}

	jwtToken := authHeader[len("Bearer"):]
	claims, err := jwtApp.Decode(jwtToken)
	if err != nil {
		logger.Errorf("JWT token: %s error %v", authHeader, err.Error())
		return realm, errors.New("Unauthorized")
	}

	err = expected.Validate(claims)
	if err != nil {
		logger.Errorf("Validate token result: %v", err.Error())
		return realm, errors.New("Unauthorized")
	}

	return "", nil
}

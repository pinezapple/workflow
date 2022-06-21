package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vfluxus/valkyrie/core"
	utilsModel "github.com/vfluxus/workflow-utils/model"
)

func DecodeJWT(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		access_token := c.Request().Header.Get("Authorization")

		ok, userid, err := JWTCheck(access_token)
		if !ok || err != nil {
			return c.JSON(http.StatusUnauthorized, &utilsModel.Response{
				Data: nil,
				Error: utilsModel.ResponseError{
					Message: err.Error(),
					Code:    http.StatusNotAcceptable,
				},
			})
		}

		ctx := c.Request().Context()
		c.SetRequest(c.Request().Clone(context.WithValue(ctx, "UserID", userid)))

		return next(c)
	}
}

func JWTCheck(token string) (ok bool, userid string, err error) {
	if token == "" {
		return false, "", errors.New("Empty access_token")
	}
	jwtApp := core.GetJWTApplication()
	e := core.GetExpectJWT()

	token = token[len("Bearer"):]

	claims, err := jwtApp.Decode(token)
	if err != nil {
		return false, "", err
	}

	config := core.GetMainConfig()
	if config.HttpServerConfig.ENV == "production" {
		err = e.Validate(claims)
		if err != nil {
			return false, "", err
		}
	}

	cl := *claims
	claimCtx, ok := cl["context"]
	if !ok {
		return false, "", fmt.Errorf("no user context")
	} else {
		mapClaim := claimCtx.(map[string]interface{})
		b, ok := mapClaim["user"]
		if !ok {
			return false, "", fmt.Errorf("no user context")
		} else {
			c := b.(map[string]interface{})
			d, ok := c["name"]
			if !ok {
				return false, "", fmt.Errorf("no user context")
			} else {
				return true, d.(string), nil
			}
		}
	}
}

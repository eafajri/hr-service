package transport

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/eafajri/hr-service.git/module/employee/internal/entity"
	"github.com/eafajri/hr-service.git/module/employee/internal/usecase"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func BasicAuthMiddleware(userUc usecase.UserUseCase) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Basic ") {
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing or invalid Authorization header")
			}

			payload, err := base64.StdEncoding.DecodeString(authHeader[len("Basic "):])
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid base64 encoding")
			}

			creds := strings.SplitN(string(payload), ":", 2)
			if len(creds) != 2 {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid basic auth format")
			}
			username := creds[0]
			password := creds[1]

			user, err := userUc.GetUserByUsernaname(username)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid username or password")
			}

			if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid username or password")
			}

			requestID := c.Request().Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = uuid.New().String()
			}

			userContext := entity.UserContext{
				UserID:    user.ID,
				Username:  user.Username,
				Role:      user.Role,
				IPAddress: c.RealIP(),
				RequestID: requestID,
			}

			c.Set("user_context", userContext)

			return next(c)
		}
	}
}

func AdminPrevilageMiddleware(userUc usecase.UserUseCase) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, ok := c.Get("user_context").(entity.UserContext)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "User ID not found in context")
			}

			if user.Role != "admin" {
				return echo.NewHTTPError(http.StatusForbidden, "Access denied: admin privileges required")
			}

			return next(c)
		}
	}
}

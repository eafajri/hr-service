package rest

import (
	"net/http"
	"time"

	moduleConfig "github.com/eafajri/hr-service.git/module/employee/config"
	"github.com/labstack/echo/v4"
)

type Rest struct {
}

func StartRest(echoInstance *echo.Echo) {

	_ = moduleConfig.NewModuleDependencies()

	restHandler := &Rest{}

	publicApi := echoInstance.Group("/public")
	publicApi.GET("/check", restHandler.CheckHealth)
}

func (h *Rest) CheckHealth(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ok", "timestamp": time.Now().Format(time.RFC3339)})
}

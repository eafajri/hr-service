package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/eafajri/hr-service.git/config"
	employeeRest "github.com/eafajri/hr-service.git/module/employee/transport/rest"
	"github.com/labstack/echo/v4"
)

func main() {
	conf := config.GetConfig()
	e := echo.New()

	employeeRest.StartRest(e)

	// Start server in goroutine
	go func() {
		addr := fmt.Sprintf(":%v", conf.ServerRestPort)
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("shutting down the server: %s", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down the server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %s", err)
	}
}

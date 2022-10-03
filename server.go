package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/microctar/licorice/app/route"
)

var router *echo.Echo

func init() {

	router = echo.New()
	router.Logger.Debug()
	// enable gzip support
	router.Use(middleware.Gzip())
	router.Use(middleware.Logger())
	// restful api
	router.GET("/clash/:link", route.ExportClashConfig)
	router.GET("/clash/:link/:rulefile", route.ExportClashConfig)

}

func RunServer() {

	go func() {
		if err := router.Start(fmt.Sprintf(":%d", port)); err != nil && err != http.ErrServerClosed {
			router.Logger.Fatal("Shutting down the server")
		}
	}()

	// graceful shutdown
	safetybolt := make(chan os.Signal, 1)

	signal.Notify(safetybolt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	<-safetybolt

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	if err := router.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

}

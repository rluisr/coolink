package main

import (
	"context"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const ginContextTimeout = 10

func main() {
	router := gin.New()
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Todo env
		AllowedMethods:   []string{"GET"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"*"},
		AllowCredentials: true,
	}))
	router.Use(setHeader())

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hi")
	})
	router.GET("/image", getImage)

	srv := &http.Server{
		Addr:    ":8080", // Todo env
		Handler: router,
	}

	go func() {
		log.Printf("listen: %s\n", srv.Addr)

		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), ginContextTimeout*time.Second)
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	defer cancel()
	log.Println("Server exiting")
}

func setHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("x-content-type-options", "nosniff")
		c.Writer.Header().Set("x-xss-protection", "1; mode=block")
		c.Writer.Header().Set("x-frame-options", "DENY")
		c.Writer.Header().Set("cache-control", "public, max-age=86400")
	}
}

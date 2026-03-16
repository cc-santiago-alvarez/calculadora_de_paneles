package main

import (
	"context"
	"embed"
	"log"
	"time"

	"github.com/dev13/calculadora-paneles-backend/internal/app"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	application, err := app.New()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	defer application.Close()

	// Start the HTTP API server in background
	server := application.StartHTTPServer()

	// Create and run Wails desktop app
	err = wails.Run(&options.App{
		Title:  "Calculadora de Paneles Solares",
		Width:  1400,
		Height: 900,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: func(ctx context.Context) {
			// HTTP server is already running
		},
		OnShutdown: func(ctx context.Context) {
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if err := server.Shutdown(shutdownCtx); err != nil {
				log.Printf("Server forced to shutdown: %v", err)
			}
			log.Println("Server stopped")
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}

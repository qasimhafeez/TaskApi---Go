package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/qasimhafeez/taskapi/internal/config"
	"github.com/qasimhafeez/taskapi/internal/httpx"
	"github.com/qasimhafeez/taskapi/internal/logger"
	"github.com/qasimhafeez/taskapi/internal/repo"
	"github.com/qasimhafeez/taskapi/internal/service"
)

	func main() {
		// Load configs and logger
		configs := config.FromEnv()
		logg := logger.New()

		// Build Dependencies
		store := repo.NewInMem()
		service := service.NewTaskService(store)
		httpAdapter := httpx.NewHandler(service, logg)

		// Create HTTP Server for timeout safety
		server:= &http.Server{
			Addr: configs.Addr,
			Handler: httpAdapter.Router(),
			ReadHeaderTimeout: 5*time.Second,
		}

		// Running server in goroutine
		go func(){
			logg.Info("Server Starting", slog.String("addr", configs.Addr))
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				// Fail Fast if server cannot start
				log.Fatal(err)
			}
		}()

		// Graceful shutdown
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		<-stop

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		logg.Info("Server Shutting Down")
		_ = server.Shutdown(ctx)

	}
	
package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/iamneuron/students-check-api/internal/config"
	"github.com/iamneuron/students-check-api/internal/http/handlers/student"
)

func main() {

	//load config

	cfg := config.MustLoad()

	//datase setup

	//setup router
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New())

	//greacefull shutdown

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	//setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}
	fmt.Printf("server started  %s", cfg.HttpServer.Addr)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("failed to start server")
		}
	}()

	<-done
	slog.Info("shuting down the server.")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := server.Shutdown(ctx)
	if err != nil {
		slog.Error("failed to shutdoen the server", slog.String("error", err.Error()))
	}
	slog.Info("server shoutdown successfuely")

}

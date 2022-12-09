package main

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
	"github.com/toannhu96/mspm_service/pkg/api"
	"github.com/toannhu96/mspm_service/pkg/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{TimestampFormat: "2006-01-02 15:04:05", FullTimestamp: true})
	logrus.SetLevel(logrus.InfoLevel)

	// singleton init
	keywordService, err := service.NewKeywordService()
	if err != nil {
		logrus.Errorf("fail to init keyword service, err=%v", err)
		return
	}
	baseHandler := api.NewBaseHandler()
	keywordHandler := api.NewKeywordHandler(baseHandler, keywordService)

	r := chi.NewRouter()
	// RequestID is a middleware that injects a request ID into the context of each
	// request. A request ID is a string of the form "host.example.com/random-0001",
	// where "random" is a base62 random string that uniquely identifies this go
	// process, and where the last number is an atomically incremented request
	// counter.
	r.Use(middleware.RequestID)
	// Recoverer is a middleware that recovers from panics, logs the panic (and a
	// backtrace), and returns a HTTP 500 (Internal Server Error) status if
	// possible. Recoverer prints a request ID if one is provided.
	r.Use(middleware.Recoverer)
	// Logger middleware
	r.Use(middleware.Logger)

	r.Mount("/", keywordHandler.Route())

	// The HTTP Server
	server := &http.Server{Addr: "0.0.0.0:8088", Handler: r}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				logrus.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			logrus.Fatal(err)
		}
		serverStopCtx()
	}()

	// Run the server
	logrus.Info("Server starting at port 8088...")
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logrus.Fatal(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}

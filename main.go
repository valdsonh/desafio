package main

import (
	"context"
	"fmt"
	"github/desafio/routes"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	port = ":8080"
)

func main() {
	router := mux.NewRouter()
	routes.InitializeRoutes(router)

	server := &http.Server{
		Handler: router,
		Addr:    port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Executa o servidor http
	go func() {
		fmt.Printf("Listening and serving on port %s\n", port)
		logrus.Info("API routes")
		router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, err := route.GetPathTemplate()
		if err != nil {
			return nil
		}
		methods, err := route.GetMethods()
		if err != nil {
			return nil
		}
		logrus.Info(fmt.Sprintf("%v %s", methods, path))
		return nil
	})
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalln(fmt.Sprintf("listen: %s\n", err))
		}
	}()

	// greacefull shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown
	logrus.Info("starting shutdown server")
	defer logrus.Info("finishing shutdown server")

	const timeout = 5
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logrus.Error("Could not shutdown server gracefully")

	}
	<-ctx.Done()
}

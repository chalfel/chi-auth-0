package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chalfel/chi-auth-0/internal/api"
	"github.com/chalfel/chi-auth-0/pkg/config"
	"github.com/chalfel/chi-auth-0/pkg/db"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/spf13/cobra"
)

func NewApiCmd() *cobra.Command {
	apiCmd := &cobra.Command{
		Use:   "api",
		Short: "Chi auth 0 api",
	}

	apiCmdServe := &cobra.Command{
		Use:   "serve",
		Short: "serve chi auth0 api",
		Run: func(cmd *cobra.Command, args []string) {
			serve()
		},
	}

	apiCmd.AddCommand(apiCmdServe)

	return apiCmd
}

type ServerConfig struct {
	Port string `env:"PORT"`
}

func serve() {
	cfg := &ServerConfig{}

	if err := config.LoadConfigFromEnv(cfg); err != nil {
		panic(err)
	}
	// execute api
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
	}))
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	apiInstance := api.NewApi(db.NewDb(), r)

	apiInstance.RegisterRoutes()

	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%s", cfg.Port),
		Handler: r,
	}

	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}

		serverStopCtx()
	}()

	fmt.Printf("Server is listening in port: %s \n", cfg.Port)
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	<-serverCtx.Done()

}

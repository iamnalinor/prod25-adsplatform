package main

import (
	"backend/config"
	"backend/internal/handler"
	"backend/internal/repo"
	"backend/internal/service"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func getDatabase(env config.Environment) (*sqlx.DB, error) {
	return sqlx.Connect("postgres", env.BuildDsn())
}

func main() {
	env := config.LoadEnvironment()
	if env.RunningInCI {
		log.Println("Detected CI environment. Requests logging and AI service will be disabled")
	}

	db, err := getDatabase(env)
	if err != nil {
		log.Fatalf("connect to database: %s\n", err)
	}

	repos := repo.NewRepositories(db)
	services, err := service.NewServices(repos, env)
	if err != nil {
		log.Fatalf("create services: %s\n", err)
	}
	h := handler.NewHandler(services)
	router := h.GetRouter(env)

	go func() {
		err := router.Run(env.ServerAddress)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("run server: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	sig := <-quit

	fmt.Printf("Received signal: %s\n", sig)
}

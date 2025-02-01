package main

import (
	"context"
	"net/http"

	"github.com/bagasadiii/maxcloud_vps/config"
	"github.com/bagasadiii/maxcloud_vps/handler"
	"github.com/bagasadiii/maxcloud_vps/repository"
	"github.com/bagasadiii/maxcloud_vps/service"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	database := config.InitDB()
	logger := config.NewLogger()

	clientRepo := repository.NewClientRepo(database, logger)
	clientService := service.NewClientService(clientRepo, logger)
	clientHandler := handler.NewClientHandler(clientService, logger)

	txSchedulerRepo := repository.NewTransactionSchedulerRepo(database, logger)
	txSchedulerService := service.NewTransactionSchedulerService(database, txSchedulerRepo, logger)

	r := mux.NewRouter()

	r.HandleFunc("/api/register", clientHandler.CreateClient).Methods("POST")
	r.HandleFunc("/api/client/{client_id}", clientHandler.GetClientInfo).Methods("GET")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go txSchedulerService.SchedulerWorkerService(ctx, 5)

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	server.ListenAndServe()
}

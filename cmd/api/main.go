package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/diegobermudez03/couples-backend/internal/config"
	"github.com/diegobermudez03/couples-backend/internal/http/api"
	"github.com/joho/godotenv"
)

//it reads ENV variables
//it runs services that are critical for the applications (migrations)
//it creates the HTTP server
func main(){
	if err := godotenv.Load(".env"); err != nil{
		log.Fatal(err)
	}

	configuration := config.NewConfig()

	//start services 

	//create API server
	//	GRACEFUL SHITDOWN
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	
	server := api.NewAPIServer(configuration)
	go func(){
		log.Printf("Server running on port %s", configuration.Port)
		if err := server.Run(); err != nil && err != http.ErrServerClosed{
			log.Fatalf("couldn't start server %s", err.Error())
		}
	}()

	//listen to cancel signals
	<-ctx.Done()
	log.Println("Interruption signal")
	if err := server.Shutdown(); err != nil{
		log.Fatalf("Server shutdown error %s", err.Error())
	}
	log.Println("Succesfully graceful shutdown")
}
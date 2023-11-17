package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/thejixer/shop-api/database"
	"github.com/thejixer/shop-api/handlers"
	"github.com/thejixer/shop-api/mailer"
	"github.com/thejixer/shop-api/redis"
	"github.com/thejixer/shop-api/server"
)

func init() {
	godotenv.Load()
	env := flag.String("env", "DEV", "enviroment")
	flag.Parse()

	os.Setenv("ENVIROMENT", *env)
	fmt.Println("##########################3")
	fmt.Println("enviroment is : ", os.Getenv("ENVIROMENT"))
	fmt.Println("##########################3")
}

func main() {

	listenAddr := os.Getenv("LISTEN_ADDR")

	store, err := database.NewPostgresStore()

	if err != nil {
		log.Fatal("could not connect to the database: ", err)
	}

	if err := store.Init(); err != nil {
		log.Fatal("could not connect to the database: ", err)
	}

	redisStore := redis.NewRedisStore()
	mailerService := mailer.NewMailerService()

	handlerService := handlers.NewHandlerService(store, redisStore, mailerService)

	apiServer := server.NewAPIServer(listenAddr, handlerService)
	apiServer.Run()

}

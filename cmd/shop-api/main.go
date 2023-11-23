package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/thejixer/shop-api/internal/database"
	"github.com/thejixer/shop-api/internal/handlers"
	"github.com/thejixer/shop-api/internal/mailer"
	"github.com/thejixer/shop-api/internal/redis"
	"github.com/thejixer/shop-api/internal/server"
)

func init() {
	godotenv.Load()
}

func seedDB(store *database.PostgresStore) {
	fmt.Println("db is being seeded")
	store.UserRepo.Create(
		"main addmin",
		os.Getenv("MAIN_ADMIN_EMAIL"),
		os.Getenv("MAIN_ADMIN_PASSWORD"),
		"admin",
		true,
	)
}

func main() {

	env := flag.String("env", "DEV", "enviroment")
	seed := flag.Bool("seed", false, "seed the db")

	flag.Parse()

	os.Setenv("ENVIROMENT", *env)
	fmt.Println("##########################3")
	fmt.Println("enviroment is : ", os.Getenv("ENVIROMENT"))
	fmt.Println("##########################3")

	listenAddr := os.Getenv("LISTEN_ADDR")

	store, err := database.NewPostgresStore()

	if err != nil {
		log.Fatal("could not connect to the database: ", err)
	}

	if err := store.Init(); err != nil {
		log.Fatal("could not connect to the database: ", err)
	}

	if *seed {
		seedDB(store)
	}

	redisStore, err := redis.NewRedisStore()
	if err != nil {
		log.Fatal("could not connect to the redis: ", err)
	}
	mailerService := mailer.NewMailerService()

	handlerService := handlers.NewHandlerService(store, redisStore, mailerService)

	apiServer := server.NewAPIServer(listenAddr, handlerService)
	apiServer.Run()

}

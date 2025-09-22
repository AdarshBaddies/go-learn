package main

import (
	"log"
	"time"

	"github.com/adarshbaddies/go-learn/account"
	"github.com/kelseyhightower/envconfig"
)

type Config struct{
	DatabaseURL string `envconfig:"DATABASE_URL"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	var r account.Repository
	// retry.ForeverSleep(2*time.Second, func(_ int)(err error){
	// 	r, err := account.NewPostgresRepository(cfg.DatabaseURL)
	// 	if err !=nil {
	// 		log.Println(err)
	// 	}
	// 	return
	// })
	// defer r.Close()
	// log.Println("Listening on port 8080...")
	// s := account.NewService(r)
	// log.Fatal(account.ListenGRPC(s, 8080))

	for {
		log.Println("Attempting to connect to the database...")
		r, err = account.NewPostgresRepository(cfg.DatabaseURL)
		if err == nil {
			log.Println("Database connection successful.")
			break
		}

		log.Printf("Failed to connect to the database: %v. Retrying in 2 seconds...\n", err)
		time.Sleep(2 * time.Second)
	}

	defer r.Close()

	log.Println("Listening on port 8080...")
	s := account.NewService(r)
	log.Fatal(account.ListenGRPC(s, 8080))
}
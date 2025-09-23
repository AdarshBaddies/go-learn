package main

import (
	"log"
	"time"

	"github.com/tinrab/retry"
	"github.com/kelseyhightower/envconfig"
	"github.com/adarshbaddies/go-learn/catalog"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	var r catalog.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error){
		r, err = catalog.NewElasticRepository(cfg.DatabaseURL)
		if err != nil {
			log.Println(err)
		}
		return
	})
	defer r.Close()
	log.Println("Listening on Port 8080 ... ")
	s := catalog.NewService(r)
	log.Fatal(catalog.ListenGRPC(s, 8080))
}
//protoc --proto_path=. --go_out=pb --go_opt=paths=source_relative --go-grpc_out=pb --go-grpc_opt=paths=source_relative catalog.proto
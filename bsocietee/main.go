package main

import (
	"log"

	"github.com/cristipercu/societee/bsocietee/cmd/api"
	"github.com/cristipercu/societee/bsocietee/cmd/config"
	"github.com/cristipercu/societee/bsocietee/cmd/db"
)

func main() {
  db, err := db.NewPGStorage(config.Envs.DbHost, config.Envs.DbUser, config.Envs.DbName, config.Envs.DbPassword, config.Envs.Sslmode, config.Envs.DbPort)

  if err != nil {
    log.Fatal(err)
  }

  server := api.NewAPIServer(":1620", db)

  if err := server.Run(); err != nil {
    log.Fatal(err)
  }
}

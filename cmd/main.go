package main

import (
	"log"

	"github.com/BeiymTech-hacknu2024/BeiymTech-hacknu2024-backend/internal/postgres"
)

func main() {
  DB, err := postgres.ConnectDB()
  if err != nil {
    log.Fatalf("Unable to connect to database: %v\n", err)
  }

}

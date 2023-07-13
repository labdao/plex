package models

import (
  "os"
  "fmt"
  "log"
  "gorm.io/gorm"
  "gorm.io/driver/postgres"
)

var DB *gorm.DB

func ConnectDatabase() {

    dsn := fmt.Sprintf(
      "host=%s user=%s password=%s dbname=%s",
      os.Getenv("PGHOST"),
      os.Getenv("PGUSER"),
      os.Getenv("PGPASSWORD"),
      os.Getenv("PGDATABASE"),
    )
    database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    // database, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})

    if err != nil {
      panic("Failed to connect to database!")
    }

    log.Print("Migrating database")
    err = database.AutoMigrate(&Job{})
    if err != nil {
      return
    }

    DB = database
}

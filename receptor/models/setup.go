package models

import (
  "gorm.io/gorm"
  "gorm.io/driver/sqlite"
  // "gorm.io/driver/postgres"
  "gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {

    database, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})

    if err != nil {
      panic("Failed to connect to database!")
    }

    err = database.AutoMigrate(&Job{})
    if err != nil {
      return
    }

    DB = database
}

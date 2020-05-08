package db

import (
  "genosha/utils/confs"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/postgres"
  "log"
)

// Pg for postgres handler
var Pg *gorm.DB

func init() {
  //Pg = createDBHandler()
}

func createDBHandler() *gorm.DB {
  log.Println("connecting to db...")
  connStr := "host=" + confs.FlagPGHost + " port=" + confs.FlagPGPort + " user=" + confs.FlagPGUser + " dbname=" + confs.FlagPGName + " password=" + confs.FlagPGPassword + " sslmode=disable"
  log.Printf("connecting string:%s\n", connStr)
  db, err := gorm.Open("postgres", connStr)
  if err != nil {
    panic("failed to connect database")
  }
  log.Println("postgres connected!")
  db.LogMode(true)
  return db
}

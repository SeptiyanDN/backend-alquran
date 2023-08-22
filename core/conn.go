package core

import (
	"alquran/kedaihelpers"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DBConnect() kedaihelpers.DBStruct {
	username := viper.GetString("database.username")
	password := viper.GetString("database.password")
	database := viper.GetString("database.name")
	host := viper.GetString("database.host")
	port := viper.GetInt("database.port")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, username, password, database)
	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		fmt.Println("Error Connecting DB => ", err)
		os.Exit(0)
	}

	maxLifetime, _ := time.ParseDuration(viper.GetString("database.max_lifetime_connection") + "s")
	db.SetMaxIdleConns(viper.GetInt("database.max_idle_dbection"))
	db.SetConnMaxLifetime(maxLifetime)
	dbs := kedaihelpers.DBStruct{Dbx: db}

	return dbs
}

func InitGorm() (*gorm.DB, error) {
	username := viper.GetString("database.username")
	password := viper.GetString("database.password")
	database := viper.GetString("database.name")
	host := viper.GetString("database.host")
	port := viper.GetInt("database.port")
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Tokyo", host, port, username, password, database)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the Database")
	}
	// Perform auto migration for tables
	err = db.AutoMigrate(
	// &surah.Surah{},
	// &ayah.Ayah{},
	// &stream.Stream{},
	)
	if err != nil {
		log.Fatal("Failed to perform auto migration")
	}
	fmt.Println("🚀 Connected Successfully to the Database")
	return db, nil
}

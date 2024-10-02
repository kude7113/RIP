package main

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"RIP/internal/app/ds"
	"RIP/internal/app/dsn"
)

func main() {
	println("fgsefvr")
	_ = godotenv.Load()
	env := dsn.FromEnv()
	println(env)
	db, err := gorm.Open(postgres.Open(env), &gorm.Config{})
	if err != nil {
		println("failed to connect database" + err.Error())
	}
	println("fgsefvr")
	// Migrate the schema
	if err = db.AutoMigrate(
		&ds.Fines{},
		&ds.Resolutions{},
		&ds.Fine_Resolution{},
		&ds.User{},
	); err != nil {
		println("cant migrate db")
	}
	println("fgsefvr")
}

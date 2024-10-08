package main

import (
	"RIP/internal/app/config"
	"RIP/internal/app/dsn"
	"RIP/internal/app/handler"
	"RIP/internal/app/pkg"
	"RIP/internal/app/repository"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	router := gin.Default()
	conf, errConf := config.NewConfig()
	if errConf != nil {
		logrus.Fatalln("Error with config reading: #{errConf}")
	}
	// через dsn парсим и помещаем в переменную
	postgresString := dsn.FromEnv()

	fmt.Println(postgresString)

	rep, err := repository.NewRepository(postgresString, logger)
	if err != nil {
		logrus.Fatalln("Error with repo: err", err)
	}

	hand := handler.NewHandler(logger, rep)
	application := pkg.NewApp(conf, router, logger, hand)
	application.StartServer()
}

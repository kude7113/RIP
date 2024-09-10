package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Fine struct {
	Id int 
	Title string
	FullInfo string
	Price int
}


func StartServer ()  {
	log.Println("Server start up")

	var Fines = []Fine{
		{Id: 1, Title: "Проезд по пешеходному переходу", FullInfo: "Как и на велосипеде, на самокате запрещено ехать по пешеходным переходам — необходимо спешиться и катить СИМ в руках.", Price: 2000},
		{Id: 2, Title: "Езда вдвоем на электросамокате", FullInfo: "На электросамокатах запрещено перевозить пассажиров, то есть ездить вдвоем на одном транспорте нельзя.", Price: 1500},
		{Id: 3, Title: "Передача управления электросамокатом ребенку", FullInfo: "Водить арендные самокаты можно только с 18 лет.", Price: 5000},
		{Id: 4, Title: "Езда на самокате пьяным", FullInfo: "Запрещено передвигаться на электросамокате в нетрезвом виде.", Price: 30000},
	}

	r := gin.Default()
	

	r.LoadHTMLGlob("templates/*")

	r.GET("/home", func(c *gin.Context) {
		c.HTML(http.StatusOK, "boot.html", Fines)
	})

	r.GET("/more", func(c *gin.Context) {
		c.HTML(http.StatusOK, "boot.html", Fines)
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	log.Println("Server down")
}
package api

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Fine struct {
	Id       int
	Title    string
	FullInfo string
	Price    int
}

func StartServer() {
	log.Println("Server start up")

	var Fines = []Fine{
		{Id: 1, Title: "Проезд по пешеходному переходу", FullInfo: "Как и на велосипеде, на самокате запрещено ехать по пешеходным переходам — необходимо спешиться и катить СИМ в руках.", Price: 2000},
		{Id: 2, Title: "Езда вдвоем на электросамокате", FullInfo: "На электросамокатах запрещено перевозить пассажиров, то есть ездить вдвоем на одном транспорте нельзя.", Price: 1500},
		{Id: 3, Title: "Передача управления электросамокатом ребенку", FullInfo: "Водить арендные самокаты можно только с 18 лет.", Price: 5000},
		{Id: 4, Title: "Езда на самокате пьяным", FullInfo: "Запрещено передвигаться на электросамокате в нетрезвом виде.", Price: 30000},
	}

	var Orders = []Fine{
		Fines[0],
		Fines[1],
		Fines[3],
	}

	r := gin.Default()

	r.Static("/static", "./static")

	r.LoadHTMLGlob("templates/*")

	r.GET("/home", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", gin.H{
			"fines": Fines,
		})
	})

	r.GET("/more/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.String(http.StatusBadRequest, "Invalid ID")
			return
		}

		var fine *Fine

		for _, f := range Fines {
			if f.Id == id {
				fine = &f
				break
			}
		}

		if fine == nil {
			c.String(http.StatusNotFound, "Fine not found")
			return
		}

		c.HTML(http.StatusOK, "more.html", gin.H{
			"fine": fine,
		})

	})

	r.GET("/home/search", func(c *gin.Context) {
		searchText := c.Query("searchInput")
		searchText = strings.ToLower(searchText)

		var filteredFines []Fine
		for _, Fine := range Fines {
			fineTitle := strings.TrimSpace(strings.ToLower(Fine.Title))
			if strings.HasPrefix(fineTitle, searchText) {
				filteredFines = append(filteredFines, Fine)
			}
		}

		// Возвращаем результаты поиска на страницу
		c.HTML(http.StatusOK, "home.html", gin.H{
			"fines":      filteredFines,
			"searchText": c.Query("searchInput"),
		})

	})

	r.GET("/order", func(c *gin.Context) {
		c.HTML(http.StatusOK, "order.html", Orders)
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	log.Println("Server down")
}

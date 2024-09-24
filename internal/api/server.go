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
	Imge     string
}

type Resolution struct {
	Id    int
	Fines map[Fine]int
}

func StartServer() {
	log.Println("Server start up")

	var Fines = []Fine{
		{Id: 1, Title: "Проезд по пешеходному переходу", FullInfo: "Как и на велосипеде, на самокате запрещено ехать по пешеходным переходам — необходимо спешиться и катить СИМ в руках.", Price: 2000, Imge: "http://127.0.0.1:9000/fines-bucket/0.jpg"},
		{Id: 2, Title: "Езда вдвоем на электросамокате", FullInfo: "На электросамокатах запрещено перевозить пассажиров, то есть ездить вдвоем на одном транспорте нельзя.", Price: 1500, Imge: "http://127.0.0.1:9000/fines-bucket/1.jpg"},
		{Id: 3, Title: "Передача управления электросамокатом ребенку", FullInfo: "Водить арендные самокаты можно только с 18 лет.", Price: 5000, Imge: "http://127.0.0.1:9000/fines-bucket/2.jpg"},
		{Id: 4, Title: "Езда на самокате пьяным", FullInfo: "Запрещено передвигаться на электросамокате в нетрезвом виде.", Price: 30000, Imge: "http://127.0.0.1:9000/fines-bucket/3.jpg"},
	}

	var Resolutions = []Resolution{
		{
			Id: 1,
			Fines: map[Fine]int{
				Fines[0]: 15,
				Fines[1]: 16,
			},
		},
	}

	r := gin.Default()

	r.Static("/static", "./static")

	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
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

	r.GET("/search", func(c *gin.Context) {
		searchText := c.Query("searchFines")
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
			"searchText": c.Query("searchFines"),
		})

	})

	r.GET("/order/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.String(http.StatusBadRequest, "Invalid ID")
			return
		}

		var resolution *Resolution

		for _, r := range Resolutions {
			if r.Id == id {
				resolution = &r
				break
			}
		}

		if resolution == nil {
			c.String(http.StatusNotFound, "Fine not found")
			return
		}

		c.HTML(http.StatusOK, "order.html", gin.H{
			"resolution": resolution,
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	log.Println("Server down")
}

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
	DopInf   string
}

type Resolution struct {
	Id    int
	Fines map[Fine]int
}

func StartServer() {
	log.Println("Server start up")

	var Fines = []Fine{
		{Id: 1, Title: "Проезд по пешеходному переходу", FullInfo: "Как и на велосипеде, на самокате запрещено ехать по пешеходным переходам — необходимо спешиться и катить СИМ в руках.", Price: 2000, Imge: "http://127.0.0.1:9000/fines-bucket/3.png", DopInf: "Колличество:"},
		{Id: 2, Title: "Езда вдвоем на электросамокате", FullInfo: "На электросамокатах запрещено перевозить пассажиров, то есть ездить вдвоем на одном транспорте нельзя.", Price: 1500, Imge: "http://127.0.0.1:9000/fines-bucket/4.png", DopInf: "Колличество людей на самокате:"},
		{Id: 3, Title: "Передача управления электросамокатом ребенку", FullInfo: "Водить арендные самокаты можно только с 18 лет.", Price: 5000, Imge: "http://127.0.0.1:9000/fines-bucket/5.png", DopInf: "Возраст ребенка:"},
		{Id: 4, Title: "Превышение скорости", FullInfo: "Запрещено передвигаться на электросамокате в нетрезвом виде.", Price: 30000, Imge: "http://127.0.0.1:9000/fines-bucket/2.png", DopInf: "Скорость:"},
		{Id: 5, Title: "Езда на самокате пьяным", FullInfo: "Запрещено передвигаться на электросамокате в нетрезвом виде.", Price: 2000, Imge: "http://127.0.0.1:9000/fines-bucket/1.png", DopInf: "Промилле:"},
		{Id: 6, Title: "Помеха в движении транспортного средства", FullInfo: "Запрещено передвигаться на электросамокате в нетрезвом виде.", Price: 4000, Imge: "http://127.0.0.1:9000/fines-bucket/6.png", DopInf: "Колличество транспортных средств:"},
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
			"fines":  Fines,
			"number": len(Resolutions[0].Fines),
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

	r.GET("/resolution/:id", func(c *gin.Context) {
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

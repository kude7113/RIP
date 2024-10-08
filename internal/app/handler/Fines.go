package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// вызываются функции из репы, которые идут в бд
// то есть как бы прослойка между эндпоинтами и данными, которые идут их бд

func (h *Handler) AllFines(ctx *gin.Context) {
	searchFines := ctx.Query("searchFines")

	userId := 1
	resCount, _ := h.Repository.GetResolutionLength(userId)
	resID, _ := h.Repository.HasRequestByUserID(userId)
	if searchFines == "" {
		fines, err := h.Repository.GetAllFines()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.HTML(http.StatusOK, "home.html", gin.H{
			"Fines":      fines,
			"searchText": searchFines,
			"number":     resCount,
			"ReqID":      resID,
		})
		return
	}
	fines, err := h.Repository.SearchFines(searchFines)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.HTML(http.StatusOK, "home.html", gin.H{
		"Fines":        fines,
		"searchText":   searchFines,
		"ReqCallCount": resCount,
		"ReqID":        resID,
	})

}

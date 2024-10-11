package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// вызываются функции из репы, которые идут в бд
// то есть как бы прослойка между эндпоинтами и данными, которые идут их бд

func (h *Handler) AllFines(ctx *gin.Context) {
	searchFines := ctx.Query("searchFines")

	userId := 1
	resCount, _ := h.Repository.GetResolutionLength(userId)
	resID, _ := h.Repository.ResolutionByUserID(userId)
	if searchFines == "" {
		fines, err := h.Repository.GetAllFines()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.HTML(http.StatusOK, "home.html", gin.H{
			"fines":      fines,
			"searchText": searchFines,
			"number":     resCount,
			"resID":      resID,
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
		"fines":      fines,
		"searchText": searchFines,
		"number":     resCount,
		"resID":      resID,
	})

}

func (h *Handler) FinesByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.String(http.StatusBadRequest, "Invalid ID")
		return
	}

	fine, err := h.Repository.GetFinesByID(id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.HTML(http.StatusOK, "more.html", gin.H{
		"fine": fine,
	})
}

func (h *Handler) AddFinesToRes(ctx *gin.Context) {
	fineID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid fine ID"})
		return
	}

	userID := 1
	h.Logger.Infof("Start adding fine ID: %d to user ID: %d", fineID, userID)

	err = h.Repository.AddFinesToResolution(userID, fineID)
	if err != nil {
		h.Logger.Errorf("Error adding fine to resolution: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.Logger.Infof("Successfully added fine ID: %d to user ID: %d", fineID, userID)
	ctx.Redirect(http.StatusFound, "/")
}

func (h *Handler) GetResolution(ctx *gin.Context) {
	resID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid resID"})
	}

	finesWithCount, err := h.Repository.GetFinesInResolutionById(resID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.HTML(http.StatusOK, "order.html", gin.H{
		"fine":  finesWithCount,
		"resID": resID,
	})
}

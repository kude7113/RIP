package repository

import (
	"RIP/internal/app/ds"
	"errors"
	"gorm.io/gorm"
	"strings"
)

func (r *Repository) GetFinesByID(id int) (*ds.Fines, error) {
	Fines := &ds.Fines{}

	err := r.db.First(Fines, "id = ?", "1").Error // find Fines with id = 1
	if err != nil {
		return nil, err
	}

	return Fines, nil
}

func (r *Repository) GetAllFines() (*[]ds.Fines, error) {
	var deliveryItems []ds.Fines
	r.db.Where("is_delete = ?", false).Find(&deliveryItems)
	return &deliveryItems, nil
}

func (r *Repository) GetResolutionLength(userId int) (int64, error) {
	var count int64
	var req ds.Resolutions
	status := ds.DraftStatus

	if err := r.db.Where("User_id = ? AND Status = ?", userId, status).First(&req).Error; err != nil {
		return 0, err
	}

	reqID := req.Resolution_ID

	err := r.db.Model(&ds.Fine_Resolution{}).Where("Resolution_id = ?", reqID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *Repository) HasRequestByUserID(userID int) (int, error) {
	var res ds.Fine_Resolution
	err := r.db.Where("User_id = ? AND Status = ?", userID, ds.DraftStatus).First(&res).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}
	return res.Resolution_ID, nil
}

func (r *Repository) SearchFines(searchText string) (*[]ds.Fines, error) {
	searchText = strings.ToLower(searchText)
	var Fines []ds.Fines
	// сохраняем данные из бд в массив
	r.db.Find(&Fines)

	var filteredFines []ds.Fines
	for _, Fine := range Fines {
		fineTitle := strings.TrimSpace(strings.ToLower(Fine.Title))
		if strings.HasPrefix(fineTitle, searchText) {
			filteredFines = append(filteredFines, Fine)
		}
	}
	return &filteredFines, nil
}

package repository

import (
	"RIP/internal/app/ds"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
)

type FinesWithCount struct {
	Fine_ID int
	Title   string
	FullInf string
	Price   int
	Imge    string
	DopInf  string
	Count   int
}

func (r *Repository) GetFinesByID(id int) (*ds.Fines, error) {
	Fine := &ds.Fines{}

	err := r.db.First(Fine, "Fine_ID = ?", id).Error // find Fines with id = 1
	if err != nil {
		return nil, err
	}

	return Fine, nil
}

func (r *Repository) GetAllFines() (*[]ds.Fines, error) {
	var deliveryItems []ds.Fines
	r.db.Model(&ds.Fines{}).Find(&deliveryItems)
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

func (r *Repository) ResolutionByUserID(userID int) (int, error) {
	var res ds.Resolutions
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

func (r *Repository) AddFinesToResolution(userID, fineID int) error {
	// Поиск существующей заявки пользователя со статусом 'черновик'
	var draftRequest ds.Resolutions
	err := r.db.Where("user_id = ? AND status = ?", userID, ds.DraftStatus).First(&draftRequest).Error

	// Если черновик не найден, создаём новый
	if errors.Is(err, gorm.ErrRecordNotFound) {
		draftRequest = ds.Resolutions{
			User_ID:      userID,
			Status:       ds.DraftStatus,
			Date_Created: time.Now(),
		}

		// Создание новой записи
		err = r.db.Create(&draftRequest).Error
		if err != nil {
			return fmt.Errorf("error creating new draft request: %w", err)
		}

		r.logger.Infof("Created new draft request ID: %d for user ID: %d", draftRequest.Resolution_ID, userID)
	} else if err != nil {
		// Если произошла ошибка запроса, возвращаем её
		return fmt.Errorf("error fetching draft request: %w", err)
	} else {
		r.logger.Infof("Found existing draft request ID: %d for user ID: %d", draftRequest.Resolution_ID, userID)
	}

	// Добавляем элемент в существующую заявку (или новую)
	fineRes := ds.Fine_Resolution{
		Fine_ID:       fineID,
		Resolution_ID: draftRequest.Resolution_ID,
	}

	// Вставляем в базу данных
	err = r.db.Create(&fineRes).Error
	if err != nil {
		return fmt.Errorf("error linking fine to draft request: %w", err)
	}

	r.logger.Infof("Fine ID: %d successfully added to Resolution ID: %d", fineID, draftRequest.Resolution_ID)

	return nil
}

func (r *Repository) GetFinesInResolutionById(resID int) (*[]FinesWithCount, error) {
	var finesWithCount []FinesWithCount

	// Получаем все записи Fine_Resolution для заданного resID
	var finesResolution []ds.Fine_Resolution
	err := r.db.Where("resolution_id = ?", resID).Find(&finesResolution).Error
	if err != nil {
		return nil, err
	}

	// Цикл по всем найденным штрафам в постановлении
	for _, fineRes := range finesResolution {
		var fine ds.Fines
		// Для каждого штрафа ищем его полную информацию
		err := r.db.Where("fine_id = ?", fineRes.Fine_ID).First(&fine).Error
		if err != nil {
			return nil, err
		}

		// Создаём объект FinesWithCount и добавляем его в результат
		fineCount := FinesWithCount{
			Fine_ID: fine.Fine_ID,
			Title:   fine.Title,
			FullInf: fine.FullInf,
			Price:   fine.Price,
			Imge:    fine.Imge,
			DopInf:  fine.DopInf,
			Count:   fineRes.Number, // Используем поле Number из finesResolution
		}
		finesWithCount = append(finesWithCount, fineCount)
	}

	return &finesWithCount, nil
}

package ds

type Fines struct {
	Fine_ID int    `gorm:"primaryKey"`
	Title   string `gorm:"type:varchar(255)"`
	FullInf string `gorm:"type:varchar(255)"`
	Price   int
	Imge    string `gorm:"type:varchar(255)"`
	DopInf  string `gorm:"type:varchar(255)"`
}

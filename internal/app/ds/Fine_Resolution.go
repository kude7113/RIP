package ds

type Fine_Resolution struct {
	ID            int `gorm:"primaryKey"`
	Resolution_ID int `gorm:"foreignKey:Resolution_ID"`
	Fine_ID       int `gorm:"foreignKey:Fine_ID"`
	Number        int
}

package ds

type User struct {
	User_ID  int    `gorm:"primaryKey"`
	Login    string `gorm:"type:varchar(255)"`
	Password string `gorm:"type:varchar(255)"`
}

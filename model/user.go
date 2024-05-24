package model

type User struct {
    ID          int64   `gorm:"primaryKey" json:"id"`
    Name        string  `json:"name"`
    Email       string  `gorm:"unique" json:"email"`
    Password    string  `json:"-"`
}

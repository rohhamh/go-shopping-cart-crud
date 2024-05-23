package model

type User struct {
    ID          int64   `gorm:"primaryKey"`
    Name        string
    Email       string  `gorm:"unique"`
    Password    string  `json:"-"`
}

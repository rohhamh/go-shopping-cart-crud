package model

import (
    "time"
    "gorm.io/datatypes"
)

type Cart struct {
    ID			int64           `gorm:"primaryKey"`
    CreatedAt	time.Time
    UpdatedAt	time.Time
    Data		datatypes.JSON
    State       string
    UserID      int64           `json:"-"`
    User        User            `json:"-"`
}

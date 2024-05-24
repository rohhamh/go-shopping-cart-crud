package model

import (
    "time"
    "gorm.io/datatypes"
)

type Cart struct {
    ID			int64           `gorm:"primaryKey" json:"id"`
    CreatedAt	time.Time       `json:"createdAt"`
    UpdatedAt	time.Time       `json:"updatedAt"`
    Data		datatypes.JSON  `json:"data"`
    State       string          `json:"state"`
    UserID      int64           `json:"-"`
    User        User            `json:"-"`
}

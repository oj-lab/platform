package model

import "time"

type MetaFields struct {
	CreateAt  time.Time `gorm:"autoCreateTime" json:"createAt"`
	UpdateAt  time.Time `gorm:"autoUpdateTime" json:"updateAt"`
	DeletedAt time.Time `gorm:"index" json:"deletedAt"`
}

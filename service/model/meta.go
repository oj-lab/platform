package model

import "time"

type MetaFields struct {
	CreateAt  *time.Time `gorm:"autoCreateTime" json:"createAt,omitempty"`
	UpdateAt  *time.Time `gorm:"autoUpdateTime" json:"updateAt,omitempty"`
	DeletedAt *time.Time `gorm:"index" json:"deletedAt,omitempty"`
}

var MetaFieldsSelection = []string{"create_at", "update_at", "deleted_at"}

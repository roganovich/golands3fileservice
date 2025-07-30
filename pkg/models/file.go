package models

import (
	"time"
	"github.com/google/uuid"
)

type File struct {
	ID            uuid.UUID  `db:"id" json:"id"`
	Title         string    `db:"title" json:"title,omitempty"`
	Filename      string     `db:"filename" json:"filename,omitempty"`
	Extension     string    `db:"extension" json:"extension,omitempty"`
	Size          int64      `db:"size" json:"size"`
	DateCreate    time.Time  `db:"date_create" json:"dateCreate"`
	IsDeleted     bool       `db:"is_deleted" json:"isDeleted"`
	FkOriginal    *uuid.UUID `db:"fk_original" json:"fkOriginal,omitempty"`
	DateExpiration *time.Time `db:"date_expiration" json:"dateExpiration,omitempty"`
}
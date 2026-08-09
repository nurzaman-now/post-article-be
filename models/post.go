package models

import (
	"post-article-be/enums"
	"time"
)

type Post struct {
	ID          int              `json:"id" binding:"omitempty"`
	Title       string           `json:"title" binding:"required"`
	Content     string           `json:"content" binding:"required"`
	Category    string           `json:"category" binding:"required"`
	CreatedDate time.Time        `json:"created_date" binding:"omitempty"`
	UpdatedDate *time.Time       `json:"updated_date" binding:"omitempty"`
	Status      enums.PostStatus `json:"status" binding:"required"`
}

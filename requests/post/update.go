package post

import "post-article-be/enums"

type UpdatePostRequest struct {
	Title    string           `json:"title" binding:"omitempty,min=20,max=200"`
	Content  string           `json:"content" binding:"omitempty,min=200"`
	Category string           `json:"category" binding:"omitempty,min=3,max=100"`
	Status   enums.PostStatus `json:"status" binding:"omitempty,oneof=publish draft thrash"`
}

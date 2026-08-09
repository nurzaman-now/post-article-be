package post

import "post-article-be/enums"

type CreatePostRequest struct {
	Title    string           `json:"title" binding:"required,min=20,max=200"`
	Content  string           `json:"content" binding:"required,min=200"`
	Category string           `json:"category" binding:"required,min=3,max=100"`
	Status   enums.PostStatus `json:"status" binding:"required,oneof=publish draft thrash"`
}

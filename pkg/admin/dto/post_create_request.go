package dto

type PostCreateRequest struct {
	Title   string `json:"title" validate:"required" form:"title"`
	Content string `json:"content" validate:"required" form:"content"`
}

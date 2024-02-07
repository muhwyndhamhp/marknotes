package dto

type TagFindOrCreateRequest struct {
	Tag  string   `json:"tag" validate:"required" form:"tag"`
	Tags []string `json:"tags" form:"tags"`
}

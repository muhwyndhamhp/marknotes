package models

import "gorm.io/gorm"

type Commenter struct {
	gorm.Model
	Name     string     `json:"name"`
	Email    string     `json:"email"`
	Comments []*Comment `json:"comments"`
	// is blocked default false
	IsBlocked bool `json:"is_blocked" gorm:"default:false"`
}

type Comment struct {
	gorm.Model
	Text        string     `json:"text"`
	CommenterID uint       `json:"commenter_id"`
	Commenter   *Commenter `json:"commenter"`
	PostID      uint       `json:"post_id"`
}

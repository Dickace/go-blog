package models

import (
	_ "github.com/jinzhu/gorm"
	"time"
)
type Post struct {
	ID uint `json:"id" gorm:"primary_key"`
	Title string `json:"title"`
	PostDate time.Time `json:"post_date"`
	IsFavorite bool `json:"is_favorite"`
}
type CreatePost struct {
	Title string `json:"title"`
	PostDate  time.Time `json:"post_date"`
}
type UpdatePost struct {
	Title string `json:"title"`
	PostDate  time.Time `json:"post_date" `
}

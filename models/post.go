package models

import (
	_ "github.com/jinzhu/gorm"
)
type Post struct {
	ID uint `json:"id" gorm:"primary_key"`
	Title string `json:"title"`
	PostDate string `json:"post_date"`
}
type CreatePost struct {
	Title string `json:"title" binding:"required"`
	PostDate  string `json:"post_date"`
}
type UpdatePost struct {
	Title string `json:"title"`
	PostDate  string `json:"post_date" `
}
package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var DB *gorm.DB

//const (
//	host     = "ec2-63-33-14-215.eu-west-1.compute.amazonaws.com"
//	port     = "5432"
//	user     = "efqfutmrzzunvf"
//	password = "d9fbc191dd683c84603b840e4cd07b2ae93dcf4af8959a4230f2a5e0ccf1b1a0"
//	dbname   = "d4obgiuuk72amd"
//)

const (
	host     = "localhost"
	port     = "5432"
	user     = "postgres"
	password = "password"
	dbname   = "blog"
)

func OpenConnection() {
	psqlConnStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := gorm.Open("postgres", psqlConnStr)
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Post{})

	DB = db
}

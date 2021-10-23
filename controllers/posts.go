package controllers

import (
	"awesomeProject/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddTelemetry(c *gin.Context){
	var input models.CreateTelemetry
	if err:= c.ShouldBindJSON(&input); err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
		return
	}
	telemtry := models.Telemetry{Interface: input.nterface, ICCID: input.ICCID, Tamb_degC: input.Tamb_degC, AX: input.AX ,AY: input.AY, AZ: input.AZ,  RSSI_dBm: input.RSSI_dBm, Latitude: input.Latitude, Longitude: input.Longitude, GNSS_data_valid: input.GNSS_data_valid }
	models.DB.Create(&telemtry)
	c.JSON(http.StatusOK, gin.H{"data": telemtry})
}

func GetTelemetries(c *gin.Context){
	var telemtries []models.Telemetry
	models.DB.Find(&telemtries)
	c.JSON(http.StatusOK, gin.H{"data": telemtries})
}

func GetPosts(c *gin.Context){
	var posts []models.Post
	models.DB.Find(&posts)
	c.JSON(http.StatusOK,gin.H{"data":posts})
}
func GetPost(c *gin.Context){
	var post models.Post
	if err:= models.DB.Where("id = ?", c.Param("id")).First(&post).Error; err!= nil{
		c.JSON(http.StatusBadRequest, gin.H{"error":"Record not found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": post})
}
func EditPost(c *gin.Context){

	var post models.Post
	if err:= models.DB.Where("id = ?", c.Param("id")).First(&post).Error;err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{"error":"record not found"})
		return
	}
	var input models.UpdatePost
	if err:= c.ShouldBindJSON(&input); err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	models.DB.Model(&post).Updates(input)
	c.JSON(http.StatusOK, gin.H{"data": post})

}
func DeletePost(c *gin.Context){
	var post models.Post
	if err:= models.DB.Where("id = ?", c.Param("id")).First(&post).Error;err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{"error":"post not found"})
		return
	}
	models.DB.Delete(&post)
	c.JSON(http.StatusOK, gin.H{"data":true})
}
func AddPost(c *gin.Context){
	var input models.CreatePost

	if err:= c.ShouldBindJSON(&input); err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	post := models.Post{Title: input.Title, PostDate: input.PostDate}
	models.DB.Create(&post)
	c.JSON(http.StatusOK, gin.H{"data": post})
}
func Favorite(c *gin.Context){
	var post models.Post

	if err:= models.DB.Where("id = ?", c.Param("id")).First(&post).Error; err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{"error":"post not found"})
	}
	if !post.IsFavorite {
		models.DB.Model(&post).Updates(map[string]interface{}{
			"IsFavorite": true,
		})
	} else {
		models.DB.Model(&post).Updates(map[string]interface{}{
			"IsFavorite": false,
		})
	}
	c.JSON(http.StatusOK, gin.H{"data": post})
}
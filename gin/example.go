package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB
var err error

type Person struct {
	ID       uint   `json:"id"`
	GameName string `json:"title"`
	Alias    string `json:"alias"`
}

func main() {
	// NOTE: See we're using = to assign the global var
	// instead of := which would assign it only in this functionias
	db, err = gorm.Open("sqlite3", "./gorm.db")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	db.AutoMigrate(&Person{})
	r := gin.Default()

	r.GET("/game", GetProjects)
	r.GET("/game/:alias", GetPerson)
	r.POST("/game", CreatePerson)
	r.PUT("/game/:alias", UpdatePerson)
	r.DELETE("/game/:alias", DeletePerson)

	r.Run(":8080")
}

func GetProjects(c *gin.Context) {
	var game []Person
	if err := db.Find(&game).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, game)
	}
}

func GetPerson(c *gin.Context) {
	alias := c.Params.ByName("alias")
	var person Person
	if err := db.Where("alias LIKE ?", "%"+alias+"%").Find(&person).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, person)
	}
}

func CreatePerson(c *gin.Context) {
	var person Person
	c.BindJSON(&person)
	db.Create(&person)
	c.JSON(200, person)
}

func UpdatePerson(c *gin.Context) {
	var person Person
	alias := c.Params.ByName("alias")
	if err := db.Where("alias = ?", alias).First(&person).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	}
	c.BindJSON(&person)
	db.Save(&person)
	c.JSON(200, person)
}

func DeletePerson(c *gin.Context) {
	alias := c.Params.ByName("alias")
	var person Person
	d := db.Where("alias = ?", alias).Delete(&person)
	fmt.Println(d)
	c.JSON(200, gin.H{"alias #" + alias: "deleted"})
}

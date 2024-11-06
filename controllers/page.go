package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maximka76667/sigma-go-rest-api/database"
	"github.com/maximka76667/sigma-go-rest-api/models"
)

func CreatePage(c *gin.Context) {
	var page models.Page

	// Bind JSON to Page struct
	if err := c.ShouldBindJSON(&page); err != nil {
		// Error on binding (incorrect format request data)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create page in database
	if err := database.DB.Create(&page).Error; err != nil {
		// Error on creating
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return newly created page
	c.JSON(http.StatusCreated, page)
}

func GetPages(c *gin.Context) {
	var pages []models.Page
	if err := database.DB.Find(&pages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pages"})
		return
	}

	c.JSON(http.StatusOK, pages)
}

func GetPageByGUID(c *gin.Context) {
	guid := c.Param("guid")
	var page models.Page

	if err := database.DB.Where("page_guid = ?", guid).First(&page).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, page)
}

func UpdatePage(c *gin.Context) {
	guid := c.Param("guid")
	var page models.Page

	// Chect if page already exists
	if err := database.DB.Where("page_guid = ?", guid).First(&page).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
	}

	// Bind JSON to Page struct
	if err := c.ShouldBindJSON(&page); err != nil {
		// Error on binding (incorrect format request data)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Save(&page).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update page"})
		return
	}

	c.JSON(http.StatusOK, page)
}

func DeletePage(c *gin.Context) {
	id := c.Param("id")
	var page models.Page

	// Chect if page already exists
	if err := database.DB.Delete(page, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
	}

	c.JSON(http.StatusOK, page)
}

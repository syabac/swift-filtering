package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"bni.co.id/swift-filtering/config"
	"bni.co.id/swift-filtering/scenario"
)

// CacheController handle route to manage settings
type CacheController struct {
}

// NewCache Create new setting controller
func NewCache() *CacheController {
	return &CacheController{}
}

// ClearSettings clear settings cache from memory and load new from database
func (sc *CacheController) ClearSettings(c *gin.Context) {
	config.ReloadSettings()

	c.JSON(http.StatusOK, gin.H{
		"message": "Settings reloaded",
		"values":  config.GetSettings(),
	})
}

// ClearRules clear rules cache from memory
func (sc *CacheController) ClearRules(c *gin.Context) {
	scenario.ClearRulesCache()

	c.JSON(http.StatusOK, gin.H{
		"message": "Rules Cache cleared",
	})
}

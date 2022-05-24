package scheduler

import (
	"log"
	"time"
	
	"bni.co.id/swift-filtering/config"
	"bni.co.id/swift-filtering/scenario"
)

// StartClearCacheJob clear cache from memory every 1-hour
func StartClearCacheJob(){
	go func() {
		for true {
			time.Sleep(time.Hour)
			log.Println("executing ClearCacheJob...")
			config.ReloadSettings()
			scenario.ClearRulesCache()
		}
	}()
	
	log.Println("ClearCacheJob started")
}
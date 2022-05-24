package config

import (
	"bni.co.id/swift-filtering/database"
)

var settingsCache Settings = Settings{}
var settingLoaded = false

type settingData struct {
	SettingName  string
	SettingValue string
}

// ReloadSettings get  settings from database
func ReloadSettings() {
	var rs []settingData
	var db = database.Open()

	db.Raw("SELECT setting_name, setting_value FROM sys_setting").Scan(&rs)

	for _, sett := range rs {
		settingsCache[sett.SettingName] = sett.SettingValue
	}
}

// LoadSettings  lo ad settings if not initialize
func LoadSettings() {
	if settingLoaded {
		return
	}

	settingLoaded = true
	ReloadSettings()
}

// GetSettings get current setting values
func GetSettings() Settings {
	return settingsCache
}

// GetSetting get setting value based on setting name
func GetSetting(name string) string {
	var val, exist = settingsCache[name]

	if exist {
		return val
	}

	return ""
}

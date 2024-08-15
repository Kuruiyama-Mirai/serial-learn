package models

import (
	"time"
)

type SaveData struct {
	ID                     uint      `gorm:"primaryKey;" json:"id"`
	DeviceID               int64     `gorm:"column:deviceId;" json:"deviceId"`
	Temperature            float64   `gorm:"column:temperature;" json:"temperature"`
	TemperatureAlarmStatus int       `gorm:"column:temperatureAlarmStatus;" json:"temperatureAlarmStatus"`
	DIStatus               int       `gorm:"column:diStatus;" json:"DIStatus"`
	SOEpointer             int       `gorm:"column:soepointer;" json:"SOEpointer"`
	VoltageImbalance       float64   `gorm:"column:voltageImbalance;" json:"voltageImbalance"`
	CurrentImbalance       float64   `gorm:"column:currentImbalance;" json:"currentImbalance"`
	CreatedAt              time.Time `gorm:"column:created_at"`
}

func (table *SaveData) TableName() string {
	return "save_data"
}

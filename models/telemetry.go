package models

import (
	_ "github.com/jinzhu/gorm"
)
type Telemetry struct {
	ID uint `json:"id" gorm:"primary_key"`
	Interface string `json:"interface"`
	ICCID string `json: "ICCID"`
	Tamb_degC int `json: "Tamb_degC"`
	AX int `json: "aX"`
	AY int `json: "aY"`
	AZ int `json: "aZ"`
	RSSI_dBm int `json: "RSSI_dBm"`
	Latitude float32 `json: "latitude"`
	Longitude float32 `json: "longitude"`
	GNSS_data_valid uint `json: "GNSS_data_valid"`
}
type CreatePost struct {
	Interface string `json:"interface"`
	ICCID string `json: "ICCID"`
	Tamb_degC int `json: "Tamb_degC"`
	AX int `json: "aX"`
	AY int `json: "aY"`
	AZ int `json: "aZ"`
	RSSI_dBm int `json: "RSSI_dBm"`
	Latitude float32 `json: "latitude"`
	Longitude float32 `json: "longitude"`
	GNSS_data_valid uint `json: "GNSS_data_valid"`
}

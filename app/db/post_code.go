package db

import (
	"gorm.io/gorm"
)

type PostCode struct {
  gorm.Model
  Code string `json:"code" gorm:"index"`
  PrefectureName string `json:"prefecture_name"`
  CityName string `json:"city_name"`
  TownAreaName string `json:"town_area_name"`
}

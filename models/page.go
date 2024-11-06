package models

import "html/template"

type Page struct {
	ID      uint          `json:"id" gorm:"primaryKey" db:"id"`
	Title   string        `json:"page_title" gorm:"column:page_title"`
	Content template.HTML `json:"page_content" gorm:"column:page_content"`
	Date    string        `json:"page_date" gorm:"column:page_date"`
	GUID    string        `json:"page_guid" gorm:"unique;column:page_guid"`
}

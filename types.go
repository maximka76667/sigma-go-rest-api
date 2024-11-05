package main

import "html/template"

type Page struct {
	Title      string
	Content    template.HTML
	RawContent string
	Date       string
	GUID       string
}

type JSONResponse struct {
	Fields map[string]string
}

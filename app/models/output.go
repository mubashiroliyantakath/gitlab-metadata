package models

import "text/template"

type Output struct {
	Template *template.Template
	FileName string
	Data     interface{}
}

package controllers

import (
	"html/template"
)

var FuncMap = template.FuncMap{
	"add": Add,
}

func Add(a int) int {
	return a + 1
}

package util

import (
	"text/template"

	"github.com/Masterminds/sprig"
)

var NameFuncMap = func() template.FuncMap {
	m := sprig.TxtFuncMap()
	return m
}()

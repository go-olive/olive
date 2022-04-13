package main

import (
	"github.com/go-olive/olive/app"
	_ "github.com/go-olive/olive/internal"
)

func main() {
	app.NewDevice().Run()
}

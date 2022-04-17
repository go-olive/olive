package main

import (
	"github.com/go-olive/olive/src/app"
	_ "github.com/go-olive/olive/src/internal"
)

func main() {
	app.NewDevice().Run()
}

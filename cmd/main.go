package main

import (
	"github.com/luxcgo/lifesaver/app"
	_ "github.com/luxcgo/lifesaver/internal"
)

func main() {
	app.NewDevice().Run()
}

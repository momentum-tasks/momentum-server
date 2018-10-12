package main

import (
	"github.com/momentum-tasks/momentum-server/api/app"
)

func main() {
	app.NewDBManager().Begin("mysql", "momentumuser:password@tcp(database:3306)/momentum?parseTime=true")
	app.NewRoutes().Begin(":5000")
}

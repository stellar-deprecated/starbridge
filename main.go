package main

import "github.com/stellar/starbridge/app"

func main() {
	app := app.NewApp(app.Config{
		Port:      8000,
		AdminPort: 6666,
	})
	go app.RunHTTPServer()
	go app.RunBackendWorker()
	ch := make(chan bool)
	<-ch
}

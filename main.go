package main

func main() {
	app := NewApp(Config{
		Port:      8000,
		AdminPort: 6666,
	})
	go app.RunHTTPServer()
	go app.RunBackendWorker()
	ch := make(chan bool)
	<-ch
}

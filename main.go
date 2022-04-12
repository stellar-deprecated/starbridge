package main

func main() {
	app := NewApp(Config{
		Port:      8000,
		AdminPort: 6666,
	})
	app.Run()
}

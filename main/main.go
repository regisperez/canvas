package main

func main() {
	a := App{}
	a.Initialize("./app.toml")
	a.Run("127.0.0.1")
}
package main

func main() {
	a := App{}
	a.Initialize("./app.toml")
	a.Run("0.0.0.0:8010")
}
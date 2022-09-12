package main

import "fmt"

func main() {
	a := App{}
	a.Initialize("./app.toml")
	a.Run("0.0.0.0:8010")
	fmt.Println("Now we are in AWS server with docker compose")
}
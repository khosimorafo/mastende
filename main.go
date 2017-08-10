package main

import "os"

func main() {

	port := os.Getenv("PORT")

	a := MastendeQL{}

	a.Initialize()

	a.Run(":" + port)
}
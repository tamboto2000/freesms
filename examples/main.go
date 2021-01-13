package main

import (
	"github.com/tamboto2000/freesms"
)

func main() {
	cl, err := freesms.NewClient()
	if err != nil {
		panic(err.Error())
	}

	if err := cl.SendMsg("085711502721", "Jangan lupa makan ya sanyaaang"); err != nil {
		panic(err.Error())
	}
}

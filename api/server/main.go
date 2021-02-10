package main

import (
	"github.com/noahfriedman-ca/quick-factor/api"
	"unit.nginx.org/go"
)

func main() {
	if e := unit.ListenAndServe(":8080", api.Router()); e != nil {
		panic(e)
	}
}

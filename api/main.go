package main

import "unit.nginx.org/go"

func main() {
	if e := unit.ListenAndServe(":8080", nil); e != nil {
		panic(e)
	}
}

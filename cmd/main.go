package main

import "zencache/internal/transport/http"

func main() {
	s := http.New(":8080")
	s.Run()
}

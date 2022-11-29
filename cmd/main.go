package main

import (
	"power-tech-test/internal"
	"sync"
)

func main() {
	pgStorage := internal.NewPGStorage("postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")

	wg := sync.WaitGroup{}
	wg.Add(200)
	for i := 1; i <= 200; i++ {
		if i > 10 && i < 20 {
			go pgStorage.Decrease(&wg)
			continue
		}
		go pgStorage.Decrease(&wg, i)
	}
	wg.Wait()
}
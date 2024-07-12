package main

import (
	"github.com/monishth/go-therm/internal/engine"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go engine.RunEngine(&wg)
	wg.Wait()
}

package main

import (
	"fmt"
	"log"
	"sync"
)

const (
	VERSION = "2.4"
)

func main() {
	fmt.Printf("[INFO] switcher version:  %s", VERSION)
	fmt.Println("[INFO] switcher is running ")

	log.Printf("[INFO] switcher version:  %s", VERSION)
	log.Println("[INFO] switcher is running ")

	wg := &sync.WaitGroup{}
	for _, v := range config.Rules {
		wg.Add(1)
		go listen(v, wg)
	}
	wg.Wait()
	fmt.Printf("[INFO] switcher exited")
	log.Printf("[INFO] switcher exited")
}

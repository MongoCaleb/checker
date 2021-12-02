package main

import (
	rstparser "checker/rstparser"
	"time"
)

func main() {
	sayHi()
	time.Sleep(time.Second)
}

func sayHi() {
	go func() {
		for i := 0; i < 10; i++ {
			rstparser.Parse()
		}
	}()
}

package main

import (
	"fmt"
	"time"
)

func main() {
	var entrada1 chan int = make(chan int)
	var entrada2 chan int = make(chan int)
	var entrada3 chan int = make(chan int)
	var saida chan int = make(chan int)

	go gerador(entrada1, 2)
	go gerador(entrada2, 3)
	go gerador(entrada3, 5)

	go merge(entrada1, entrada2, entrada3, saida)

	go consumer(saida, 1)

	fin := make(chan struct{})
	<-fin
}

func gerador(entrada chan int, multiplicador int) {
	for i := 0; true; i++ {
		time.Sleep(500 * time.Millisecond)
		entrada <- i * multiplicador
	}
}

func consumer(saida chan int, t int) {
	for {
		x := <-saida
		fmt.Print(x, " ")
	}
}

func merge(entrada1, entrada2, entrada3, out chan int) {
	v1 := <-entrada1
	v2 := <-entrada2
	v3 := <-entrada3
	for {
		min := v1 // acha o menor
		if v2 < min {
			min = v2
		}
		if v3 < min {
			min = v3
		}
		out <- min
		// le proximo valor da serie do menor
		if min == v1 {
			v1 = <-entrada1
		}
		if min == v2 {
			v2 = <-entrada2
		}
		if min == v3 {
			v3 = <-entrada3
		}
	}
}

//MCC - Fernando Dotti

package main

import (
	"fmt"
	"strconv"
)

const (
	PHILOSOPHERS = 5
	FORKS        = 5
)

func philosopher(id int, first_fork chan struct{}) {
	for {
		fmt.Println(strconv.Itoa(id) + " senta")
		<-first_fork // pega
		fmt.Println(strconv.Itoa(id) + " pegou direita")
		fmt.Println(strconv.Itoa(id) + " come")
		first_fork <- struct{}{} // devolve
		fmt.Println(strconv.Itoa(id) + " levanta e pensa")
	}
}

func main() {
	var fork_channels [FORKS]chan struct{}
	for i := 0; i < FORKS; i++ {
		fork_channels[i] = make(chan struct{}, 1)
		fork_channels[i] <- struct{}{} // no inicio garfo esta livre
	}
	for i := 0; i < (PHILOSOPHERS); i++ {
		fmt.Println("Filosofo " + strconv.Itoa(i))
		go philosopher(i, fork_channels[i])
	}
	<-make(chan struct{})
}

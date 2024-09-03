// MCC - Fernando Dotti

package main

import (
	"fmt"
	"strconv"
)

const (
	PHILOSOPHERS = 5
	FORKS        = 5
)

type Fork struct {}

func philosopher(id int, first_fork, second_fork chan Fork) {
	for {
		// 2nd Solution to deadlock problem
		// If you get the first fork but not the second, put the first fork back and try again.
		// If you get the second fork but not the first, put the second fork back and try again.
		// Only if you get both forks, eat and then put them back.

		fmt.Println(strconv.Itoa(id) + " senta")
		<-first_fork // pega
		fmt.Println(strconv.Itoa(id) + " pegou direita")
		select {
		case <-second_fork:
			fmt.Println(strconv.Itoa(id) + " pegou esquerda")
			fmt.Println(strconv.Itoa(id) + " come")
			first_fork <- Fork{} // devolve
			second_fork <- Fork{}
			fmt.Println(strconv.Itoa(id) + " levanta e pensa")
		default:
			fmt.Println(strconv.Itoa(id) + " nao pegou esquerda")
			first_fork <- Fork{} // devolve
			fmt.Println(strconv.Itoa(id) + " levanta e pensa")
		}
	}
}

func main() {
	var fork_channels [FORKS]chan Fork
	for i := 0; i < FORKS; i++ {
		fork_channels[i] = make(chan Fork, 1)
		fork_channels[i] <- Fork{} // no inicio garfo esta livre
	}
	for i := 0; i < (PHILOSOPHERS); i++ {
		fmt.Println("Filosofo " + strconv.Itoa(i))
		// 1st Solution to deadlock problem: odd philosophers pick up left fork first, even philosophers pick up right fork first.
		if i%2 == 0 {
			go philosopher(i, fork_channels[i], fork_channels[(i+1)%PHILOSOPHERS])
		} else {
			go philosopher(i, fork_channels[(i+1)%PHILOSOPHERS], fork_channels[i])
		}
	}
	<- make(chan struct{})
}
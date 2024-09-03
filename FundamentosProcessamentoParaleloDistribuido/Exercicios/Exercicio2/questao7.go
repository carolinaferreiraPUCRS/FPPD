package main

import (
	"fmt"
)

func SomaSC() {
	var sharedTest int = 0
	var ch_fim chan struct{} = make(chan struct{})
	var sem chan struct{} = make(chan struct{}, 1)

	for i := 0; i < 100; i++ {
		go func() {
			for k := 0; k < 100; k++ {
				chan <- struct{}
				sharedTest = sharedTest + 1
				<-chan
			}
			ch_fim <- struct{}
		}()
	}
	fmt.Println("Resultado ", sharedTest)
}

func main() {
	for i := 0; i < 10; i++ {
		go SemaSC()
	}
	fmt.Println("Criei 20 processos")
	for i := 0; i < 20; i++ {
		<-ch_fim
	}
	fmt.Println("Processos acabaram. Resultado ", sharedTest)
}
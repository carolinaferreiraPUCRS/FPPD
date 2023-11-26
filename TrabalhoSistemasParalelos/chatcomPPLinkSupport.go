// Construido como parte da disciplina: Sistemas Distribuidos - PUCRS - Escola Politecnica
//  Professor: Fernando Dotti  (https://fldotti.github.io/)

/*
LANCAR 2 PROCESSOS EM SHELL's DIFERENTES, PARA CADA PROCESSO, O SEU PROPRIO ENDERECO EE O PRIMEIRO DA LISTA
go run chatComPPLink.go   127.0.0.1:5001  127.0.0.1:6001
go run chatComPPLink.go   127.0.0.1:6001  127.0.0.1:5001
ou, claro, fazer distribuido trocando os ip's
*/

package main

import (
	PP2PLink "SD/PP2PLink"
	// "bufio"
	"fmt"
	"os"
	"time"
)

const (
	nrSizes  = 6  // numero de magnitudes dos valores primos
	nrPrimos = 10 // numero de valores primos para cada magnitude
)

func contaPrimosConc(s [nrPrimos]int) (int, time.Duration) {
	end := make(chan int)
	start := time.Now()
	contador := 0
	for i := 0; i < nrPrimos; i++ {
		go func(i int) {
			if isPrime(s[i]) {
				contador++ //+1
			}
			end <- 1
		}(i)
	}
	for i := 0; i < nrPrimos; i++ {
		<-end
	}
	return contador, time.Since(start) //contador
}

func isPrime(p int) bool {
	if p%2 == 0 {
		return false
	}
	for i := 3; i*i <= p; i += 2 {
		if p%i == 0 {
			return false
		}
	}
	return true
}

func sendResponse(response string, address string, lk *PP2PLink) {
	req := PP2PLink.PP2PLink_Req_Message{
		To:      address,
		Message: response}

	lk.Req <- req
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage:   go run pp2plTest.go thisProcessIpAddress:port otherProcessIpAddress:port")
		fmt.Println("Example: go run pp2plTest.go  127.0.0.1:8050    127.0.0.1:8051")
		fmt.Println("Example: go run pp2plTest.go  127.0.0.1:8051    127.0.0.1:8050")
		return
	}

	addresses := os.Args[1:]
	fmt.Println("Chat PPLink - addresses: ", addresses)

	lk := PP2PLink.NewPP2PLink(addresses[0], false)

	for {
		primesStr := <-lk.Ind //string que possui array de numeros --> "2, 3, 5, 7..."
		fmt.Println("                                            Rcv: ", primesStr)
		var primes [10]int
		array := [10]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10} 
		primes = array // primesStr.split(", ") --> [2, 3, 5, 7]
		primeCount, duration := contaPrimosConc(primes)
		message := fmt.Sprintln(primeCount, " ", duration)
		go sendResponse(message, addresses[1], lk)
	}

	// recebe array de numeros --> "15, 16, 17, 18"
	// calcula quais são primos
	// retorna array de primos

}

// Nome dos integrantes: Arthur Both, Carolina Ferreira, Felipe Freitas, Gabriel Ferreira e Mateus Caçabuena.

package main

import (
	PP2PLink "SD/PP2PLink"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func contaPrimosConc(s []int) (int, time.Duration) {
	end := make(chan int)
	start := time.Now()
	contador := 0
	for i := 0; i < len(s); i++ {
		go func(i int) {
			if isPrime(s[i]) {
				contador++
			}
			end <- 1
		}(i)
	}
	for i := 0; i < len(s); i++ {
		<-end
	}
	return contador, time.Since(start)
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

func sendResponse(response string, address string, lk *PP2PLink.PP2PLink) {
	req := PP2PLink.PP2PLink_Req_Message{
		To:      address,
		Message: response}

	lk.Req <- req
}

// Esta função é executada quando o programa é iniciado
// Recebe os endereços dos processos que serão conectados
// Cria um link entre os processos
// Recebe uma mensagem do processo conectado
// Calcula a primalidade dos números recebidos
// Envia a resposta com quantidade de primos e duração do cálculo
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:   go run chatComPPLinkSupport.go thisProcessIpAddress:port otherProcessIpAddress:port")
		thisProcess := "127.0.0.1:8051"
		otherProcess := "127.0.0.1:8050"
		fmt.Printf("Example: go run chatComPPLinkSupport.go %s %s\n", thisProcess, otherProcess)
		return
	}

	addresses := os.Args[1:]
	fmt.Println("Chat PPLink - addresses: ", addresses)

	lk := PP2PLink.NewPP2PLink(addresses[0], false)

	for {
		msg := <-lk.Ind
		// Mensagem que possui os números por se computar primalidade --> "2,3,5,7..."
		message := msg.Message
		// Separar mensagem em duas partes --> "... Magnitude = 3||2,3,5,7..."
		pieces := strings.Split(message, " || ")
		// Separar magnitude e números --> ["Magnitude = 3", "2,3,5,7"]
		magnitudeMessage, numbersStr := pieces[0], pieces[1]
		// Separar valor magnitude --> ["Magnitude", "3"]
		magnitude := strings.Split(magnitudeMessage, " = ")[1]
		// Transformar em array --> numbers.split(",") --> ["2", "3", "5", "7"]
		numbersArrayStr := strings.Split(numbersStr, ",")
		// Converter para números --> [2, 3, 5, 7]
		var primes []int
		for _, n := range numbersArrayStr {
			number, _ := strconv.Atoi(n)
			primes = append(primes, number)
		}
		// Calcular primalidade
		primeCount, duration := contaPrimosConc(primes)
		response := fmt.Sprintf("Magnitude, duration, primeCount = %s || %s || %d", magnitude, duration, primeCount)
		go sendResponse(response, addresses[1], lk)
	}
}

// por Fernando Dotti - PUCRS -
// utilize este programa para gerar nPr números primos de tamanho conforme número de 9s em rand.Intn(9999999999999)
//

package main

import (
	"fmt"
	"math/rand"
	"time"
)

const N = 2000
const NPr = 10

var primos = make([]int, NPr)

func main() {
	fmt.Println("------ DIFFERENT prime count IMPLEMENTATIONS -------")
	slice := generateSlice(N)
	p := contaPrimosSeq(slice)

	fmt.Println("  ------ n primos :  ", p)
	fmt.Println("  Primos :  ", primos)
}

// Generates a slice of size, size filled with random numbers
func generateSlice(size int) []int {
	slice := make([]int, size, size)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		slice[i] = rand.Intn(9999999999999999)
	}
	return slice
}

func contaPrimosSeq(s []int) int {
	result := 0
	for i := 0; i < N && result < 10; i++ {
		if isPrime(s[i]) {
			fmt.Println(" ")
			fmt.Println("  ------ primos :  ", s[i])
			if result < NPr {
				primos[result] = s[i]
			}
			result++
		} else {
			fmt.Print(" .")
		}
	}
	return result
}

// Is p prime?
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

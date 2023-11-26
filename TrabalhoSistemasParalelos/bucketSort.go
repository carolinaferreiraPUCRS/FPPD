// Nome dos integrantes: Arthur Both, Carolina Ferreira, Felipe Freitas, Gabriel Ferreira e Mateus Caçabuena.

package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"slices"
	"sync"
	"time"
)

func generateRandomSlice(size, interval int) []int {
    slice := make([]int, size, size)
    rand.Seed(time.Now().UnixNano())
    for i := 0; i < size; i++ {
        slice[i] = rand.Intn(interval)
    }
    return slice
}

func main() {
	for i := 0; i < 16; i++ {
		exec(1000000, 1000000000, i + 1)
	}
}

func exec(numberCount, interval, nrBuckets int) {
	// var numberCount, interval, nrBuckets int = 1000000, 1000000000, 0
	
	// fmt.Print("Digite o numero de elementos a serem ordenados: ")
	// fmt.Scan(&numberCount)
	// for numberCount <= 0 {
	// 	fmt.Print("O numero de elementos deve ser positivo. Digite novamente: ")
	// 	fmt.Scan(&numberCount)
	// }

	// fmt.Print("Digite o limite superior dos numeros a serem ordenados: ")
	// fmt.Scan(&interval)
	// for interval <= 0 {
	// 	fmt.Print("O limite superior deve ser positivo. Digite novamente: ")
	// 	fmt.Scan(&interval)
	// }

	// fmt.Print("Digite o numero de buckets (processadores) a serem usados: ")
	// fmt.Scan(&nrBuckets)
	// for nrBuckets <= 0 {
	// 	fmt.Print("O numero de buckets deve ser positivo. Digite novamente: ")
	// 	fmt.Scan(&nrBuckets)
	// }
	runtime.GOMAXPROCS(nrBuckets)
	var wg sync.WaitGroup
	wg.Add(nrBuckets)

	numbersToSort := generateRandomSlice(numberCount, interval)
	buckets := make([][]int, nrBuckets)
	finalArray := make([]int, 0)

	fmt.Printf("Sorting %d numbers from 0 to %d using %d buckets\n", numberCount, interval - 1, nrBuckets)

	startTime := time.Now()
	for i := 0; i < nrBuckets; i++ {
		go func(i int) {
			for _, number := range numbersToSort {
				// fmt.Println("Bucket", i, "checking number", number)
				if number >= interval / nrBuckets * i && number < interval / nrBuckets * (i + 1) {
					// fmt.Println("Number", number, "in bucket", i)
					buckets[i] = append(buckets[i], number)
				}
			}
			// fmt.Println("Bucket", i, ": ", buckets[i])
			slices.Sort(buckets[i])
			// fmt.Println("Sorted bucket", i, ": ", buckets[i])
			wg.Done()
		}(i)
	}

	wg.Wait()
	for i := 0; i < nrBuckets; i++ {
		finalArray = append(finalArray, buckets[i]...)
	}
	timeElapsed := time.Since(startTime)
	fmt.Println("Time elapsed:", timeElapsed)
	// fmt.Println("Is final array sorted?", slices.IsSorted(finalArray))
	// if numberCount <= 1000 {
	// 	fmt.Println("Final array:", finalArray)
	// } else {
	// 	fmt.Println("Final array:", finalArray[:50], "...", finalArray[len(finalArray) - 50:])
	// }
}

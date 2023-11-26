// por Fernando Dotti - PUCRS -
// 		este programa calcula o tempo para detectar que os valores dos diversos arrays em "todosPrimos" são primos
//      note que os diferentes arrays tem primos com diferentes magnitudes

package main

import (
	"fmt"
	"runtime"
	"time"
)

const (
	nrSizes  = 6  // numero de magnitudes dos valores primos
	nrPrimos = 10 // numero de valores primos para cada magnitude
)

func main() {
	fmt.Print("Digite o numero de processadores a serem usados: ")
	var n int
	fmt.Scan(&n)
	runtime.GOMAXPROCS(n) // usando n processadores

	//  valores primos com respectivamente 3, 6, 9, 13, 18 casas
	//  use o programa AchaNPrimos para achar primos com determinado número de casas

	primos3 := [nrPrimos]int{101, 883, 359, 941, 983, 859, 523, 631, 181, 233}
	primos6 := [nrPrimos]int{547369, 669437, 683251, 610279, 851117, 655439, 937351, 419443, 128467, 316879}
	primos9 := [nrPrimos]int{550032733, 429415309, 109543211, 882936113, 546857209, 756170741, 699422809, 469062577, 117355333, 617320027}
	primos13 := [nrPrimos]int{7069402558433, 960246047869, 5738081989711, 5358141480883, 2569391599009, 4135462531597, 7807787948171, 130788041233, 2708131414819, 1571981553097}
	primos16 := [nrPrimos]int{2207749090466833, 9361721528139247, 2657959759011013, 3551950148669023, 3460183118669741, 5503892014624961, 4067979800826917, 7848969908399551, 2806933754138389, 5211072635754109}
	primos18 := [nrPrimos]int{383376390724197361, 882611655919772761, 533290385325847007, 17969611178168479, 903013501582628521, 541906710014517121, 281512690206248899, 403936627075987639, 775148726422474717, 942319117335957539}

	todosPrimos := [nrSizes][nrPrimos]int{primos3, primos6, primos9, primos13, primos16, primos18}

	for p := 0; p < nrSizes; p++ {
		fmt.Println("  ****** Valores avaliados :  ", todosPrimos[p])
		res := contaPrimosSeq(todosPrimos[p])
		fmt.Println(" ")
		fmt.Println("\n  ------ Tempo Seq: ", res)
		end := make(chan int)
		res = contaPrimosConc(todosPrimos[p], end)
		fmt.Println("\n  ------ Tempo Conc: ", res)
	}
}

func contaPrimosSeq(s [nrPrimos]int) time.Duration {
	start := time.Now()
	for i := 0; i < nrPrimos; i++ {
		if isPrime(s[i]) {
			fmt.Print(" .")
		}
	}
	return time.Since(start)
}

func contaPrimosConc(s [nrPrimos]int, end chan int) time.Duration {
	start := time.Now()
	for i := 0; i < nrPrimos; i++ {
		go func(i int) {
			if isPrime(s[i]) {
				fmt.Print(" .") //+1
			}
			end <- 1
		}(i)
	}
	for i := 0; i < nrPrimos; i++ {
		<-end
	}
	return time.Since(start) //contador
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

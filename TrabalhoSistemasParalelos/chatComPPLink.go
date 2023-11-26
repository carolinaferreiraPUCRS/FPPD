// Nome dos integrantes: Arthur Both, Carolina Ferreira, Felipe Freitas, Gabriel Ferreira e Mateus Caçabuena.

package main

import (
	PP2PLink "SD/PP2PLink"
	"fmt"
	"os"
	"strings"
)

const (
	nrPrimos = 10 // numero de valores primos para cada magnitude
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:   go run chatComPPLink.go thisProcessIpAddress:port otherProcessIpAddress:port")
		thisProcess := "127.0.0.1:8050"
		otherProcess := "127.0.0.1:8051"
		fmt.Printf("Example: go run chatComPPLink.go %s %s\n", thisProcess, otherProcess)
		return
	}

	addresses := os.Args[1:]
	fmt.Println("Chat PPLink - addresses: ", addresses)

	lk := PP2PLink.NewPP2PLink(addresses[0], false)

	go func() {
		for {
			m := <-lk.Ind
			fmt.Println("Rcv: ", m)
		}
	}()

	go func() {
		sizes := [6]int{3, 6, 9, 13, 16, 18}
		primos3 := [nrPrimos]int{101, 883, 359, 941, 983, 859, 523, 631, 181, 233}
		primos6 := [nrPrimos]int{547369, 669437, 683251, 610279, 851117, 655439, 937351, 419443, 128467, 316879}
		primos9 := [nrPrimos]int{550032733, 429415309, 109543211, 882936113, 546857209, 756170741, 699422809, 469062577, 117355333, 617320027}
		primos13 := [nrPrimos]int{7069402558433, 960246047869, 5738081989711, 5358141480883, 2569391599009, 4135462531597, 7807787948171, 130788041233, 2708131414819, 1571981553097}
		primos16 := [nrPrimos]int{2207749090466833, 9361721528139247, 2657959759011013, 3551950148669023, 3460183118669741, 5503892014624961, 4067979800826917, 7848969908399551, 2806933754138389, 5211072635754109}
		primos18 := [nrPrimos]int{383376390724197361, 882611655919772761, 533290385325847007, 17969611178168479, 903013501582628521, 541906710014517121, 281512690206248899, 403936627075987639, 775148726422474717, 942319117335957539}

		todosPrimos := [len(sizes)][nrPrimos]int{primos3, primos6, primos9, primos13, primos16, primos18}

		for p := 0; p < len(sizes); p++ {
			messageBeginning := fmt.Sprintf("Enviando %d valores para %s. Magnitude = %d", nrPrimos, addresses[1], sizes[p])
			numbersStr := fmt.Sprint(todosPrimos[p])
			numbersFieldsArray := strings.Fields(numbersStr)
			numbersStr = strings.Join(numbersFieldsArray, ",")
			numbersStr = strings.Trim(numbersStr, "[]")
			message := messageBeginning + " || " + numbersStr
			fmt.Println("Snd: ", message)
			req := PP2PLink.PP2PLink_Req_Message{
				To:      addresses[1],
				Message: message}

			lk.Req <- req
		}
	}()

	<- make(chan struct{})
}

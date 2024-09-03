// por Fernando Dotti - fldotti.github.io - PUCRS - Escola Politécnica
// PROBLEMA:
//   o dorminhoco especificado no arquivo Ex1-ExplanacaoDoDorminhoco.pdf nesta pasta
// ESTE ARQUIVO
//   Um template para criar um anel generico.
//   Adapte para o problema do dorminhoco.
//   Nada está dito sobre como funciona a ordem de processos que batem.
//   O ultimo leva a rolhada ...
//   ESTE  PROGRAMA NAO FUNCIONA.    É UM RASCUNHO COM DICAS.

package main

import (
	"fmt"
	"math/rand"
	"time"
)

const NJ = 5 // número de jogadores
const M = 4  // número de cartas

type carta string // carta é uma string

var ch [NJ]chan carta // NJ canais de itens tipo carta
var bateu chan int
var proximaRodada chan bool // Canal para sinalizar o início de uma nova rodada

func jogador(id int, in chan carta, out chan carta, cartasIniciais []carta, inicioJogo chan bool) {
	mao := cartasIniciais // estado local - as cartas na mão do jogador
	nroDeCartas := M      // quantas cartas ele tem
	cartaRecebida := carta("")
	fmt.Printf("Jogador %d cartas na mão: %v\n", id, mao)

	<-inicioJogo // Aguarda o sinal para começar o jogo

	bateuAntes := false

	for {
		if (len(bateu) != 0 && !bateuAntes) { //se alguém bateu e não possui o id dentro do canal
			fmt.Printf("Jogador %d bateu também!\n", id)
			if(len(bateu) == NJ - 1) { //se todos bateram, o último leva a rolhada
				fmt.Printf("Jogador %d foi o último a bater e perdeu o jogo. Portanto, rolhada nele!!!\n", id)
			}
			bateu <- id
			return
		}

		if bateuAntes {
			time.Sleep(time.Millisecond * 100) // Aguarda um curto período antes de verificar novamente
			continue
		}

		if nroDeCartas == 5 {
			fmt.Println(id, " joga") // escreve seu identificador
			// escolhe alguma carta para passar adiante...
			indiceAleatorio := rand.Intn(nroDeCartas)
			cartaParaSair := mao[indiceAleatorio]
			
			if possuiQuatroLetrasIguais(mao) { // se possui 4 iguais passa a quinta diferente para o próximo
				cartaParaSair = cartaDiferente(mao)
			} 

			fmt.Printf("Jogador %d escolheu a carta %s\n", id, cartaParaSair)

			// manda carta escolhida para o próximo jogador
			out <- cartaParaSair
			mao = append(mao[:indiceAleatorio], mao[indiceAleatorio+1:]...)
			fmt.Printf("Minha nova mão (%d): %v\n", id, mao)
			nroDeCartas--

			if possuiQuatroLetrasIguais(mao) && nroDeCartas == 4 {
				fmt.Printf("Jogador %d bate!\n", id)
				bateuAntes = true
				bateu <- id
				proximaRodada <- false
			}
			
			// Verifica se todas as cartas foram jogadas
			if nroDeCartas == 0 {
				// Sinaliza o início de uma nova rodada
				proximaRodada <- true
			}
			time.Sleep(time.Millisecond * 100)
		} else {
			select {
			case cartaRecebida = <-in:
				fmt.Printf("Jogador %d recebeu a carta %s\n", id, cartaRecebida)
				mao = append(mao, cartaRecebida)
				nroDeCartas++
				fmt.Printf("Minha nova mão (%d): %v\n", id, mao)
				//time.Sleep(time.Millisecond * 100)
			default:
				// Não fazer nada, aguarda carta ser recebida
			}
		}
	}
}

func possuiQuatroLetrasIguais(mao []carta) bool {
	contagem := make(map[carta]int)
	for _, c := range mao {
		contagem[c]++
		if contagem[c] >= 4 {
			return true
		}
	}
	return false
}

func cartaDiferente(mao []carta) carta {
	contagem := make(map[carta]int)
	for _, c := range mao {
		contagem[c]++
	}

	// Percorre as cartas na mão
	for _, c := range mao {
		// Se a contagem da carta for igual a 1, significa que é a carta diferente
		if contagem[c] == 1 {
			return c
		}
	}

	// Se todas as cartas forem iguais, retorna uma carta vazia ou um valor de erro, dependendo do seu caso de uso
	return ""
}

func main() {
	rand.Seed(time.Now().UnixNano())
	bateu = make(chan int, 5)
	for i := 0; i < NJ; i++ {
		ch[i] = make(chan carta)
	}

	inicioJogo := make(chan bool)    // Canal para sinalizar o início do jogo
	proximaRodada = make(chan bool) // Canal para sinalizar o início de uma nova rodada

	for i := 0; i < NJ; i++ {
		cartasEscolhidas := make([]carta, 4)
		for j := 0; j < 4; j++ { // escolhe aleatoriamente (tira) cartas do baralho, passa cartas para jogador
			indiceAleatorio := rand.Intn(4) // Gera aleatoriamente um índice entre 0 e 3
			cartasEscolhidas[j] = carta('A' + indiceAleatorio)
		}
		go jogador(i, ch[i], ch[(i+1)%NJ], cartasEscolhidas, inicioJogo)
		time.Sleep(time.Millisecond * 100)
	}
	fmt.Println()

	time.Sleep(time.Millisecond * 100)

	// Inicie o jogo sinalizando para os jogadores que podem começar
	for i := 0; i < NJ; i++ {
		inicioJogo <- true
	}

	// Comece o jogo passando uma carta inicial
	cartaInicial := carta("@")
	ch[0] <- cartaInicial

	for i:=0 ; i < NJ; i++ {
		<-proximaRodada
	}
}
// por Fernando Dotti - PUCRS
// dado abaixo um exemplo de estrutura em arvore, uma arvore inicializada
// e uma operação de caminhamento, pede-se fazer:
//   1.a) a operação que soma todos elementos da arvore.
//        func soma(r *Nodo) int {...}
//   1.b) uma operação concorrente que soma todos elementos da arvore
//   2.a) a operação de busca de um elemento v, dizendo true se encontrou v na árvore, ou falso
//        func busca(r* Nodo, v int) bool {}...}
//   2.b) a operação de busca concorrente de um elemento, que informa imediatamente
//        por um canal se encontrou o elemento (sem acabar a busca), ou informa
//        que nao encontrou ao final da busca
//   3.a) a operação que escreve todos pares em um canal de saidaPares e
//        todos impares em um canal saidaImpares, e ao final avisa que acabou em um canal fin
//        func retornaParImpar(r *Nodo, saidaP chan int, saidaI chan int, fin chan struct{}){...}
//   3.b) a versao concorrente da operação acima, ou seja, os varios nodos sao testados
//        concorrentemente se pares ou impares, escrevendo o valor no canal adequado
//
//  ABAIXO: RESPOSTAS A QUESTOES 1a e b
//  APRESENTE A SOLUÇÃO PARA AS DEMAIS QUESTÕES

package main

import (
	"fmt"
)

type Nodo struct {
	v int
	e *Nodo
	d *Nodo
}

func caminhaERD(r *Nodo) {
	if r != nil {
		caminhaERD(r.e)
		fmt.Print(r.v, ", ")
		caminhaERD(r.d)
	}
}

// -------- SOMA ----------
// soma sequencial recursiva
func soma(r *Nodo) int {
	if r != nil {
		//fmt.Print(r.v, ", ")
		return r.v + soma(r.e) + soma(r.d)
	}
	return 0
}

// funcao "wraper" retorna valor
// internamente dispara recursao com somaConcCh
// usando canais
func somaConc(r *Nodo) int {
	s := make(chan int)
	go somaConcCh(r, s)
	return <-s
}
func somaConcCh(r *Nodo, s chan int) {
	if r != nil {
		s1 := make(chan int)
		go somaConcCh(r.e, s1)
		go somaConcCh(r.d, s1)
		s <- (r.v + <-s1 + <-s1)
	} else {
		s <- 0
	}
}

// 2.a) a operação de busca de um elemento v, dizendo true se encontrou v na árvore, ou falso
//
//	func busca(r* Nodo, v int) bool {}...}
func busca(r *Nodo, v int) bool {
	if r != nil {
		if r.v == v {
			return true
		}
		return busca(r.e, v) || busca(r.d, v)
	}
	return false
}

// 2.b) a operação de busca concorrente de um elemento, que informa imediatamente
//
//	por um canal se encontrou o elemento (sem acabar a busca), ou informa
//	que nao encontrou ao final da busca
func buscaConc(r *Nodo, v int) bool {
	s := make(chan bool)
	go buscaConcCh(r, v, s)
	return <-s
}

func buscaConcCh(r *Nodo, v int, s chan bool) {
	if r != nil {
		if r.v == v {
			s <- true
		}
		s1 := make(chan bool)
		go buscaConcCh(r.e, v, s1)
		go buscaConcCh(r.d, v, s1)
		s <- (<-s1 || <-s1)
	}
	s <- false
}

// 3.a) a operação que escreve todos pares em um canal de saidaPares e
//
//	todos impares em um canal saidaImpares, e ao final avisa que acabou em um canal fin
//	func retornaParImpar(r *Nodo, saidaP chan int, saidaI chan int, fin chan struct{}){...}
func retornaParImpar(r *Nodo, saidaP chan int, saidaI chan int, fin chan struct{}) {
	if r != nil {
		preencheSaida(r, saidaP, saidaI)
	}
	fin <- struct{}{}
}

func preencheSaida(r *Nodo, saidaP chan int, saidaI chan int) {
	if r != nil {
		preencheSaida(r.e, saidaP, saidaI)
		if r.v%2 == 0 {
			saidaP <- r.v
		} else {
			saidaI <- r.v
		}
		preencheSaida(r.d, saidaP, saidaI)
	}
	if r == nil {
		return
	}

}

// 3.b) a versao concorrente da operação acima, ou seja, os varios nodos sao testados
// concorrentemente se pares ou impares, escrevendo o valor no canal adequado
func retornaParImparConc(r *Nodo, saidaP chan int, saidaI chan int, fin chan struct{}) {
	terminou := make(chan bool)

	go func() {
		preencheSaidaConc(r, saidaP, saidaI, terminou)
		// Após a conclusão, sinal para indicar que terminou
		terminou <- true
	}()

	// Aguarda o sinal de conclusão
	<-terminou
	fin <- struct{}{}
}

func preencheSaidaConc(r *Nodo, saidaP chan int, saidaI chan int, terminou chan bool) {
	if r != nil {
		terminouE := make(chan bool)
		terminouD := make(chan bool)

		go func() {
			preencheSaidaConc(r.e, saidaP, saidaI, terminouE)
			terminouE <- true
		}()

		go func() {
			preencheSaidaConc(r.d, saidaP, saidaI, terminouD)
			terminouD <- true
		}()

		if r.v%2 == 0 {
			saidaP <- r.v
		} else {
			saidaI <- r.v
		}

		// Aguarda os sinais de conclusão das threads filhas
		<-terminouE
		<-terminouD
	}
	// Sinal para indicar que terminou
	terminou <- true
}

func main() {
	root := &Nodo{v: 10, // raiz
		e: &Nodo{v: 5, // esquerda de 10
			e: &Nodo{v: 3, // esquerda de 5
				e: &Nodo{v: 1, e: nil, d: nil},  // esquerda de 3
				d: &Nodo{v: 4, e: nil, d: nil}}, // direita de 3
			d: &Nodo{v: 7, // direita de 5
				e: &Nodo{v: 6, e: nil, d: nil},   // esquerda de 7
				d: &Nodo{v: 8, e: nil, d: nil}}}, // direita de 7
		d: &Nodo{v: 15, // direita de 10
			e: &Nodo{v: 13, // esquerda de 15
				e: &Nodo{v: 12, e: nil, d: nil},  // esquerda de 13
				d: &Nodo{v: 14, e: nil, d: nil}}, // direita de 13
			d: &Nodo{v: 18, // direita de 18
				e: &Nodo{v: 17, e: nil, d: nil},   // esquerda de 18
				d: &Nodo{v: 19, e: nil, d: nil}}}} // direita de 18

	fmt.Println()
	fmt.Print("Valores na árvore: ")
	caminhaERD(root)
	fmt.Println()
	fmt.Println()

	fmt.Println("Soma: ", soma(root))
	fmt.Println("SomaConc: ", somaConc(root))
	fmt.Println()
	fmt.Println()

	fmt.Println("Busca: ", busca(root, 13))
	fmt.Println("Busca: ", busca(root, 10))
	fmt.Println("Busca: ", busca(root, 20))
	fmt.Println("Busca: ", busca(root, 0))
	fmt.Println()
	fmt.Println()

	fmt.Println("BuscaConc: ", buscaConc(root, 13))
	fmt.Println("BuscaConc: ", buscaConc(root, 10))
	fmt.Println("BuscaConc: ", buscaConc(root, 20))
	fmt.Println("BuscaConc: ", buscaConc(root, 0))
	fmt.Println()
	fmt.Println()

	fmt.Println("Retorna Par e Impar: ")
	saidaP := make(chan int)
	saidaI := make(chan int)
	fin := make(chan struct{})
	go retornaParImpar(root, saidaP, saidaI, fin)
	//go retornaParImparConc(root, saidaP, saidaI, fin)
	/*
		caso queira testar o concorrente, basta comentar o
		anterior e descomentar o concorrente.
	*/
	for {
		select {
		case p := <-saidaP:
			fmt.Println("Par: ", p)
		case i := <-saidaI:
			fmt.Println("Impar: ", i)
		case <-fin:
			fmt.Println("Fim")
			return
		}
	}

}

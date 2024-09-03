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

// a operação de busca de um elemento v, dizendo true se encontrou v na árvore, ou falso
//
//	func busca(r* Nodo, v int) bool {}...}
func busca(r *Nodo, v int) bool {
	if r != nil {
		if r.v == v {
			return true
		}
		return busca(r.d, v) || busca(r.e, v)
	}
	return false
}

// 3.a) a operação que escreve todos pares em um canal de saidaPares e
//
//	todos impares em um canal saidaImpares, e ao final avisa que acabou em um canal fin
//	func retornaParImpar(r *Nodo, saidaP chan int, saidaI chan int, fin chan struct{}){...}
func retornaParImpar(r *Nodo, saidaP chan int, saidaI chan int, fin chan struct{}) {
	fmt.Println("Nodo: ", r.v)
	if r != nil {
		retornaParImpar(r.e, saidaP, saidaI, fin)
		retornaParImpar(r.d, saidaP, saidaI, fin)
		if r.v%2 == 0 {
			fmt.Println("Par: ", r.v)
			saidaP <- r.v

		} else {
			fmt.Println("Impar: ", r.v)
			saidaI <- r.v
		}
	} else {
		fmt.Println("É nulo")
		fin <- struct{}{}
	}
}

// 3.b) a versao concorrente da operação acima, ou seja, os varios nodos sao testados
//
//	concorrentemente se pares ou impares, escrevendo o valor no canal adequado
func retornaParImparConc(r *Nodo) int {
	saidaP := make(chan int)
	saidaI := make(chan int)
	fin := make(chan struct{})
	go retornaParImparConcCh(r, saidaP, saidaI, fin)
	fmt.Println("Aguardando fim...")
	return <-saidaP
}
func retornaParImparConcCh(r *Nodo, saidaP chan int, saidaI chan int, fin chan struct{}) {
	if r != nil {
		fmt.Println("Nodo: ", r.v)
		p1 := make(chan int)
		i1 := make(chan int)
		f1 := make(chan struct{})
		go retornaParImparConcCh(r.e, p1, i1, f1)
		go retornaParImparConcCh(r.d, p1, i1, f1)
		if r.v%2 == 0 {
			saidaP <- r.v
			fmt.Println("Par: ", r.v)
		} else {
			saidaI <- r.v
			fmt.Println("Impar: ", r.v)
		}
	} else {
		fin <- struct{}{}

	}

}

// ---------   agora vamos criar a arvore e usar as funcoes acima

func main() {
	root := &Nodo{v: 10,
		e: &Nodo{v: 5,
			e: &Nodo{v: 3,
				e: &Nodo{v: 1, e: nil, d: nil},
				d: &Nodo{v: 4, e: nil, d: nil}},
			d: &Nodo{v: 7,
				e: &Nodo{v: 6, e: nil, d: nil},
				d: &Nodo{v: 8, e: nil, d: nil}}},
		d: &Nodo{v: 15,
			e: &Nodo{v: 13,
				e: &Nodo{v: 12, e: nil, d: nil},
				d: &Nodo{v: 14, e: nil, d: nil}},
			d: &Nodo{v: 18,
				e: &Nodo{v: 17, e: nil, d: nil},
				d: &Nodo{v: 19, e: nil, d: nil}}}}

	fmt.Println()
	fmt.Print("Valores na árvore: ")
	caminhaERD(root)
	fmt.Println()
	fmt.Println()

	fmt.Println("Soma: ", soma(root))
	fmt.Println("SomaConc: ", somaConc(root))
	fmt.Println()
	fmt.Println("Contém 13?: ", busca(root, 13))
	fmt.Println("Contém 20?: ", busca(root, 20))
	fmt.Println()
	saidaP := make(chan int, 10)
	saidaI := make(chan int, 10)
	fin := make(chan struct{})
	retornaParImpar(root, saidaP, saidaI, fin)

	for i := 0; i < 21; i++ {
		<-fin
	}

	close(saidaP)
	close(saidaI)
	close(fin)

	fmt.Println("Pares e ímpares processados com sucesso.")

}

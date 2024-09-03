// por Fernando Dotti - fldotti.github.io - PUCRS - Escola Politécnica
// servidor com criacao dinamica de thread de servico
// Problema:
//   considere um servidor que recebe pedidos por um canal (representando uma conexao)
//   ao receber o pedido, sabe-se através de qual canal (conexao) responder ao cliente.
//   Abaixo uma solucao sequencial para o servidor.
// Exercicio
//   deseja-se tratar os clientes concorrentemente, e nao sequencialmente.
//   como ficaria a solucao ?
// Veja abaixo a resposta ...
//   quantos clientes podem estar sendo tratados concorrentemente ?
//   resposta: Não há um limite definido para o número de clientes que podem ser tratados concorrentemente,
//	 o número de goroutines concorrentes é dinâmico e pode aumentar à medida que mais solicitações de clientes chegam
//
// Exercicio:
//   agora suponha que o seu servidor pode estar tratando no maximo 10 clientes concorrentemente.
//   como voce faria ?
//  resposta: Usando um canal bufferizado com tamanho 10 (código)
//

package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

const (
	NCL  = 100
	Pool = 10
)

type Request struct {
	v      int
	ch_ret chan int
}

// ------------------------------------
// cliente
func cliente(i int, req chan Request, sem chan struct{}) {
	var v, r int
	my_ch := make(chan int)
	for {
		v = rand.Intn(1000)
		req <- Request{v, my_ch}
		r = <-my_ch
		fmt.Println("cli: ", i, " req: ", v, "  resp:", r, "Processo: ", strconv.Itoa(len(sem)))
		// Liberar o semáforo após a conclusão da tarefa
		<-sem
		time.Sleep(60 * time.Second)
	}
}

// ------------------------------------
// servidor
// thread de servico calcula a resposta e manda direto pelo canal de retorno informado pelo cliente
func trataReq(id int, req Request) {
	fmt.Println("                                 trataReq ", id)
	req.ch_ret <- req.v * 2
	//esperar tempo aleatorio 
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
}

// servidor que dispara threads de servico
func servidorConc(in chan Request, sem chan struct{}) {
	// servidor fica em loop eterno recebendo pedidos e criando um processo concorrente para tratar cada pedido
	var j int = 0
	for {
		j++
		req := <-in
		sem <- struct{}{}
		go trataReq(j, req)
	}
}

// ------------------------------------
// main
func main() {
	fmt.Println("------ Servidores - criacao dinamica -------")
	serv_chan := make(chan Request, Pool) // CANAL POR ONDE SERVIDOR RECEBE PEDIDOS
	sem := make(chan struct{}, Pool) // Use um canal bufferizado com tamanho 10
	go servidorConc(serv_chan, sem)      // LANÇA PROCESSO SERVIDOR
	for i := 0; i < NCL; i++ {      // LANÇA DIVERSOS CLIENTES
		go cliente(i, serv_chan, sem)
	}
	<-make(chan int)
}
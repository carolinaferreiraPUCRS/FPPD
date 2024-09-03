package main

import (
    "fmt"
    "time"
)

var (
    semPCR1 = NewSemaphore(10)
    semPCR2 = NewSemaphore(5)
)

type Semaphore struct {
    v    int           // valor do semaforo: negativo significa proc bloqueado
    fila chan struct{} // canal para bloquear os processos se v < 0
    sc   chan struct{} // canal para atomicidade das operacoes wait e signal
}

func NewSemaphore(init int) *Semaphore {
    s := &Semaphore{
        v:    init,                   // valor inicial de creditos
        fila: make(chan struct{}),    // canal sincrono para bloquear processos
        sc:   make(chan struct{}, 1), // usaremos este como semaforo para SC, somente 0 ou 1
    }
    return s
}

func (s *Semaphore) Wait() {
    fmt.Println("W_ Haviam ", s.v, " creditos")
    s.sc <- struct{}{} // SC do semaforo feita com canal
    s.v--              // decrementa valor
    if s.v < 0 {       // se negativo era 0 ou menor, tem que bloquear
        <-s.sc               // antes de bloq, libera acesso
        s.fila <- struct{}{} // bloqueia proc
    } else {
        <-s.sc // libera acesso
    }
    fmt.Println("W_ Agora tem ", s.v, " creditos")
}

func (s *Semaphore) Signal() {
    fmt.Println("S_ Haviam ", s.v, " creditos")
    s.sc <- struct{}{} // entra sc
    s.v++
    if s.v <= 0 { // tem processo bloqueado ?
        <-s.fila // desbloqueia
    }
    <-s.sc // libera SC para outra op
    fmt.Println("Agora tem ", s.v, " creditos")
}

func thread(id int) {
    for {
        fmt.Println(id, "thread is threading in its thread")

        semPCR1.Wait()
            fmt.Println(id, "thread is threading in its PCR1 thread")
        semPCR1.Signal()

        semPCR2.Wait()
            fmt.Println(id, "thread is threading in its PCR2 thread")
        semPCR2.Signal()
    }
}

func main() {
    for i := 0; i < 10; i++ {
        go thread(i)
    }

    <- time.After(1000 * time.Millisecond)
}
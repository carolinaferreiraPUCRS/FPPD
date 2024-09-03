package main

import (
    "fmt"
    "time"
)

type ma struct{}

var (
    chanPCR1 = make(chan ma, 10)
    amountPCR1 = 0
    chanPCR2 = make(chan ma, 5)
    amountPCR2 = 0
    mutexPCR1 = make(chan ma, 1)
    mutexPCR2 = make(chan ma, 1)
)

func thread(id int) {
    for {
        fmt.Println(id, "thread is threading in its thread")

        chanPCR1 <- ma{}
            mutexPCR1 <- ma{}
                amountPCR1++
                fmt.Printf("thread %d is threading in its PCR1 thread (%d/%d)\n", id, amountPCR1, 10)
            <- mutexPCR1
        <- chanPCR1

        mutexPCR1 <- ma{}
            amountPCR1--
        <- mutexPCR1

        chanPCR2 <- ma{}
            mutexPCR2 <- ma{}
                amountPCR2++
                fmt.Printf("thread %d is threading in its PCR2 thread (%d/%d)\n", id, amountPCR2, 5)
            <- mutexPCR2
        <- chanPCR2

        mutexPCR2 <- ma{}
            amountPCR2--
        <- mutexPCR2
    }
}

func main() {
    for i := 0; i < 10; i++ {
        go thread(i)
    }

    <- time.After(1000 * time.Millisecond)
}
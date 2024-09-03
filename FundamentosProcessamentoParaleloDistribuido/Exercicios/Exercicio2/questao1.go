
package main

import (
"fmt"
)

var sharedTest int = 0 // variavel compartilhada
var ch_fim chan struct{} = make(chan struct{})
func MyFunc(inc int) {
 for k := 0; k < 1000; k++ {
 sharedTest = sharedTest + inc
 }
 ch_fim <- struct{}{}
}
func main() {
 for i := 0; i < 10; i++ {
 go MyFunc(1)
 go MyFunc(-1)
 }
 fmt.Println("Criei 20 processos")
 for i := 0; i < 20; i++ {
 <-ch_fim
 }
 fmt.Println("Processos acabaram. Resultado ", sharedTest)
}
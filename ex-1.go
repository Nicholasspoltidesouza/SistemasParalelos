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
	tamanhoDosNumeros = 6  // numero de magnitudes dos valores primos
	primos            = 10 // numero de valores primos para cada magnitude
)

func contaPrimosSeq(s [primos]int) time.Duration {
	start := time.Now()
	for i := 0; i < primos; i++ {
		if isPrime(s[i]) {
		}
	}
	return time.Since(start)
}

func contaPrimosConc(s [primos]int, end chan int) time.Duration {
	start := time.Now()
	for i := 0; i < primos; i++ {
		go func(i int) {
			if isPrime(s[i]) {
			}
			end <- 1
		}(i)
	}
	for i := 0; i < primos; i++ {
		<-end
	}
	return time.Since(start)
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

func main() {
	var n int
	n = 2 // numero de processadores
		runtime.GOMAXPROCS(n) // usando n processadores

	//  valores primos com respectivamente 3, 6, 9, 13, 18 casas
	//  use o programa AchaNPrimos para achar primos com determinado número de casas

	primos3 := [primos]int{101, 883, 359, 941, 983, 859, 523, 631, 181, 233}
	primos6 := [primos]int{547369, 669437, 683251, 610279, 851117, 655439, 937351, 419443, 128467, 316879}
	primos9 := [primos]int{550032733, 429415309, 109543211, 882936113, 546857209, 756170741, 699422809, 469062577, 117355333, 617320027}
	primos13 := [primos]int{7069402558433, 960246047869, 5738081989711, 5358141480883, 2569391599009, 4135462531597, 7807787948171, 130788041233, 2708131414819, 1571981553097}
	primos16 := [primos]int{2207749090466833, 9361721528139247, 2657959759011013, 3551950148669023, 3460183118669741, 5503892014624961, 4067979800826917, 7848969908399551, 2806933754138389, 5211072635754109}
	primos18 := [primos]int{383376390724197361, 882611655919772761, 533290385325847007, 17969611178168479, 903013501582628521, 541906710014517121, 281512690206248899, 403936627075987639, 775148726422474717, 942319117335957539}

	numerosPrimos := [tamanhoDosNumeros][primos]int{primos3, primos6, primos9, primos13, primos16, primos18}

	for n := 0; n < tamanhoDosNumeros; n++ {
		fmt.Println("Valores: ", numerosPrimos[n])
		resultado := contaPrimosSeq(numerosPrimos[n])
		fmt.Println(" ")
		fmt.Println("\nTempo Sequencial: ", resultado)
		fim := make(chan int)
		resultado = contaPrimosConc(numerosPrimos[n], fim)
		fmt.Println("\nTempo Concorrente: ", resultado)
	}
}

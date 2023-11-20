package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

func contaPrimosSeq(slice []int) int {
	var n int
	for _, v := range slice {
		if ehPrimo(v) {
			n++
		}
	}
	return n
}

func contaPrimosConc(slice []int, numProcessadores int) int {
	var n int
	c := make(chan int)

	for _, v := range slice {
		go func(v int) {
			if ehPrimo(v) {
				c <- 1
			} else {
				c <- 0
			}
		}(v)
	}

	for i := 0; i < len(slice); i++ {
		n += <-c
	}

	close(c)
	return n
}

func ehPrimo(n int) bool {
	if n <= 1 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func generateSlice(size int) []int {
	slice := make([]int, size)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		slice[i] = rand.Intn(2*size) + 1
	}
	return slice
}

func main() {
	const tamanho = 2000000
	processadores := []int{1, 2, 4}

	for _, numProcessadores := range processadores {
		fmt.Printf("------ conta primos de tamanho %d com %d processadores -------\n", tamanho, numProcessadores)
		slice := generateSlice(tamanho)

		// Sequencial
		start := time.Now()
		p := contaPrimosSeq(slice)
		fmt.Println(" -> sequencial ------ secs: ", time.Since(start).Seconds())
		fmt.Println(" ------ n primos : ", p)

		// Concorrente
		start1 := time.Now()
		runtime.GOMAXPROCS(numProcessadores)
		p = contaPrimosConc(slice, numProcessadores)
		fmt.Println(" -> concorrente ------ secs: ", time.Since(start1).Seconds())
		fmt.Println(" ------ n primos : ", p)
		fmt.Println()
	}
}

package main

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"
)

func gerarVetorAleatorio(tamanho, min, max int) []int {
	rand.Seed(time.Now().UnixNano())
	vetor := make([]int, tamanho)
	for i := 0; i < tamanho; i++ {
		vetor[i] = rand.Intn(max-min+1) + min
	}
	return vetor
}

func ordenacaoSequencial(array []int, tamanhoBaldes int) []int {
	var max, min int
	for _, n := range array {
		if n < min {
			min = n
		}
		if n > max {
			max = n
		}
	}
	nBaldes := int(max-min)/tamanhoBaldes + 1
	baldes := make([][]int, nBaldes)
	for i := 0; i < nBaldes; i++ {
		baldes[i] = make([]int, 0)
	}

	for _, n := range array {
		idx := int(n-min) / tamanhoBaldes
		baldes[idx] = append(baldes[idx], n)
	}

	ordenado := make([]int, 0)
	for _, balde := range baldes {
		if len(balde) > 0 {
			insertionSort(balde)
			ordenado = append(ordenado, balde...)
		}
	}

	return ordenado
}

func insertionSort(array []int) {
	for i := 0; i < len(array); i++ {
		temp := array[i]
		j := i - 1
		for ; j >= 0 && array[j] > temp; j-- {
			array[j+1] = array[j]
		}
		array[j+1] = temp
	}
}

func ordenacaoParalela(vetor []int, nbaldes int, wg *sync.WaitGroup) {
	defer wg.Done()

	baldes := make([][]int, nbaldes)
	minValor := vetor[0]
	maxValor := vetor[0]

	// Encontrar os valores mínimo e máximo no vetor
	for _, valor := range vetor {
		if valor < minValor {
			minValor = valor
		}
		if valor > maxValor {
			maxValor = valor
		}
	}

	// Determinar o intervalo de cada balde
	intervaloBalde := (maxValor - minValor + 1) / nbaldes

	// Distribuir valores nos baldes
	for _, valor := range vetor {
		indiceBalde := (valor - minValor) / intervaloBalde
		if indiceBalde == nbaldes {
			indiceBalde--
		}
		baldes[indiceBalde] = append(baldes[indiceBalde], valor)
	}

	// Ordenar cada balde
	for i := 0; i < nbaldes; i++ {
		sort.Ints(baldes[i])
	}

	// Concatenar baldes ordenados
	resultado := make([]int, 0, len(vetor))
	for _, balde := range baldes {
		resultado = append(resultado, balde...)
	}

	copy(vetor, resultado)
}

func main() {
	tamanho := 1000000      // Tamanho do vetor
	minValor := 0           // Valor mínimo no vetor
	maxValor := 999999999    // Valor máximo no vetor

	vetorSequencial := gerarVetorAleatorio(tamanho, minValor, maxValor)
	vetorParalelo := make([]int, tamanho)
	copy(vetorParalelo, vetorSequencial)

	nbaldesLista := []int{2, 4, 6, 8, 16}

	fmt.Printf("Tamanho do vetor: %d\n", tamanho)

	// Execução sequencial
	tempoInicioSequencial := time.Now()
	ordenacaoSequencial(vetorSequencial, tamanho)
	tempoSequencial := time.Since(tempoInicioSequencial)
	fmt.Printf("Tempo sequencial: %s\n", tempoSequencial)

	// Execução paralela variando o número de baldes
	for _, nbaldes := range nbaldesLista {
		fmt.Printf("\nNúmero de Baldes (NB): %d\n", nbaldes)

		tempoInicioParalelo := time.Now()
		var wg sync.WaitGroup
		wg.Add(nbaldes)

		// Divide o vetor em partes e classifica cada parte em uma goroutine separada
		for i := 0; i < nbaldes; i++ {
			inicio := i * (tamanho / nbaldes)
			fim := (i + 1) * (tamanho / nbaldes)
			if i == nbaldes-1 {
				fim = tamanho
			}
			go ordenacaoParalela(vetorParalelo[inicio:fim], nbaldes, &wg)
		}

		wg.Wait()

		// Ordena novamente o vetor final resultante da combinação dos baldes
		sort.Ints(vetorParalelo)

		tempoParalelo := time.Since(tempoInicioParalelo)
		fmt.Printf("Tempo paralelo: %s\n", tempoParalelo)

		// Calcula o speedup em relação à versão sequencial
		speedup := float64(tempoSequencial) / float64(tempoParalelo)
		fmt.Printf("Speedup: %.2f\n", speedup)
	}
}

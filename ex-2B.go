package main

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"
)

func gerarVetorAleatorio(tamanho, min, max int) []int {
	rand.Seed(time.Now().UnixNano()) // semente para geração de números aleatórios
	vetor := make([]int, tamanho) // cria o vetor
	for i := 0; i < tamanho; i++ { // preenche o vetor com números aleatórios
		vetor[i] = rand.Intn(max-min+1) + min
	}
	return vetor
}

func ordenacaoSequencial(array []int, tamanhoBaldes int) []int {
	var max, min int // valor máximo e mínimo no vetor
	for _, n := range array { // encontra o valor máximo e mínimo no vetor
		if n < min { 
			min = n	
		}
		if n > max {
			max = n
		}
	}
	nBaldes := int(max-min)/tamanhoBaldes + 1 // número de baldes
	baldes := make([][]int, nBaldes) // cria os baldes
	for i := 0; i < nBaldes; i++ { // inicializa os baldes
		baldes[i] = make([]int, 0) 
	}

	for _, n := range array { // distribui os valores nos baldes
		idx := int(n-min) / tamanhoBaldes
		baldes[idx] = append(baldes[idx], n) 
	}

	ordenado := make([]int, 0) // vetor ordenado
	for _, balde := range baldes {
		if len(balde) > 0 { // ordena cada balde e concatena no vetor ordenado
			sort.Ints(balde) // ordena o balde
			ordenado = append(ordenado, balde...) // concatena o balde no vetor ordenado
		}
	}

	return ordenado // retorna o vetor ordenado
}

func ordenacaoParalela(vetor []int, nbaldes int, wg *sync.WaitGroup) {
	defer wg.Done() // decrementa o contador de goroutines

	baldes := make([][]int, nbaldes) // cria os baldes
	minValor := vetor[0] // valor mínimo no vetor
	maxValor := vetor[0] // valor máximo no vetor

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
		if indiceBalde == nbaldes { // se o valor for igual ao máximo, coloca no último balde
			indiceBalde--
		}
		baldes[indiceBalde] = append(baldes[indiceBalde], valor) // adiciona o valor no balde
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

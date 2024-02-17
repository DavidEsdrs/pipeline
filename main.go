package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

func readFromFolder(folderName string) chan string {
	out := make(chan string, 8)

	go func() {
		entries, err := os.ReadDir(folderName) // lê o diretório e retorna todos os "entries" (arquivos/pastas)

		if err != nil { // interrompe a execução imediatamente ao encontrar um erro
			panic(err)
		}

		for _, e := range entries { // para cada arquivo encontrado...
			out <- "files/" + e.Name() // ...envie o nome para o canal de saída
		}

		close(out)
	}()

	return out
}

func readFromFile(filesNames chan string) chan string {
	out := make(chan string, 8)
	var wg sync.WaitGroup

	process := func(fname string) {
		defer wg.Done()
		f, err := os.Open(fname)

		if err != nil {
			panic(err)
		}

		scanner := bufio.NewScanner(f)

		for scanner.Scan() {
			line := scanner.Text()
			out <- line
		}
	}

	go func() {
		for f := range filesNames { // enquanto o canal estiver aberto...
			wg.Add(1)
			go process(f) // lançamos uma goroutine para cada arquivo
		}

		wg.Wait()

		close(out)
	}()

	return out
}

func processLine(lines chan string) chan int {
	out := make(chan int, 8)

	go func() {
		for l := range lines {
			strgs := strings.Split(l, " ")
			out <- len(strgs)
		}

		close(out)
	}()

	return out
}

func countFromLine(in chan int) int64 {
	result := int64(0)

	for n := range in {
		atomic.AddInt64(&result, int64(n))
	}

	return result
}

func Execute() {
	start := time.Now()

	readOutput := readFromFolder("files")
	fileOutput := readFromFile(readOutput)
	countOutput := processLine(fileOutput)
	result := countFromLine(countOutput)

	duration := time.Since(start)

	fmt.Printf("%v - it took %vms", result, duration.Milliseconds())
}

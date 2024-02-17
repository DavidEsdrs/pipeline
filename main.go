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
	out := make(chan string, 10)

	go func() { // lançamos uma goroutine que lançará o resultado da operação no canal out
		entries, err := os.ReadDir(folderName) // lê o diretório e retorna todos as "entries" (arquivos/pastas)

		if err != nil { // interrompe a execução imediatamente ao encontrar um erro
			panic(err)
		}

		for _, e := range entries { // para cada arquivo encontrado...
			out <- "files/" + e.Name() // ...envie o nome para o canal de saída
		}

		close(out) // feche o canal ao fim do processo
	}()

	return out
}

// o canal de input daqui é o canal de output da função anterior
func readFile(filesNames chan string) chan *os.File {
	out := make(chan *os.File, 10) // iremos retornar os arquivos lidos por esse canal

	go func() {
		// mesmo padrão de antes. Aqui iremos consumir o que chegar do canal de input
		for fname := range filesNames {
			f, err := os.Open(fname)

			if err != nil {
				panic(err)
			}

			out <- f
		}

		close(out)
	}()

	return out // retornamos imediatamente. Aqui será enviado o conteúdo dos arquivos linha a linha
}

func readFileContent(files chan *os.File) chan string {
	out := make(chan string, 1024)
	var wg sync.WaitGroup

	// podemos definir a função "process" dentro da própria função para termos um
	// código mais legível (assim também evitamos ter que passar mais argumentos para essa função)
	process := func(f *os.File) {
		defer wg.Done() // ao fim da função, indicaremos ao WaitGroup que cocluímos a execução dessa goroutine
		defer f.Close() // e fecharemos o arquivo para economizar recursos do sistema

		scanner := bufio.NewScanner(f) // o scanner permite a leitura do arquivo linha a linha

		for scanner.Scan() { // enquanto houver linha para ler...
			line := scanner.Text() // ...guardamos seu conteúdo dentro de "line"...
			out <- line            // ...e enviamos para o canal de saída
		}
	}

	go func() {
		for f := range files { // enquanto o canal "filesName" estiver aberto...
			wg.Add(1)
			go process(f) // lançamos uma goroutine para cada arquivo (ainda iremos definir a função process)
		}

		wg.Wait() // esperamos até que todas as goroutines lançadas concluam seu processo...

		close(out) // ...e então fechamos o canal
	}()

	return out
}

func countWords(lines chan string) chan int {
	out := make(chan int, 1024)

	go func() {
		var wg sync.WaitGroup

		for l := range lines { // enquanto o canal lines estiver aberto...
			wg.Add(1)

			go func(l string) {
				defer wg.Done()
				strgs := strings.Split(l, " ") // ...toda string "l" que chegar de lines, dividimos pelos espaços...
				out <- len(strgs)              // ...e enviamos o seu tamanho para o canal de output
			}(l)
		}

		wg.Wait()
		close(out)
	}()

	return out
}

func countFromLine(in chan int) int64 {
	result := int64(0)

	for n := range in { // enquanto o canal de input estiver aberto, continuamos a receber os valores

		// adicionamos o valor a result atomicamente para garantir que não haja acessos
		// concorrentes à result
		atomic.AddInt64(&result, int64(n))
	}

	return result
}

func main() {
	start := time.Now()

	readOutput := readFromFolder("files")
	filesOutput := readFile(readOutput)
	contentOutput := readFileContent(filesOutput)
	countOutput := countWords(contentOutput)
	result := countFromLine(countOutput)

	duration := time.Since(start)

	fmt.Printf("%v palavras contadas em %vms", result, duration.Milliseconds())
}

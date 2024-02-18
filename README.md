# Pipeline de dados em Golang

Esse é um exemplo completo de uma pipeline em Golang.

## Qual o objetivo?

Entender padrões e design para utilização de goroutines, canais, WaitGroup e atomic em Golang.

## O que foi feito?

Quando executado, todos os arquivos de texto da pasta "files" serão lidos 
concorrentemente (e, muito provavelmente, em paralelo) e serão contados o 
número de palavras de cada arquivo. No fim, será printado o número de palavras
contadas e quanto tempo levou para contar.

## Como foi feito?

O padrão de pipeline de dados foi aplicado. A pipeline consiste em três etapas
principais e uma adicional:

1° `readFromFolder`:
Lemos a pasta que contém os arquivos que iremos processar.

2° `readFile`:
Abrimos os arquivos e enviamos para a próxima etapa por meio de um canal.

3° `readFileContent`:
Consumimos o contéudo do canal anterior e lemos o seu conteúdo, linha a linha.

4° `processLine`:
Dividimos as linhas pelas palavras e contamos a quantidade.

5° `countFromLine`:
Agregamos os resultados que vem da função anterior e somamos ao resultado, que é
retornado ao fim da função.

6° `main`:
Compomos as etapas e calculamos o tempo que levou todo o processo.

## Como executar

Clone o respositório em uma pasta local. Você pode utilizar a build que está disponível na pasta `builds` para um resultado mais realista:
```sh
// no windows powershell
.\builds\pipeline.exe

// num terminal linux
chmod +x pipeline
./pipeline
```
Você pode executar a build de desenvolvimento, para isso, você deve ter a [runtime do Go instalada na sua máquina](https://go.dev/dl/) (o projeto foi desenvolvido utilzando go v1.21). Após instalar, execute o código principal:
```sh
go run main.go
```
Em ambos os casos, você obterá um resultado semelhante a esse:
```sh
30102 palavras contadas em 2ms
```

## Considerações

- Essa pipeline NÃO necessariamente é mais performática que uma versão equivalente
síncrona. Goroutines brilham em cenários em que mais dados são processados!
- Abusamos um pouco dos canais e goroutines no processo. Por exemplo, a primeira
função `readFromFolder` não precisaria ser uma etapa em si da pipeline (utilzando
o canal para saída), ele poderia facilmente ser síncrono.

## Desafio

- Desafio 1: Crie uma nova etapa (ou mude uma das etapas, o que achar melhor) para contar o
número de caracteres dos arquivos. A pipeline, ao fim do processo deve retornar
uma struct `Result`:
```go
type Result struct {
  CharCount int64 // quantidade de caracteres
  WordCount int64 // quantidade de palavras
}
```
Faça as contagens de forma concorrente entre si.

- Crie o código síncrono equivalente à essa pipeline e faça o benchmarking de ambas
as versões utilizando uma ferramente como o [hyperfine](https://github.com/sharkdp/hyperfine). Teste o resultado com diferentes quantidades de arquivos, para isso, você pode colocar, por exemplo, 100 arquivos dentro da pasta `files` e criar um contador para controlar a quantidade máxima de arquivos que devem ser lidos.
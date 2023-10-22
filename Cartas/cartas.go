package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2"
)

type Fila struct {
	elementos []string
}

func nova_fila() *Fila {
	return &Fila{
		elementos: make([]string, 0),
	}
}
func (f *Fila) enfileirar(item string) {
	f.elementos = append(f.elementos, item)
}
func (f *Fila) desenfileirar() (string, bool) {
	if f.fila_vazia() {
		return "", false
	}
	item := f.elementos[0]
	f.elementos = f.elementos[1:]
	return item, true
}
func (f *Fila) fila_vazia() bool {
	return len(f.elementos) == 0
}

var mu sync.Mutex

// variavel que iremos utilizar como a senha correta gerada
var senha_correta string

// Função capaz de gerar senha e colocá-la na fila
func GerarSenha() {
	mu.Lock()
	senha_aleatoria := rand.Intn(10) + 1
	senha_correta = strconv.Itoa(senha_aleatoria)
	mu.Unlock()
}

func main() {
	var cartas_retiradas = nova_fila()
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		c.Type("html", "utf-8") // a codificação UTF-8 é para consertar problemas de texto
		// Abrindo o arquivo HTML da pasta
		file, err := os.Open("static/inicio.html")
		if err != nil {
			// Lidar com erros, como arquivo não encontrado
			return c.SendString("Erro ao carregar o arquivo HTML")
		}
		defer file.Close()

		// Lendo o conteúdo do arquivo
		htmlContent, err := io.ReadAll(file)
		if err != nil {
			// Lidar com erros de leitura
			return c.SendString("Erro ao ler o arquivo HTML")
		}

		cartasNaFila := []string{}
		for !cartas_retiradas.fila_vazia() {
			item, _ := cartas_retiradas.desenfileirar()
			cartasNaFila = append(cartasNaFila, "<li>"+item+"</li>")
		}

		//senhaCorreta é adicionada ao %s e html enviadr como resosta
		return c.SendString(fmt.Sprintf(string(htmlContent), senha_correta, strings.Join(cartasNaFila, "\n")))
	})

	//rota para a ação de gerar nova senha
	app.Post("/gerar_senha", func(c *fiber.Ctx) error {
		go GerarSenha()
		mu.Lock()
		defer mu.Unlock()
		// redireciona a raiz após a geração
		return c.Redirect("/")
	})
	//rota para confirmar senha
	app.Post("/confirmar_senha", func(c *fiber.Ctx) error {
		//obtem senha digitada pelo usuário e faz a comparação com a correta
		senha_digitada := c.FormValue("senha_usuario")
		//retorna ou para a pagina onde mostra a carta ou para uma mensagem de erro dependendo da comparação
		c.Type("html", "utf-8")
		// Abrindo o arquivo HTML da pasta
		file, err := os.Open("static/confirma-senha.html")
		if err != nil {
			// Lidar com erros, como arquivo não encontrado
			return c.SendString("Erro ao carregar o arquivo HTML")
		}
		defer file.Close()

		// Lendo o conteúdo do arquivo
		htmlContent, err := io.ReadAll(file)
		if err != nil {
			// Lidar com erros de leitura
			return c.SendString("Erro ao ler o arquivo HTML")
		}
		// conversão para inteiro, utlizamos '_' para ignorar o segundo valor de retorno ( variável do tipo error que indica se a conversão foi bem-sucedida, podemos implementar futuramenete)
		index, _ := strconv.Atoi(senha_correta)
		// indices começam em 0, por isso a subtração
		index--
		//lista de cartas possiveis
		cartas := []string{"As de Copas", "2 de Copas", "8 de Espadas", "As de Ouro", "K de Paus", "6 de Espadas", "J de Copas", "3 de Ouro", "As de Paus", "J de Ouro"}
		// carta recebe a carta sorteada
		carta := cartas[index]
		if senha_digitada == senha_correta {
			cartas_retiradas.enfileirar(carta)
			return c.SendString(fmt.Sprintf(string(htmlContent), carta))
		} else {
			return c.SendString("Senha incorreta. Tente novamente!")
		}
	})
	//definição rota segunda_pagina
	app.Get("/page-2", func(c *fiber.Ctx) error {
		c.Type("html", "utf-8")
		// Abrindo o arquivo HTML da pasta
		file, err := os.Open("static/page-2.html")
		if err != nil {
			// Lidar com erros, como arquivo não encontrado
			return c.SendString("Erro ao carregar o arquivo HTML")
		}
		defer file.Close()

		// Lendo o conteúdo do arquivo
		htmlContent, err := io.ReadAll(file)
		if err != nil {
			// Lidar com erros de leitura
			return c.SendString("Erro ao ler o arquivo HTML")
		}

		return c.SendString(fmt.Sprintf(string(htmlContent)))
	})

	app.Listen(":3000")
}

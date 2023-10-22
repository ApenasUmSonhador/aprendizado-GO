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

// Definição de uma estrutura de fila
type Fila struct {
	elementos []string
}

// Função para enfileirar um item na fila
func (f *Fila) enfileirar(item string) {
	mu.Lock()
	f.elementos = append(f.elementos, item)
	mu.Unlock()
}

// Função para desenfileirar um item da fila
func (f *Fila) desenfileirar() (string, bool) {
	mu.Lock()
	if f.fila_vazia() {
		mu.Unlock()
		return "", false
	}
	item := f.elementos[0]
	f.elementos = f.elementos[1:]
	mu.Unlock()
	return item, true
}

// Função para verificar se a fila está vazia
func (f *Fila) fila_vazia() bool {
	return len(f.elementos) == 0
}

// Função para carregar conteúdo HTML de um arquivo
func carregarHTML(arquivo string) (string, error) {
	file, err := os.Open(arquivo)
	if err != nil {
		return "", err
	}
	defer file.Close()

	htmlContent, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(htmlContent), nil
}

var mu sync.Mutex

// Variável que será usada para armazenar a senha correta gerada
var senha_correta string

// Função para gerar uma senha e adicioná-la à fila
func GerarSenha() {
	mu.Lock()
	senha_aleatoria := rand.Intn(10) + 1
	senha_correta = strconv.Itoa(senha_aleatoria)
	mu.Unlock()
}

func main() {
	var cartas_retiradas = &Fila{
		elementos: make([]string, 0),
	}
	cartasNaFila := []string{}

	app := fiber.New()

	// Rota principal para exibir a página inicial
	app.Get("/", func(c *fiber.Ctx) error {
		c.Type("html", "utf-8")
		htmlContent, err := carregarHTML("static/inicio.html")
		if err != nil {
			return c.SendString("Erro ao carregar o arquivo HTML")
		}
		for !cartas_retiradas.fila_vazia() {
			item, _ := cartas_retiradas.desenfileirar()
			cartasNaFila = append(cartasNaFila, "<li>"+item+"</li>")
		}
		// Senha correta é adicionada ao %s e o HTML é enviado como resposta
		return c.SendString(fmt.Sprintf(string(htmlContent), senha_correta, strings.Join(cartasNaFila, "\n")))
	})

	// Rota para a ação de gerar uma nova senha
	app.Post("/gerar_senha", func(c *fiber.Ctx) error {
		go GerarSenha()
		mu.Lock()
		defer mu.Unlock()
		// Redireciona para a raiz após a geração
		return c.Redirect("/")
	})

	// Rota para confirmar a senha
	app.Post("/confirmar_senha", func(c *fiber.Ctx) error {
		// Obtém a senha digitada pelo usuário e compara com a senha correta
		senha_digitada := c.FormValue("senha_usuario")
		c.Type("html", "utf-8")
		htmlContent, err := carregarHTML("static/confirma-senha.html")
		if err != nil {
			return c.SendString("Erro ao ler o arquivo HTML")
		}
		// Conversão para inteiro, ignorando o segundo valor de retorno (possível erro de conversão)
		index, _ := strconv.Atoi(senha_correta)
		index--
		cartas := []string{"As de Copas", "2 de Copas", "8 de Espadas", "As de Ouro", "K de Paus", "6 de Espadas", "J de Copas", "3 de Ouro", "As de Paus", "J de Ouro"}
		carta := cartas[index]
		if senha_digitada == senha_correta {
			cartas_retiradas.enfileirar(carta)
			return c.SendString(fmt.Sprintf(string(htmlContent), carta))
		} else {
			return c.SendString("Senha incorreta. Tente novamente!")
		}
	})

	// Rota para a segunda página
	app.Get("/page-2", func(c *fiber.Ctx) error {
		c.Type("html", "utf-8")
		htmlContent, err := carregarHTML("static/page-2.html")
		if err != nil {
			return c.SendString("Erro ao ler o arquivo HTML")
		}
		return c.SendString(fmt.Sprintf(string(htmlContent)))
	})

	app.Listen(":3000")
}

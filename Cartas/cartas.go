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

// Função para gerar uma senha e adicioná-la à fila
func GerarSenha(senhas *[]int, senhasUsadas map[int]bool, cartasAssociadas map[int]string) int {
	mu.Lock()
	defer mu.Unlock()
	if len(*senhas) == 0 {
		return -1 // Retorna um valor especial para indicar que não há senhas disponíveis
	}
	senha_aleatoria := rand.Intn(len(*senhas))
	senha_correta := (*senhas)[senha_aleatoria]

	// Verifica se a senha gerada já foi usada
	if _, ok := senhasUsadas[senha_correta]; ok {
		// Se a senha já foi usada, gere outra
		return GerarSenha(senhas, senhasUsadas, cartasAssociadas)
	}

	// Marca a senha como usada
	senhasUsadas[senha_correta] = true

	// Verifica se a senha já está associada a uma carta
	if _, ok := cartasAssociadas[senha_correta]; ok {
		// Se a senha já está associada a uma carta, gere outra
		return GerarSenha(senhas, senhasUsadas, cartasAssociadas)
	}

	(*senhas) = append((*senhas)[:senha_aleatoria], (*senhas)[senha_aleatoria+1:]...)

	return senha_correta
}

var mu sync.Mutex

// Variável que será usada para armazenar a senha correta gerada
var senha_correta int

func main() {
	var cartas_retiradas = &Fila{
		elementos: make([]string, 0),
	}
	cartasNaFila := []string{}
	// Maneira elegante de escrever as cartas
	valores := []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
	naipes := []string{" de Copas", " de Espadas", " de Ouros", " de Paus"}
	cartas := []string{}
	for _, valor := range valores {
		for _, naipe := range naipes {
			carta := valor + naipe
			cartas = append(cartas, carta)
		}
	}
	// Declarando os tokens (maneira simplificada)
	senhas := []int{}
	for i := 1; i < 53; i++ {
		senhas = append(senhas, i)
	}

	// Map para rastrear as senhas usadas
	senhasUsadas := make(map[int]bool)
	// Map para rastrear as cartas associadas a senhas
	cartasAssociadas := make(map[int]string)

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
		SenhaToHTML := strconv.Itoa(senha_correta)
		if SenhaToHTML == "0"{
			SenhaToHTML = ""
		}
		return c.SendString(fmt.Sprintf(string(htmlContent), string(SenhaToHTML), strings.Join(cartasNaFila, "\n")))
	})

	// Rota para a ação de gerar uma nova senha
	app.Post("/gerar_senha", func(c *fiber.Ctx) error {
		senha := GerarSenha(&senhas, senhasUsadas, cartasAssociadas)
		if senha == -1 {
			return c.SendString("Todas as senhas já foram usadas")
		}
		mu.Lock()
		senha_correta = senha
		mu.Unlock()
		// Redireciona para a raiz após a geração
		return c.Redirect("/")
	})

	// Rota para confirmar a senha
	app.Post("/confirmar_senha", func(c *fiber.Ctx) error {
		// Obtém a senha digitada pelo usuário
		senha_digitada := c.FormValue("senha_usuario")
		c.Type("html", "utf-8")
		htmlContent, err := carregarHTML("static/confirma-senha.html")
		if err != nil {
			return c.SendString("Erro ao ler o arquivo HTML")
		}

		// Converte a senha correta para string
		senha_correta_str := strconv.Itoa(senha_correta)

		if senha_digitada == senha_correta_str {
			carta := cartas[senha_correta-1]

			// Verifica se a carta já foi retirada
			if _, ok := cartasAssociadas[senha_correta]; ok {
				return c.SendString("Carta já retirada. Tente novamente!")
			}

			// Marca a carta como retirada
			cartas_retiradas.enfileirar(carta)
			// Associa a senha à carta
			cartasAssociadas[senha_correta] = carta

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

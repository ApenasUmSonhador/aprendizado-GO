package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"

	"github.com/gofiber/fiber/v2"
)

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
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		c.Type("html", "utf-8") // a codificação UTF-8 é para consertar problemas de texto
		htmlContent := `
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8"> 
			<title>Gerar Senha</title>
		</head>
		<body>
			<form action="/gerar_senha" method="post">
				<input type="submit" value="Obter Senha">
			</form>			
			<p>Sua senha é: <span id="senha">%s</span></p>	
			<a href="/nova_pagina">Ir para a Segunda Página</a>		
		</body>
		</html>
		`
		//senhaCorreta é adicionada ao %s e html enviadr como resosta
		return c.SendString(fmt.Sprintf(htmlContent, senha_correta))

	})
	//rota para a ação de gerar nova senha
	app.Post("/gerar_senha", func(c *fiber.Ctx) error {
		GerarSenha()
		// redireciona a raiz após a geração
		return c.Redirect("/")
	})
	//rota para confirmar senha
	app.Post("/confirmar_senha", func(c *fiber.Ctx) error {
		//obtem senha digitada pelo usuário e faz a comparação com a correta
		senha_digitada := c.FormValue("senha_usuario")
		//retorna ou para a pagina onde mostra a carta ou para uma mensagem de erro dependendo da comparação
		c.Type("html", "utf-8")
		htmlContent := `
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>Nova Página</title>
		</head>
		<body>
			<form action="/confirmar_senha" method="post">
			<input type="text" name="senha_usuario" placeholder="Digite sua senha">
			<input type="submit" value="Confirmar">
			</form>
			<h2>Sua carta é: %s</h2>
		</body>
		</html>
		`
		// conversão para inteiro, utlizamos '_' para ignorar o segundo valor de retorno ( variável do tipo error que indica se a conversão foi bem-sucedida, podemos implementar futuramenete)
		index, _ := strconv.Atoi(senha_correta)
		// indices começam em 0, por isso a subtração
		index--
		//lista de cartas possiveis
		cartas := []string{"As de Copas", "2 de Copas", "8 de Espadas", "As de Ouro", "K de Paus", "6 de Espadas", "J de Copas", "3 de Ouro", "As de Paus", "J de Ouro"}
		// carta recebe a carta sorteada
		carta := cartas[index]
		if senha_digitada == senha_correta {
			return c.SendString(fmt.Sprintf(htmlContent, carta))
		} else {
			return c.SendString("Senha incorreta. Tente novamente!")
		}
	})
	//definição rota segunda_pagina
	app.Get("/nova_pagina", func(c *fiber.Ctx) error {
		c.Type("html", "utf-8")
		htmlContent := `
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>Nova Página</title>
		</head>
		<body>
			<form action="/confirmar_senha" method="post">
			<input type="text" name="senha_usuario" placeholder="Digite sua senha">
			<input type="submit" value="Confirmar">
			</form>
		</body>
		</html>
		`
		return c.SendString(fmt.Sprintf(htmlContent))
	})

	app.Listen(":3000")
}

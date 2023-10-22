# Progresso Desafio GO CEOS - Cartas
## Caio:
### [Link do último commit](https://github.com/ApenasUmSonhador/aprendizado-GO/commit/efa1c8f2d887edc29621cb1c4db4a12e04d4fe1d)
Caio foi o primeiro a mexer no código, Arthur e Lucas estavam muito ocupados sexta-feira. 
Portanto, por estar mais habituado a trabalhar com JavaScript, optou por fazer o esqueleto da estrutura do código que viria a ser utilizado pelos demais.
Tentou colocar em código o que entendeu do desafio proposto, mas claro que ainda faltavam ideias que seriam abordadas por Lucas e Arthur.
## Lucas:
### [Link do último commit](https://github.com/ApenasUmSonhador/aprendizado-GO/commit/8e6ddb9c1e50d74fed038cf028b495ce97809073)
Foi quem realmente começou a por as mãos no código Go e, seguindo conselhos e em conversação constante com Arthur e Caio, desenvolveu esse código sem lock e unlock, coisa que havia sido implementada por Arthur quando sua branch estava intacta.
Porém, por problemas de git, sua branch teve que ser apagada e seu último código que modificou foi o apresentado pelo commit.
## Arthur:
### [Link da sua sequência de commits](https://github.com/ApenasUmSonhador/aprendizado-GO/commits/main)
Arthur não tem um computador próprio e infelizmente, após apresentar um trabalho de FBD teve que se ausentar da UFC, dificultando seu acesso a computadores, porém, no sábado foi para a casa de sua namorada e utilizou o computador dela para trabalhar no projeto.
Arthur aproveitou o último código de Lucas e começou a fazer um papel de revisor e a ter ideia que documentou em um grupo privado com os membros Caio e Lucas.
Arthur implementou locks e unlocks, passou um pente fino no código, corrigindo problemas, por exemplo, de 2 tokens apontarem para uma mesma carta e descobriu maneiras de carregar um HTML sem ter que escrevê-lo explicitamente em "contentHTML".
Por ter mais familiaridade com Git e GitHub, já deu merge na Main e foi o último a mexer no código.
## Resumo:
**Caio:** Iniciou os trabalhos fazendo esqueleto em JavaScript. <br>
**Lucas:** Partiu do esqueleto de Caio e começou a desenvolver a lógica em Go, implementou as estrutura de fila e o esquema de adição das cartas retiradas. <br>
**Arthur:** Revisou código, fazendo algumas melhorias pontuais, não permitiu que um token aponte para mais que uma carta ,implementou Locks e Unlocks.

# Explicando código:
# Explicação do Código

O código a seguir é um servidor web escrito em Go que oferece funcionalidades relacionadas à geração de senhas exclusivas e à associação dessas senhas a cartas de um baralho. O servidor utiliza o framework Fiber para criar um aplicativo da web. Vou explicar o código, destacando a função `GerarSenha` e sua recursão.

```go
package main

import (
    // Importações de pacotes
)

// Definição de uma estrutura de fila
type Fila struct {
    elementos []string
}

// Função para enfileirar um item na fila
func (f *Fila) enfileirar(item string) {
    // Implementação
}

// Função para desenfileirar um item da fila
func (f *Fila) desenfileirar() (string, bool) {
    // Implementação
}

// Função para verificar se a fila está vazia
func (f *Fila) fila_vazia() bool {
    // Implementação
}

// Função para carregar conteúdo HTML de um arquivo
func carregarHTML(arquivo string) (string, error) {
    // Implementação
}

// Função para gerar uma senha e adicioná-la à fila
func GerarSenha(senhas *[]int, senhasUsadas map[int]bool, cartasAssociadas map[int]string) int {
    // Implementação com recursão
}

// Variável que será usada para armazenar a senha correta gerada
var senha_correta int

func main() {
    // Configuração inicial

    // Rota principal para exibir a página inicial

    // Rota para a ação de gerar uma nova senha

    // Rota para confirmar a senha

    // Rota para a segunda página

    // Início do servidor
}
```
## Explicação do Código
O código é um servidor web que oferece funcionalidades relacionadas à geração de senhas únicas e à associação dessas senhas a cartas de um baralho. Aqui está uma visão geral das principais partes do código:

**Definição de uma Estrutura de Fila:** A estrutura Fila é definida para representar uma fila que será usada para armazenar as cartas retiradas. Ela possui um slice de elementos que funciona como uma fila.

**Funções de Fila (enfileirar, desenfileirar e fila_vazia):** Três funções para operações em uma fila são definidas. enfileirar adiciona um item à fila, desenfileirar remove um item da fila e fila_vazia verifica se a fila está vazia. Essas funções são usadas para rastrear as cartas retiradas.

**Função carregarHTML:** Esta função é responsável por carregar o conteúdo HTML de um arquivo. Ela abre o arquivo, lê seu conteúdo e retorna como uma string. É utilizada para carregar modelos HTML usados nas respostas da aplicação.

**Função GerarSenha:** Esta é uma função crítica no programa. Ela gera senhas únicas e as adiciona a uma fila. A função funciona da seguinte forma:

- Bloqueia o mutex (mu) para garantir a exclusão mútua.
- Verifica se ainda existem senhas disponíveis no slice de senhas.
- Gera uma senha aleatória selecionando um índice no slice.
- Verifica se a senha gerada já foi usada. Se sim, gera outra senha chamando a própria função recursivamente. Isso garante que as senhas sejam exclusivas.
- Marca a senha como usada no mapa senhasUsadas.
- Verifica se a senha já está associada a uma carta. Se estiver, gera outra senha.
- Remove a senha gerada do slice de senhas para evitar duplicações.
- Libera o mutex e retorna a senha gerada.
- Variável senha_correta: Esta variável é usada para armazenar a senha correta gerada a cada iteração.

**Configuração Inicial:** O código define uma série de estruturas de dados e configurações iniciais, incluindo as cartas, senhas, mapas para rastrear senhas usadas e cartas associadas, e configurações do servidor Fiber.

**Rotas HTTP:**

- A rota principal (/) é usada para exibir a página inicial e inclui a senha correta, bem como as cartas na fila de retirada.
- A rota /gerar_senha é usada para gerar uma nova senha única. A função GerarSenha é chamada aqui.
- A rota /confirmar_senha permite aos usuários confirmar uma senha digitada e associá-la a uma carta retirada.
- A rota /page-2 exibe uma segunda página da aplicação.
- Início do Servidor: O servidor é iniciado na porta :3000.

## Função GerarSenha e Recursão
A função GerarSenha é crucial para o funcionamento do programa. Ela gera senhas exclusivas e garante que as senhas não sejam duplicadas ou associadas a cartas já retiradas. A recursão é usada para tratar casos em que a senha gerada não atende aos critérios necessários (já foi usada ou está associada a uma carta). Nesses casos, a função chama a si mesma recursivamente até que uma senha adequada seja gerada. Isso é fundamental para manter a integridade do sistema e garantir que as senhas sejam exclusivas.<br>

Em resumo, o código cria um servidor web que gera senhas únicas e permite a confirmação e retirada de cartas associadas a essas senhas. A função GerarSenha, com seu mecanismo de recursão, desempenha um papel central no programa, garantindo a exclusividade das senhas e evitando conflitos.

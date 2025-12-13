# Megobot 🤖

> ⚠️ **Work in Progress**: This is a simple abstraction layer over [discordgo](https://github.com/bwmarrin/discordgo) to make slash command handling easier and more straightforward. It's an ongoing project focused on simplifying common bot development patterns.

A lightweight wrapper around discordgo that reduces boilerplate and simplifies slash command creation and management.

---

## Features

- 🚀 Simple and intuitive API
- ⚡ Automatic command registration and cleanup
- ⏱️ Built-in timeout handling for long-running commands
- 🔄 Automatic deferred responses (thinking state)
- 🎯 Clean command option mapping
- 🛡️ Graceful shutdown handling

## Installation

```bash
go get github.com/LucasGabrielBravo/megobot
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/bwmarrin/discordgo"
    "github.com/LucasGabrielBravo/megobot"
)

func main() {
    // Create bot instance
    bot, err := megobot.New("YOUR_BOT_TOKEN", "YOUR_GUILD_ID")
    if err != nil {
        log.Fatal(err)
    }

    // Optional: Set custom timeout (default: 30 seconds)
    bot.SetTimeout(60 * time.Second)

    // Add a simple command
    bot.AddCommand(
        "hello",
        "Says hello to the user",
        []*discordgo.ApplicationCommandOption{
            {
                Type:        discordgo.ApplicationCommandOptionString,
                Name:        "name",
                Description: "Your name",
                Required:    true,
            },
        },
        func(s *discordgo.Session, i *discordgo.InteractionCreate, options map[string]*discordgo.ApplicationCommandInteractionDataOption) string {
            name := options["name"].StringValue()
            return fmt.Sprintf("Hello, %s! 👋", name)
        },
    )

    // Start the bot
    ctx := context.Background()
    if err := bot.Start(ctx); err != nil {
        log.Fatal(err)
    }
}
```

## API Reference

### `New(token string, guildId string) (*Megobot, error)`

Creates a new bot instance.

**Parameters:**
- `token`: Your Discord bot token
- `guildId`: The guild ID where commands will be registered (use empty string for global commands)

### `SetTimeout(timeout time.Duration)`

Sets the maximum execution time for command handlers (default: 30 seconds).

### `AddCommand(name, description string, options []*discordgo.ApplicationCommandOption, handler func)`

Registers a new slash command.

**Handler Signature:**
```go
func(s *discordgo.Session, i *discordgo.InteractionCreate, options map[string]*discordgo.ApplicationCommandInteractionDataOption) string
```

The handler receives:
- `s`: Discord session
- `i`: Interaction data
- `options`: Map of command options by name

Returns the response message as a string.

### `Start(ctx context.Context) error`

Starts the bot and blocks until the context is cancelled or a shutdown signal is received (SIGINT/SIGTERM).

## Advanced Example

```go
bot.AddCommand(
    "calculate",
    "Performs a calculation",
    []*discordgo.ApplicationCommandOption{
        {
            Type:        discordgo.ApplicationCommandOptionInteger,
            Name:        "a",
            Description: "First number",
            Required:    true,
        },
        {
            Type:        discordgo.ApplicationCommandOptionInteger,
            Name:        "b",
            Description: "Second number",
            Required:    true,
        },
        {
            Type:        discordgo.ApplicationCommandOptionString,
            Name:        "operation",
            Description: "Operation to perform",
            Required:    true,
            Choices: []*discordgo.ApplicationCommandOptionChoice{
                {Name: "Add", Value: "add"},
                {Name: "Subtract", Value: "sub"},
                {Name: "Multiply", Value: "mul"},
                {Name: "Divide", Value: "div"},
            },
        },
    },
    func(s *discordgo.Session, i *discordgo.InteractionCreate, options map[string]*discordgo.ApplicationCommandInteractionDataOption) string {
        a := options["a"].IntValue()
        b := options["b"].IntValue()
        op := options["operation"].StringValue()

        var result int64
        switch op {
        case "add":
            result = a + b
        case "sub":
            result = a - b
        case "mul":
            result = a * b
        case "div":
            if b == 0 {
                return "❌ Cannot divide by zero!"
            }
            result = a / b
        }

        return fmt.Sprintf("📊 Result: %d", result)
    },
)
```

## How It Works

1. **Automatic Deferral**: When a command is invoked, Megobot automatically sends a deferred response (shows "thinking" state)
2. **Async Execution**: Your handler runs asynchronously in a goroutine
3. **Timeout Protection**: If the handler takes longer than the configured timeout, an error message is sent
4. **Response Editing**: Once your handler completes, the deferred message is edited with the final response
5. **Clean Shutdown**: Commands are automatically unregistered when the bot shuts down

## Requirements

- Go 1.16+
- [discordgo](https://github.com/bwmarrin/discordgo) library

## License

MIT License

---

# Megobot 🤖

> ⚠️ **Trabalho em Andamento**: Este é uma simples camada de abstração sobre o [discordgo](https://github.com/bwmarrin/discordgo) para tornar o tratamento de comandos slash mais fácil e direto. É um projeto em desenvolvimento focado em simplificar padrões comuns de desenvolvimento de bots.

Um wrapper leve ao redor do discordgo que reduz código repetitivo e simplifica a criação e gerenciamento de comandos slash.

---

## Funcionalidades

- 🚀 API simples e intuitiva
- ⚡ Registro e limpeza automática de comandos
- ⏱️ Tratamento integrado de timeout para comandos demorados
- 🔄 Respostas diferidas automáticas (estado "pensando")
- 🎯 Mapeamento limpo de opções de comando
- 🛡️ Desligamento gracioso

## Instalação

```bash
go get github.com/LucasGabrielBravo/megobot
```

## Início Rápido

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/bwmarrin/discordgo"
    "github.com/LucasGabrielBravo/megobot"
)

func main() {
    // Criar instância do bot
    bot, err := megobot.New("SEU_TOKEN_DO_BOT", "ID_DO_SEU_SERVIDOR")
    if err != nil {
        log.Fatal(err)
    }

    // Opcional: Definir timeout customizado (padrão: 30 segundos)
    bot.SetTimeout(60 * time.Second)

    // Adicionar um comando simples
    bot.AddCommand(
        "ola",
        "Diz olá para o usuário",
        []*discordgo.ApplicationCommandOption{
            {
                Type:        discordgo.ApplicationCommandOptionString,
                Name:        "nome",
                Description: "Seu nome",
                Required:    true,
            },
        },
        func(s *discordgo.Session, i *discordgo.InteractionCreate, options map[string]*discordgo.ApplicationCommandInteractionDataOption) string {
            nome := options["nome"].StringValue()
            return fmt.Sprintf("Olá, %s! 👋", nome)
        },
    )

    // Iniciar o bot
    ctx := context.Background()
    if err := bot.Start(ctx); err != nil {
        log.Fatal(err)
    }
}
```

## Referência da API

### `New(token string, guildId string) (*Megobot, error)`

Cria uma nova instância do bot.

**Parâmetros:**
- `token`: Token do seu bot do Discord
- `guildId`: ID do servidor onde os comandos serão registrados (use string vazia para comandos globais)

### `SetTimeout(timeout time.Duration)`

Define o tempo máximo de execução para os handlers de comando (padrão: 30 segundos).

### `AddCommand(name, description string, options []*discordgo.ApplicationCommandOption, handler func)`

Registra um novo comando slash.

**Assinatura do Handler:**
```go
func(s *discordgo.Session, i *discordgo.InteractionCreate, options map[string]*discordgo.ApplicationCommandInteractionDataOption) string
```

O handler recebe:
- `s`: Sessão do Discord
- `i`: Dados da interação
- `options`: Map de opções do comando por nome

Retorna a mensagem de resposta como string.

### `Start(ctx context.Context) error`

Inicia o bot e bloqueia até que o contexto seja cancelado ou um sinal de desligamento seja recebido (SIGINT/SIGTERM).

## Exemplo Avançado

```go
bot.AddCommand(
    "calcular",
    "Realiza um cálculo",
    []*discordgo.ApplicationCommandOption{
        {
            Type:        discordgo.ApplicationCommandOptionInteger,
            Name:        "a",
            Description: "Primeiro número",
            Required:    true,
        },
        {
            Type:        discordgo.ApplicationCommandOptionInteger,
            Name:        "b",
            Description: "Segundo número",
            Required:    true,
        },
        {
            Type:        discordgo.ApplicationCommandOptionString,
            Name:        "operacao",
            Description: "Operação a realizar",
            Required:    true,
            Choices: []*discordgo.ApplicationCommandOptionChoice{
                {Name: "Somar", Value: "add"},
                {Name: "Subtrair", Value: "sub"},
                {Name: "Multiplicar", Value: "mul"},
                {Name: "Dividir", Value: "div"},
            },
        },
    },
    func(s *discordgo.Session, i *discordgo.InteractionCreate, options map[string]*discordgo.ApplicationCommandInteractionDataOption) string {
        a := options["a"].IntValue()
        b := options["b"].IntValue()
        op := options["operacao"].StringValue()

        var resultado int64
        switch op {
        case "add":
            resultado = a + b
        case "sub":
            resultado = a - b
        case "mul":
            resultado = a * b
        case "div":
            if b == 0 {
                return "❌ Não é possível dividir por zero!"
            }
            resultado = a / b
        }

        return fmt.Sprintf("📊 Resultado: %d", resultado)
    },
)
```

## Como Funciona

1. **Diferimento Automático**: Quando um comando é invocado, o Megobot envia automaticamente uma resposta diferida (mostra estado "pensando")
2. **Execução Assíncrona**: Seu handler roda assincronamente em uma goroutine
3. **Proteção de Timeout**: Se o handler demorar mais que o timeout configurado, uma mensagem de erro é enviada
4. **Edição de Resposta**: Quando seu handler completa, a mensagem diferida é editada com a resposta final
5. **Desligamento Limpo**: Comandos são automaticamente desregistrados quando o bot desliga

## Requisitos

- Go 1.16+
- Biblioteca [discordgo](https://github.com/bwmarrin/discordgo)

## Licença

MIT License
package main

import (
	"flag"
	"fmt"
	"golang-project/cambio"
	"golang-project/server"
	"os"
	"time"
)

func main() {
	if tz := os.Getenv("TZ"); tz != "" {
		loc, err := time.LoadLocation(tz)
		if err == nil {
			time.Local = loc
		}
	}

	// Flag para escolher o modo de execução
	serverMode := flag.Bool("server", false, "Executar em modo servidor")
	port := flag.String("port", "8080", "Porta do servidor")
	flag.Parse()

	if *serverMode {
		// Modo servidor - API REST + Interface React
		server.StartServerChi(*port)
	} else {
		// Modo CLI original
		runCLIMode()
	}
}

func runCLIMode() {
	servico := cambio.NewServicoTaxasSimples()

	fmt.Println("=== SISTEMA DE CÂMBIO SIMPLIFICADO ===")

	err := servico.InicializarTaxas()
	if err != nil {
		fmt.Printf("Erro ao carregar taxas: %v\n", err)
		return
	}

	servico.ExibirTaxas()

	fmt.Println("\n=== EXEMPLOS DE CONVERSÃO ===")

	conversoes := []struct {
		valor   float64
		origem  string
		destino string
	}{
		{1000.0, "USD", "BRL"},
		{500.0, "EUR", "BRL"},
		{100.0, "BRL", "USD"},
		{1000.0, "GBP", "JPY"},
	}

	for _, conv := range conversoes {
		valorConvertido, err := servico.CalcularConversao(conv.valor, conv.origem, conv.destino)
		if err != nil {
			fmt.Printf("Erro na conversão %s->%s: %v\n", conv.origem, conv.destino, err)
			continue
		}

		fmt.Printf("💱 %s %.2f = %s %.2f\n",
			conv.origem, conv.valor,
			conv.destino, valorConvertido)
	}
}

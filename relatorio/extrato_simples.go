package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"golang-project/cambio"
	"log"
	"os"
	"strconv"
	"time"
)

var servicoCambio = cambio.NewServicoTaxasSimples()
var taxasCambioFallback = map[string]map[string]float64{
	"USD": {"BRL": 5.42, "EUR": 0.92, "GBP": 0.79, "JPY": 149.50},
	"EUR": {"BRL": 5.89, "USD": 1.09, "GBP": 0.86, "JPY": 162.80},
	"BRL": {"USD": 0.18, "EUR": 0.17, "GBP": 0.15, "JPY": 27.50},
	"GBP": {"BRL": 6.85, "USD": 1.27, "EUR": 1.16, "JPY": 189.20},
	"JPY": {"BRL": 0.036, "USD": 0.0067, "EUR": 0.0061, "GBP": 0.0053},
}

type Transacao struct {
	ID           string  `json:"id"`
	Data         string  `json:"data"`
	Hora         string  `json:"hora"`
	NomeCliente  string  `json:"nome_cliente"`
	CpfCnpj      string  `json:"cpf_cnpj"`
	MoedaOrigem  string  `json:"moeda_origem,omitempty"`
	MoedaDestino string  `json:"moeda_destino,omitempty"`
	ValorOrigem  float64 `json:"valor_origem"`
	ValorDestino float64 `json:"valor_destino"`
	Status       string  `json:"status"`
	Comissao     float64 `json:"comissao"`
	TipoOperacao string  `json:"tipo_operacao"`
	Canal        string  `json:"canal"`
	Observacoes  string  `json:"observacoes"`
}

type DadosCambio struct {
	TransacoesCambio []Transacao `json:"transacoes_cambio"`
}

func main() {

	fmt.Println("Inicializando sistema...")
	err := servicoCambio.InicializarTaxas()
	if err != nil {
		fmt.Printf("⚠️ Erro ao carregar taxas da API: %v\n", err)
		fmt.Println("O sistema funcionará com taxas de fallback.")
	}

	for {
		fmt.Println("\n=== SISTEMA DE TRANSAÇÕES DE CÂMBIO ===")
		fmt.Println("1. Ver transações")
		fmt.Println("2. Inserir nova transação")
		fmt.Println("3. Recarregar taxas de câmbio")
		fmt.Println("4. Ver taxas atuais")
		fmt.Println("5. Sair")
		fmt.Print("Escolha uma opção: ")

		var opcao string
		fmt.Scanln(&opcao)

		switch opcao {
		case "1":
			verTransacoes()
		case "2":
			inserirTransacao()
		case "3":
			recarregarTaxasCambio()
		case "4":
			servicoCambio.ExibirTaxas()
		case "5":
			fmt.Println("Saindo...")
			return
		default:
			fmt.Println("Opção inválida! Digite 1, 2, 3, 4 ou 5")
		}
	}
}

func verTransacoes() {
	arquivo, err := os.ReadFile("../transacoes_cambio.json")
	if err != nil {
		log.Fatal("Erro ao ler arquivo:", err)
	}

	var dados DadosCambio
	err = json.Unmarshal(arquivo, &dados)
	if err != nil {
		log.Fatal("Erro ao processar JSON:", err)
	}

	fmt.Println("\n=== RESUMO DE TRANSAÇÕES DE CÂMBIO ===")

	total := 0
	concluidas := 0
	valorTotal := 0.0
	comissaoTotal := 0.0

	for _, transacao := range dados.TransacoesCambio {
		total++

		fmt.Printf("Cliente: %s\n", transacao.NomeCliente)
		fmt.Printf("Data: %s | Status: %s\n", transacao.Data, transacao.Status)

		if transacao.MoedaOrigem != "" && transacao.MoedaDestino != "" {
			fmt.Printf("Conversão: %s %.2f → %s %.2f | Comissão: R$ %.2f\n",
				transacao.MoedaOrigem, transacao.ValorOrigem,
				transacao.MoedaDestino, transacao.ValorDestino, transacao.Comissao)
		} else {
			fmt.Printf("Valor: %.2f → %.2f | Comissão: R$ %.2f\n",
				transacao.ValorOrigem, transacao.ValorDestino, transacao.Comissao)
		}
		fmt.Println("----------------------------------------")

		if transacao.Status == "concluida" {
			concluidas++
			valorTotal += transacao.ValorDestino
			comissaoTotal += transacao.Comissao
		}
	}

	fmt.Printf("\nRESUMO FINAL:\n")
	fmt.Printf("Total de transações: %d\n", total)
	fmt.Printf("Transações concluídas: %d\n", concluidas)
	fmt.Printf("Valor total movimentado: R$ %.2f\n", valorTotal)
	fmt.Printf("Total de comissões: R$ %.2f\n", comissaoTotal)
}

func calcularConversao(valorOrigem float64, moedaOrigem, moedaDestino string) float64 {
	if moedaOrigem == moedaDestino {
		return valorOrigem
	}

	if servicoCambio.EstaoCarregadas() {
		valorConvertido, err := servicoCambio.CalcularConversao(valorOrigem, moedaOrigem, moedaDestino)
		if err == nil {
			fmt.Printf("Usando taxa carregada na inicialização\n")
			return valorConvertido
		}
		fmt.Printf("Erro ao usar taxa carregada: %v\n", err)
	}

	fmt.Printf("Usando taxas de fallback\n")

	if taxas, existe := taxasCambioFallback[moedaOrigem]; existe {
		if taxa, existeTaxa := taxas[moedaDestino]; existeTaxa {
			return valorOrigem * taxa
		}
	}

	fmt.Printf("Taxa de câmbio não encontrada para %s -> %s. Usando valor 1:1\n", moedaOrigem, moedaDestino)
	return valorOrigem
}

func inserirTransacao() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("\n=== INSERIR NOVA TRANSAÇÃO ===")

	fmt.Print("Nome do cliente: ")
	scanner.Scan()
	nomeCliente := scanner.Text()

	fmt.Print("CPF/CNPJ: ")
	scanner.Scan()
	cpfCnpj := scanner.Text()

	fmt.Print("Moeda de origem (USD, EUR, BRL, etc.): ")
	scanner.Scan()
	moedaOrigem := scanner.Text()

	fmt.Print("Valor de origem: ")
	scanner.Scan()
	valorOrigemStr := scanner.Text()
	valorOrigem, _ := strconv.ParseFloat(valorOrigemStr, 64)

	fmt.Print("Moeda de destino (USD, EUR, BRL, etc.): ")
	scanner.Scan()
	moedaDestino := scanner.Text()

	valorDestino := calcularConversao(valorOrigem, moedaOrigem, moedaDestino)
	fmt.Printf("Valor convertido: %s %.2f → %s %.2f\n",
		moedaOrigem, valorOrigem, moedaDestino, valorDestino)

	fmt.Print("Comissão: ")
	scanner.Scan()
	comissaoStr := scanner.Text()
	comissao, _ := strconv.ParseFloat(comissaoStr, 64)

	fmt.Print("Tipo de operação (compra/venda): ")
	scanner.Scan()
	tipoOperacao := scanner.Text()

	fmt.Print("Canal (presencial/internet_banking/mobile_app/telefone): ")
	scanner.Scan()
	canal := scanner.Text()

	fmt.Print("Status (concluida/pendente/cancelada): ")
	scanner.Scan()
	status := scanner.Text()

	fmt.Print("Observações: ")
	scanner.Scan()
	observacoes := scanner.Text()

	agora := time.Now()
	novaTransacao := Transacao{
		ID:           fmt.Sprintf("TXN-%s-%03d", agora.Format("2006"), time.Now().Unix()%1000),
		Data:         agora.Format("2006-01-02"),
		Hora:         agora.Format("15:04:05"),
		NomeCliente:  nomeCliente,
		CpfCnpj:      cpfCnpj,
		MoedaOrigem:  moedaOrigem,
		MoedaDestino: moedaDestino,
		ValorOrigem:  valorOrigem,
		ValorDestino: valorDestino,
		Status:       status,
		Comissao:     comissao,
		TipoOperacao: tipoOperacao,
		Canal:        canal,
		Observacoes:  observacoes,
	}

	// Ler arquivo atual
	arquivo, err := os.ReadFile("../transacoes_cambio.json")
	if err != nil {
		log.Fatal("Erro ao ler arquivo:", err)
	}

	var dados DadosCambio
	err = json.Unmarshal(arquivo, &dados)
	if err != nil {
		log.Fatal("Erro ao processar JSON:", err)
	}

	dados.TransacoesCambio = append(dados.TransacoesCambio, novaTransacao)

	novoJSON, err := json.MarshalIndent(dados, "", "  ")
	if err != nil {
		log.Fatal("Erro ao converter para JSON:", err)
	}

	err = os.WriteFile("../transacoes_cambio.json", novoJSON, 0644)
	if err != nil {
		log.Fatal("Erro ao salvar arquivo:", err)
	}

	fmt.Println("\nTransação inserida com sucesso!")
	fmt.Printf("ID da transação: %s\n", novaTransacao.ID)
}

// recarregarTaxasCambio força o recarregamento das taxas
func recarregarTaxasCambio() {
	fmt.Println("\n=== RECARREGANDO TAXAS DE CÂMBIO ===")

	err := servicoCambio.RecarregarTaxas()
	if err != nil {
		fmt.Printf("Erro ao recarregar taxas: %v\n", err)
		fmt.Println("O sistema continuará usando as taxas de fallback.")
		return
	}

	fmt.Println("Taxas recarregadas com sucesso!")
	servicoCambio.ExibirTaxas()
}

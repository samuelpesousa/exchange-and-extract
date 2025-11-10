package cambio

import (
	"fmt"
)

type ServicoTaxasCambio struct {
	cliente *CambioClient
	cache   *GerenciadorCache
}

func NewServicoTaxasCambio() *ServicoTaxasCambio {
	return &ServicoTaxasCambio{
		cliente: NewCambioClient(),
		cache:   NewGerenciadorCache(),
	}
}

func (s *ServicoTaxasCambio) ObterTaxasAtualizadas() (map[string]map[string]float64, error) {
	if taxas, valido := s.cache.CarregarCache(); valido {
		return taxas, nil
	}

	fmt.Println("Buscando taxas atualizadas da API...")
	taxas, err := s.cliente.BuscarTaxasParaTodasMoedas()
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar taxas da API: %w", err)
	}

	err = s.cache.SalvarCache(taxas)
	if err != nil {
		fmt.Printf("Aviso: erro ao salvar cache: %v\n", err)

	}

	return taxas, nil
}

func (s *ServicoTaxasCambio) ObterTaxasRapidas() (map[string]map[string]float64, error) {
	taxas, valido := s.cache.CarregarCache()
	if !valido {
		return nil, fmt.Errorf("cache não disponível ou expirado")
	}
	return taxas, nil
}

func (s *ServicoTaxasCambio) ForcarAtualizacao() (map[string]map[string]float64, error) {
	fmt.Println("Forçando atualização das taxas...")

	taxas, err := s.cliente.BuscarTaxasParaTodasMoedas()
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar taxas da API: %w", err)
	}

	// Salvar no cache
	err = s.cache.SalvarCache(taxas)
	if err != nil {
		fmt.Printf("Aviso: erro ao salvar cache: %v\n", err)
	}

	return taxas, nil
}

func (s *ServicoTaxasCambio) CalcularConversaoComAPI(valor float64, moedaOrigem, moedaDestino string) (float64, error) {
	taxas, err := s.ObterTaxasAtualizadas()
	if err != nil {
		return 0, fmt.Errorf("erro ao obter taxas: %w", err)
	}

	return s.cliente.CalcularConversao(valor, moedaOrigem, moedaDestino, taxas)
}

func (s *ServicoTaxasCambio) LimparCache() error {
	return s.cache.LimparCache()
}

func (s *ServicoTaxasCambio) ExibirStatusTaxas() {
	taxas, valido := s.cache.CarregarCache()
	if !valido {
		fmt.Println("Nenhuma taxa em cache. Execute uma atualização primeiro.")
		return
	}

	fmt.Println("\n=== STATUS DAS TAXAS DE CÂMBIO ===")
	for moedaOrigem, conversoes := range taxas {
		fmt.Printf("\n%s:", moedaOrigem)
		for moedaDestino, taxa := range conversoes {
			fmt.Printf("  %s: %.4f", moedaDestino, taxa)
		}
		fmt.Println()
	}
}

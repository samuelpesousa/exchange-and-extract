package cambio

import (
	"fmt"
	"sync"
)

type ServicoTaxasSimples struct {
	cliente    *CambioClient
	taxas      map[string]map[string]float64
	carregadas bool
	mutex      sync.RWMutex
}

func NewServicoTaxasSimples() *ServicoTaxasSimples {
	return &ServicoTaxasSimples{
		cliente:    NewCambioClient(),
		carregadas: false,
	}
}

func (s *ServicoTaxasSimples) InicializarTaxas() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.carregadas {
		fmt.Println("Taxas já carregadas na memória")
		return nil
	}

	fmt.Println("Carregando taxas de câmbio da API (apenas na inicialização)...")
	taxas, err := s.cliente.BuscarTaxasParaTodasMoedas()
	if err != nil {
		return fmt.Errorf("erro ao carregar taxas na inicialização: %w", err)
	}

	s.taxas = taxas
	s.carregadas = true

	fmt.Println("Taxas carregadas com sucesso!")
	return nil
}

func (s *ServicoTaxasSimples) ObterTaxas() (map[string]map[string]float64, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if !s.carregadas {
		return nil, fmt.Errorf("taxas não foram carregadas. Execute InicializarTaxas() primeiro")
	}

	return s.taxas, nil
}

func (s *ServicoTaxasSimples) CalcularConversao(valor float64, moedaOrigem, moedaDestino string) (float64, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if !s.carregadas {
		return 0, fmt.Errorf("taxas não foram carregadas. Execute InicializarTaxas() primeiro")
	}

	return s.cliente.CalcularConversao(valor, moedaOrigem, moedaDestino, s.taxas)
}

func (s *ServicoTaxasSimples) EstaoCarregadas() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.carregadas
}

func (s *ServicoTaxasSimples) ExibirTaxas() {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if !s.carregadas {
		fmt.Println("❌ Taxas não carregadas. Execute a inicialização primeiro.")
		return
	}

	fmt.Println("\n=== TAXAS DE CÂMBIO CARREGADAS ===")
	for moedaOrigem, conversoes := range s.taxas {
		fmt.Printf("\n%s:", moedaOrigem)
		for moedaDestino, taxa := range conversoes {
			fmt.Printf("  %s: %.4f", moedaDestino, taxa)
		}
		fmt.Println()
	}
}

func (s *ServicoTaxasSimples) RecarregarTaxas() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	fmt.Println("Recarregando taxas da API...")
	taxas, err := s.cliente.BuscarTaxasParaTodasMoedas()
	if err != nil {
		return fmt.Errorf("erro ao recarregar taxas: %w", err)
	}

	s.taxas = taxas
	s.carregadas = true

	fmt.Println(" Taxas recarregadas com sucesso!")
	return nil
}

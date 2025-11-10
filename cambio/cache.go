package cambio

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type CacheData struct {
	Timestamp       int64                         `json:"timestamp"`
	TaxasCambio     map[string]map[string]float64 `json:"taxas_cambio"`
	ValidadePeriodo int64                         `json:"validade_periodo"`
}

const (
	CACHE_FILE     = "taxas_cambio_cache.json"
	CACHE_VALIDADE = 3600
)

type GerenciadorCache struct {
	arquivo  string
	validade int64
}

func NewGerenciadorCache() *GerenciadorCache {
	return &GerenciadorCache{
		arquivo:  CACHE_FILE,
		validade: CACHE_VALIDADE,
	}
}

func (g *GerenciadorCache) CarregarCache() (map[string]map[string]float64, bool) {
	// Verificar se o arquivo existe
	if _, err := os.Stat(g.arquivo); os.IsNotExist(err) {
		return nil, false
	}

	data, err := os.ReadFile(g.arquivo)
	if err != nil {
		fmt.Printf("Erro ao ler cache: %v\n", err)
		return nil, false
	}

	var cache CacheData
	err = json.Unmarshal(data, &cache)
	if err != nil {
		fmt.Printf("Erro ao deserializar cache: %v\n", err)
		return nil, false
	}

	agora := time.Now().Unix()
	if agora-cache.Timestamp > g.validade {
		fmt.Println("Cache expirado, buscando dados atualizados...")
		return nil, false
	}

	fmt.Printf("Cache válido encontrado (atualizado há %d segundos)\n", agora-cache.Timestamp)
	return cache.TaxasCambio, true
}

func (g *GerenciadorCache) SalvarCache(taxas map[string]map[string]float64) error {
	cache := CacheData{
		Timestamp:       time.Now().Unix(),
		TaxasCambio:     taxas,
		ValidadePeriodo: g.validade,
	}

	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return fmt.Errorf("erro ao serializar cache: %w", err)
	}

	err = os.WriteFile(g.arquivo, data, 0644)
	if err != nil {
		return fmt.Errorf("erro ao salvar cache: %w", err)
	}

	fmt.Println("Cache atualizado com sucesso")
	return nil
}

func (g *GerenciadorCache) LimparCache() error {
	if _, err := os.Stat(g.arquivo); os.IsNotExist(err) {
		return nil
	}

	err := os.Remove(g.arquivo)
	if err != nil {
		return fmt.Errorf("erro ao limpar cache: %w", err)
	}

	fmt.Println("Cache limpo com sucesso")
	return nil
}

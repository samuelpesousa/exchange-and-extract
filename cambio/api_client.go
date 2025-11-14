package cambio

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type SimpleExchangeResponse struct {
	Result   string             `json:"result,omitempty"`
	BaseCode string             `json:"base_code,omitempty"`
	Rates    map[string]float64 `json:"rates"`
}

type ExchangeRateResponse struct {
	Result             string             `json:"result"`
	Documentation      string             `json:"documentation"`
	TermsOfUse         string             `json:"terms_of_use"`
	TimeLastUpdateUnix int64              `json:"time_last_update_unix"`
	TimeLastUpdateUtc  string             `json:"time_last_update_utc"`
	TimeNextUpdateUnix int64              `json:"time_next_update_unix"`
	TimeNextUpdateUtc  string             `json:"time_next_update_utc"`
	BaseCode           string             `json:"base_code"`
	ConversionRates    map[string]float64 `json:"conversion_rates"`
}

type CambioClient struct {
	baseURL string
	timeout time.Duration
}

func NewCambioClient() *CambioClient {
	return &CambioClient{
		baseURL: "https://api.fxratesapi.com/latest",
		timeout: 15 * time.Second,
	}
}

func (c *CambioClient) BuscarTaxasCambio(moedaBase string) (map[string]float64, error) {
	client := &http.Client{
		Timeout: c.timeout,
	}

	url := fmt.Sprintf("%s?base=%s", c.baseURL, moedaBase)

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer requisição para API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API retornou status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta da API: %w", err)
	}

	var simpleData SimpleExchangeResponse
	err = json.Unmarshal(body, &simpleData)
	if err == nil && len(simpleData.Rates) > 0 {
		return simpleData.Rates, nil
	}

	var exchangeData ExchangeRateResponse
	err = json.Unmarshal(body, &exchangeData)
	if err != nil {
		return nil, fmt.Errorf("erro ao deserializar resposta da API: %w", err)
	}

	if exchangeData.Result != "" && exchangeData.Result != "success" {
		return nil, fmt.Errorf("API retornou resultado: %s", exchangeData.Result)
	}

	if len(exchangeData.ConversionRates) > 0 {
		return exchangeData.ConversionRates, nil
	}

	return nil, fmt.Errorf("resposta da API não contém taxas válidas")
}

func (c *CambioClient) BuscarTaxasParaTodasMoedas() (map[string]map[string]float64, error) {
	moedas := []string{"USD", "EUR", "BRL", "GBP", "JPY"}
	taxasCompletas := make(map[string]map[string]float64)

	var mu sync.Mutex
	var wg sync.WaitGroup
	errors := make(chan error, len(moedas))

	for _, moedaBase := range moedas {
		wg.Add(1)
		go func(moeda string) {
			defer wg.Done()

			fmt.Printf("Buscando taxas para %s...\n", moeda)

			taxas, err := c.BuscarTaxasCambio(moeda)
			if err != nil {
				errors <- fmt.Errorf("erro ao buscar taxas para %s: %w", moeda, err)
				return
			}

			taxasFiltradas := make(map[string]float64)
			for _, moedaDestino := range moedas {
				if moedaDestino != moeda {
					if taxa, existe := taxas[moedaDestino]; existe {
						taxasFiltradas[moedaDestino] = taxa
					}
				}
			}

			mu.Lock()
			taxasCompletas[moeda] = taxasFiltradas
			mu.Unlock()

			// Sleep para evitar rate limiting (menor que antes pois é paralelo)
			time.Sleep(100 * time.Millisecond)
		}(moedaBase)
	}

	// Aguardar todas as goroutines terminarem
	wg.Wait()
	close(errors)

	// Verificar se houve algum erro
	select {
	case err := <-errors:
		if err != nil {
			return nil, err
		}
	default:
		// Nenhum erro
	}

	return taxasCompletas, nil
}

func (c *CambioClient) CalcularConversao(valor float64, moedaOrigem, moedaDestino string, taxas map[string]map[string]float64) (float64, error) {
	if moedaOrigem == moedaDestino {
		return valor, nil
	}

	if taxasMoeda, existe := taxas[moedaOrigem]; existe {
		if taxa, existeTaxa := taxasMoeda[moedaDestino]; existeTaxa {
			return valor * taxa, nil
		}
	}

	return 0, fmt.Errorf("taxa de câmbio não encontrada para %s -> %s", moedaOrigem, moedaDestino)
}

import React, { useState, useEffect } from "react";
import axios from "axios";
import {
  RefreshCw,
  TrendingUp,
  Wifi,
  WifiOff,
  Loader2,
  ArrowRightLeft,
  AlertCircle,
  DollarSign,
  CheckCircle,
  ShoppingCart,
  TrendingDown,
} from "lucide-react";

interface OperationData {
  tipo: string;
  moedaOrigem: string;
  moedaDestino: string;
  valorOrigem: string;
}

interface OperationResult {
  id: number;
  data_transacao: string;
  tipo: string;
  moeda_origem: string;
  moeda_destino: string;
  valor_origem: number;
  valor_destino: number;
  taxa_cambio: number;
  status: string;
}

interface LoadingState {
  operation: boolean;
  rates: boolean;
  refresh: boolean;
}

interface ErrorState {
  operation?: string;
  rates?: string;
}

interface ApiStatus {
  online: boolean;
}

const ExchangeRate: React.FC = () => {
  const [operation, setOperation] = useState<OperationData>({
    tipo: "Compra",
    moedaOrigem: "",
    moedaDestino: "",
    valorOrigem: "",
  });

  const [result, setResult] = useState<OperationResult | null>(null);
  const [loading, setLoading] = useState<LoadingState>({
    operation: false,
    rates: false,
    refresh: false,
  });
  const [error, setError] = useState<ErrorState>({});
  const [exchangeRates, setExchangeRates] = useState<Record<string, number>>(
    {}
  );
  const [lastUpdate, setLastUpdate] = useState<string>("");
  const [apiStatus, setApiStatus] = useState<ApiStatus>({ online: false });
  const [availableCurrencies, setAvailableCurrencies] = useState<string[]>([]);

  const API_BASE = "/api";

  const checkApiStatus = async () => {
    try {
      await axios.get(`${API_BASE}/health`);
      setApiStatus({ online: true });
    } catch {
      setApiStatus({ online: false });
    }
  };

  const loadRates = async () => {
    setLoading((prev) => ({ ...prev, rates: true }));
    setError((prev) => ({ ...prev, rates: undefined }));

    try {
      const response = await axios.get(`${API_BASE}/taxas`);
      const rates = response.data.taxas;

      const flatRates: Record<string, number> = {};
      const currenciesSet = new Set<string>();

      Object.keys(rates).forEach((fromCurrency) => {
        currenciesSet.add(fromCurrency);
        Object.keys(rates[fromCurrency]).forEach((toCurrency) => {
          currenciesSet.add(toCurrency);
          flatRates[`${fromCurrency}_${toCurrency}`] =
            rates[fromCurrency][toCurrency];
        });
      });

      setExchangeRates(flatRates);
      setAvailableCurrencies(Array.from(currenciesSet).sort());
      setLastUpdate(new Date().toLocaleString("pt-BR"));
      setApiStatus({ online: true });
    } catch (err) {
      setError((prev) => ({
        ...prev,
        rates: "Erro ao carregar taxas de câmbio",
      }));
      setApiStatus({ online: false });
    } finally {
      setLoading((prev) => ({ ...prev, rates: false }));
    }
  };

  const refreshRates = async () => {
    setLoading((prev) => ({ ...prev, refresh: true }));
    await loadRates();
    setLoading((prev) => ({ ...prev, refresh: false }));
  };

  const processOperation = async (e: React.FormEvent) => {
    e.preventDefault();

    if (
      !operation.moedaOrigem ||
      !operation.moedaDestino ||
      !operation.valorOrigem
    ) {
      setError((prev) => ({
        ...prev,
        operation: "Preencha todos os campos obrigatórios",
      }));
      return;
    }

    const valor = parseFloat(operation.valorOrigem);
    if (isNaN(valor) || valor <= 0) {
      setError((prev) => ({
        ...prev,
        operation: "Valor deve ser um número positivo",
      }));
      return;
    }

    setLoading((prev) => ({ ...prev, operation: true }));
    setError((prev) => ({ ...prev, operation: undefined }));
    setResult(null);

    try {
      const response = await axios.post(`${API_BASE}/transacoes`, {
        tipo: operation.tipo,
        moeda_origem: operation.moedaOrigem,
        moeda_destino: operation.moedaDestino,
        valor_origem: valor,
      });

      setResult(response.data);
      setApiStatus({ online: true });

      setOperation({
        tipo: operation.tipo,
        moedaOrigem: "",
        moedaDestino: "",
        valorOrigem: "",
      });
    } catch (err: any) {
      const errorMessage =
        err.response?.data?.error || "Erro ao processar operação";
      setError((prev) => ({ ...prev, operation: errorMessage }));
      setApiStatus({ online: false });
    } finally {
      setLoading((prev) => ({ ...prev, operation: false }));
    }
  };

  useEffect(() => {
    loadRates();
    checkApiStatus();

    const ratesInterval = setInterval(loadRates, 60000);
    const statusInterval = setInterval(checkApiStatus, 30000);

    return () => {
      clearInterval(ratesInterval);
      clearInterval(statusInterval);
    };
  }, []);

  const calculateEstimatedValue = (): number | null => {
    if (
      !operation.moedaOrigem ||
      !operation.moedaDestino ||
      !operation.valorOrigem
    ) {
      return null;
    }

    const valor = parseFloat(operation.valorOrigem);
    if (isNaN(valor)) return null;

    const rateKey = `${operation.moedaOrigem}_${operation.moedaDestino}`;
    const rate = exchangeRates[rateKey];

    return rate ? valor * rate : null;
  };

  const estimatedValue = calculateEstimatedValue();

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-white py-8 px-4">
      <div className="max-w-4xl mx-auto">
        <div className="bg-white rounded-lg shadow-md p-6 mb-6">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-3">
              <TrendingUp className="w-8 h-8 text-blue-600" />
              <div>
                <h1 className="text-2xl font-bold text-gray-800">
                  Operações de Câmbio
                </h1>
                <p className="text-sm text-gray-600">
                  Sistema de Compra e Venda de Moedas
                </p>
              </div>
            </div>
            <div className="flex items-center space-x-4">
              <div className="flex items-center space-x-2">
                {apiStatus.online ? (
                  <>
                    <Wifi className="w-5 h-5 text-green-500" />
                    <span className="text-sm text-green-600 font-medium">
                      Online
                    </span>
                  </>
                ) : (
                  <>
                    <WifiOff className="w-5 h-5 text-red-500" />
                    <span className="text-sm text-red-600 font-medium">
                      Offline
                    </span>
                  </>
                )}
              </div>
              <button
                onClick={refreshRates}
                disabled={loading.refresh}
                className="flex items-center space-x-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors"
              >
                <RefreshCw
                  className={`w-4 h-4 ${loading.refresh ? "animate-spin" : ""}`}
                />
                <span>Atualizar</span>
              </button>
            </div>
          </div>
          {lastUpdate && (
            <p className="text-xs text-gray-500 mt-2">
              Última atualização: {lastUpdate}
            </p>
          )}
        </div>

        {error.rates && (
          <div className="bg-red-50 border border-red-200 rounded-lg p-4 mb-6">
            <div className="flex items-center space-x-2">
              <AlertCircle className="w-5 h-5 text-red-600" />
              <p className="text-red-700">{error.rates}</p>
            </div>
          </div>
        )}

        <div className="grid md:grid-cols-2 gap-6">
          <div className="bg-white rounded-lg shadow-md p-6">
            <h2 className="text-xl font-bold text-gray-800 mb-4 flex items-center">
              <ArrowRightLeft className="w-5 h-5 mr-2 text-blue-600" />
              Nova Operação
            </h2>

            <form onSubmit={processOperation} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Tipo de Operação
                </label>
                <div className="grid grid-cols-2 gap-3">
                  <button
                    type="button"
                    onClick={() =>
                      setOperation((prev) => ({ ...prev, tipo: "Compra" }))
                    }
                    className={`flex items-center justify-center space-x-2 px-4 py-3 rounded-lg border-2 transition-all ${
                      operation.tipo === "Compra"
                        ? "border-blue-600 bg-blue-50 text-blue-700"
                        : "border-gray-300 bg-white text-gray-700 hover:border-gray-400"
                    }`}
                  >
                    <ShoppingCart className="w-4 h-4" />
                    <span className="font-medium">Compra</span>
                  </button>
                  <button
                    type="button"
                    onClick={() =>
                      setOperation((prev) => ({ ...prev, tipo: "Venda" }))
                    }
                    className={`flex items-center justify-center space-x-2 px-4 py-3 rounded-lg border-2 transition-all ${
                      operation.tipo === "Venda"
                        ? "border-blue-600 bg-blue-50 text-blue-700"
                        : "border-gray-300 bg-white text-gray-700 hover:border-gray-400"
                    }`}
                  >
                    <TrendingDown className="w-4 h-4" />
                    <span className="font-medium">Venda</span>
                  </button>
                </div>
              </div>

              <div>
                <label
                  htmlFor="moedaOrigem"
                  className="block text-sm font-medium text-gray-700 mb-2"
                >
                  Moeda de Origem *
                </label>
                <select
                  id="moedaOrigem"
                  value={operation.moedaOrigem}
                  onChange={(e) =>
                    setOperation((prev) => ({
                      ...prev,
                      moedaOrigem: e.target.value,
                    }))
                  }
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  required
                >
                  <option value="">Selecione a moeda</option>
                  {availableCurrencies.map((currency) => (
                    <option key={currency} value={currency}>
                      {currency}
                    </option>
                  ))}
                </select>
              </div>

              <div>
                <label
                  htmlFor="moedaDestino"
                  className="block text-sm font-medium text-gray-700 mb-2"
                >
                  Moeda de Destino *
                </label>
                <select
                  id="moedaDestino"
                  value={operation.moedaDestino}
                  onChange={(e) =>
                    setOperation((prev) => ({
                      ...prev,
                      moedaDestino: e.target.value,
                    }))
                  }
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  required
                >
                  <option value="">Selecione a moeda</option>
                  {availableCurrencies.map((currency) => (
                    <option key={currency} value={currency}>
                      {currency}
                    </option>
                  ))}
                </select>
              </div>

              <div>
                <label
                  htmlFor="valorOrigem"
                  className="block text-sm font-medium text-gray-700 mb-2"
                >
                  Valor ({operation.moedaOrigem || "Origem"}) *
                </label>
                <div className="relative">
                  <DollarSign className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" />
                  <input
                    type="number"
                    id="valorOrigem"
                    value={operation.valorOrigem}
                    onChange={(e) =>
                      setOperation((prev) => ({
                        ...prev,
                        valorOrigem: e.target.value,
                      }))
                    }
                    placeholder="0.00"
                    step="0.01"
                    min="0"
                    className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    required
                  />
                </div>
              </div>

              {estimatedValue !== null && (
                <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
                  <div className="flex items-center justify-between">
                    <span className="text-sm font-medium text-blue-700">
                      Valor Estimado:
                    </span>
                    <span className="text-lg font-bold text-blue-900">
                      {estimatedValue.toFixed(2)} {operation.moedaDestino}
                    </span>
                  </div>
                  <p className="text-xs text-blue-600 mt-1">
                    Taxa:{" "}
                    {exchangeRates[
                      `${operation.moedaOrigem}_${operation.moedaDestino}`
                    ]?.toFixed(4)}
                  </p>
                </div>
              )}

              {error.operation && (
                <div className="bg-red-50 border border-red-200 rounded-lg p-3">
                  <div className="flex items-center space-x-2">
                    <AlertCircle className="w-4 h-4 text-red-600 flex-shrink-0" />
                    <p className="text-sm text-red-700">{error.operation}</p>
                  </div>
                </div>
              )}

              <button
                type="submit"
                disabled={loading.operation || !apiStatus.online}
                className="w-full flex items-center justify-center space-x-2 px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors font-medium"
              >
                {loading.operation ? (
                  <>
                    <Loader2 className="w-5 h-5 animate-spin" />
                    <span>Processando...</span>
                  </>
                ) : (
                  <>
                    <CheckCircle className="w-5 h-5" />
                    <span>Confirmar Operação</span>
                  </>
                )}
              </button>
            </form>
          </div>

          <div className="bg-white rounded-lg shadow-md p-6">
            <h2 className="text-xl font-bold text-gray-800 mb-4 flex items-center">
              <CheckCircle className="w-5 h-5 mr-2 text-green-600" />
              Resultado da Operação
            </h2>

            {result ? (
              <div className="space-y-4">
                <div className="bg-green-50 border border-green-200 rounded-lg p-4">
                  <div className="flex items-center space-x-2">
                    <CheckCircle className="w-5 h-5 text-green-600" />
                    <span className="text-green-700 font-medium">
                      Operação realizada com sucesso!
                    </span>
                  </div>
                </div>

                <div className="space-y-3">
                  <div className="flex items-center justify-between py-2 border-b">
                    <span className="text-sm text-gray-600">
                      ID da Transação:
                    </span>
                    <span className="text-sm font-mono font-medium text-gray-900">
                      #{result.id}
                    </span>
                  </div>

                  <div className="flex items-center justify-between py-2 border-b">
                    <span className="text-sm text-gray-600">Tipo:</span>
                    <span
                      className={`text-sm font-medium ${
                        result.tipo === "Compra"
                          ? "text-blue-600"
                          : "text-orange-600"
                      }`}
                    >
                      {result.tipo}
                    </span>
                  </div>

                  <div className="flex items-center justify-between py-2 border-b">
                    <span className="text-sm text-gray-600">Data:</span>
                    <span className="text-sm font-medium text-gray-900">
                      {new Date(result.data_transacao).toLocaleString("pt-BR")}
                    </span>
                  </div>

                  <div className="flex items-center justify-between py-2 border-b">
                    <span className="text-sm text-gray-600">Moedas:</span>
                    <span className="text-sm font-medium text-gray-900">
                      {result.moeda_origem} → {result.moeda_destino}
                    </span>
                  </div>

                  <div className="flex items-center justify-between py-2 border-b">
                    <span className="text-sm text-gray-600">Valor Origem:</span>
                    <span className="text-sm font-medium text-gray-900">
                      {result.valor_origem.toFixed(2)} {result.moeda_origem}
                    </span>
                  </div>

                  <div className="flex items-center justify-between py-2 border-b">
                    <span className="text-sm text-gray-600">
                      Taxa de Câmbio:
                    </span>
                    <span className="text-sm font-medium text-gray-900">
                      {result.taxa_cambio.toFixed(4)}
                    </span>
                  </div>

                  <div className="flex items-center justify-between py-3 bg-blue-50 rounded-lg px-3 mt-2">
                    <span className="text-sm font-medium text-blue-700">
                      Valor Final:
                    </span>
                    <span className="text-lg font-bold text-blue-900">
                      {result.valor_destino.toFixed(2)} {result.moeda_destino}
                    </span>
                  </div>

                  <div className="flex items-center justify-between py-2">
                    <span className="text-sm text-gray-600">Status:</span>
                    <span
                      className={`px-3 py-1 rounded-full text-xs font-medium ${
                        result.status === "Concluída"
                          ? "bg-green-100 text-green-700"
                          : result.status === "Pendente"
                          ? "bg-yellow-100 text-yellow-700"
                          : "bg-red-100 text-red-700"
                      }`}
                    >
                      {result.status}
                    </span>
                  </div>
                </div>

                <button
                  onClick={() => setResult(null)}
                  className="w-full px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors font-medium mt-4"
                >
                  Nova Operação
                </button>
              </div>
            ) : (
              <div className="flex flex-col items-center justify-center py-12 text-center">
                <div className="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mb-4">
                  <ArrowRightLeft className="w-8 h-8 text-gray-400" />
                </div>
                <p className="text-gray-500">
                  Preencha o formulário e confirme a operação para ver os
                  resultados
                </p>
              </div>
            )}
          </div>
        </div>

        <div className="grid md:grid-cols-3 gap-4 mt-6">
          <div className="bg-white rounded-lg shadow-md p-4">
            <div className="flex items-center space-x-3">
              <div className="w-10 h-10 bg-blue-100 rounded-full flex items-center justify-center">
                <ShoppingCart className="w-5 h-5 text-blue-600" />
              </div>
              <div>
                <p className="text-sm text-gray-600">Compra de Moeda</p>
                <p className="text-xs text-gray-500">
                  Adquirir moeda estrangeira
                </p>
              </div>
            </div>
          </div>

          <div className="bg-white rounded-lg shadow-md p-4">
            <div className="flex items-center space-x-3">
              <div className="w-10 h-10 bg-orange-100 rounded-full flex items-center justify-center">
                <TrendingDown className="w-5 h-5 text-orange-600" />
              </div>
              <div>
                <p className="text-sm text-gray-600">Venda de Moeda</p>
                <p className="text-xs text-gray-500">
                  Vender moeda estrangeira
                </p>
              </div>
            </div>
          </div>

          <div className="bg-white rounded-lg shadow-md p-4">
            <div className="flex items-center space-x-3">
              <div className="w-10 h-10 bg-green-100 rounded-full flex items-center justify-center">
                <CheckCircle className="w-5 h-5 text-green-600" />
              </div>
              <div>
                <p className="text-sm text-gray-600">Transação Segura</p>
                <p className="text-xs text-gray-500">Registro automático</p>
              </div>
            </div>
          </div>
        </div>

        {/* Exchange Rates Cards */}
        <div className="mt-6">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-xl font-bold text-gray-800 flex items-center">
              <TrendingUp className="w-5 h-5 mr-2 text-blue-600" />
              Cotações em Tempo Real
            </h2>
            {!loading.rates && Object.keys(exchangeRates).length > 0 && (
              <span className="text-sm text-gray-500">
                {Object.keys(exchangeRates).length} pares disponíveis
              </span>
            )}
          </div>

          {loading.rates ? (
            <div className="bg-white rounded-lg shadow-md p-12 flex items-center justify-center">
              <div className="text-center">
                <Loader2 className="w-12 h-12 text-blue-600 animate-spin mx-auto mb-3" />
                <p className="text-gray-600">Carregando cotações...</p>
              </div>
            </div>
          ) : Object.keys(exchangeRates).length === 0 ? (
            <div className="bg-white rounded-lg shadow-md p-12 text-center">
              <AlertCircle className="w-12 h-12 text-gray-400 mx-auto mb-3" />
              <p className="text-gray-500">
                Nenhuma cotação disponível no momento
              </p>
            </div>
          ) : (
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
              {Object.keys(exchangeRates)
                .sort()
                .map((rateKey) => {
                  const [from, to] = rateKey.split("_");
                  const rate = exchangeRates[rateKey];

                  return (
                    <div
                      key={rateKey}
                      className="bg-white rounded-lg shadow-md hover:shadow-lg transition-all duration-200 p-4 border border-gray-100 hover:border-blue-300 cursor-pointer group"
                      onClick={() => {
                        setOperation((prev) => ({
                          ...prev,
                          moedaOrigem: from,
                          moedaDestino: to,
                        }));
                        window.scrollTo({ top: 0, behavior: "smooth" });
                      }}
                    >
                      <div className="flex items-center justify-between mb-3">
                        <div className="flex items-center space-x-2">
                          <div className="w-8 h-8 bg-blue-100 rounded-full flex items-center justify-center group-hover:bg-blue-200 transition-colors">
                            <span className="text-xs font-bold text-blue-700">
                              {from}
                            </span>
                          </div>
                          <ArrowRightLeft className="w-4 h-4 text-gray-400 group-hover:text-blue-600 transition-colors" />
                          <div className="w-8 h-8 bg-green-100 rounded-full flex items-center justify-center group-hover:bg-green-200 transition-colors">
                            <span className="text-xs font-bold text-green-700">
                              {to}
                            </span>
                          </div>
                        </div>
                      </div>

                      <div className="space-y-1">
                        <p className="text-xs text-gray-500 uppercase tracking-wide">
                          Taxa de Câmbio
                        </p>
                        <p className="text-2xl font-bold text-gray-900 font-mono">
                          {rate.toFixed(4)}
                        </p>
                      </div>

                      <div className="mt-3 pt-3 border-t border-gray-100">
                        <p className="text-xs text-gray-600">
                          1 {from} = {rate.toFixed(4)} {to}
                        </p>
                      </div>

                      <div className="mt-2 opacity-0 group-hover:opacity-100 transition-opacity">
                        <p className="text-xs text-blue-600 font-medium">
                          Clique para usar esta cotação
                        </p>
                      </div>
                    </div>
                  );
                })}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default ExchangeRate;

import React, { useState, useEffect } from "react";
import axios from "axios";
import {
  BarChart3,
  Calendar,
  Filter,
  TrendingUp,
  CheckCircle,
  Clock,
  DollarSign,
  Download,
  FileText,
  Search,
  X,
  AlertCircle,
} from "lucide-react";

interface Transaction {
  id: number;
  date: string;
  time: string;
  type: string;
  fromCurrency: string;
  toCurrency: string;
  amount: string;
  convertedAmount: string;
  rate: string;
  status: string;
}

interface Filters {
  dateFrom: string;
  dateTo: string;
  currency: string;
  type: string;
}

const Extract: React.FC = () => {
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [filters, setFilters] = useState<Filters>({
    dateFrom: "",
    dateTo: "",
    currency: "",
    type: "",
  });

  // Buscar transações da API
  const fetchTransactions = async () => {
    try {
      setLoading(true);
      setError(null);

      // Construir query parameters
      const params = new URLSearchParams();
      if (filters.dateFrom) params.append("data_inicio", filters.dateFrom);
      if (filters.dateTo) params.append("data_fim", filters.dateTo);
      if (filters.type) params.append("tipo", filters.type);
      if (filters.currency) {
        params.append("moeda_origem", filters.currency);
        // Também buscar por moeda destino
      }

      const response = await axios.get(`/api/transacoes?${params.toString()}`);

      // Transformar dados da API para o formato do componente
      const apiTransactions = response.data.transactions || [];
      const formattedTransactions = apiTransactions.map((t: any) => ({
        id: t.id,
        date: formatDateFromISO(t.data_transacao),
        time: formatTimeFromISO(t.data_transacao),
        type: t.tipo,
        fromCurrency: t.moeda_origem,
        toCurrency: t.moeda_destino,
        amount: t.valor_origem.toFixed(2),
        convertedAmount: t.valor_destino.toFixed(2),
        rate: t.taxa_cambio.toFixed(4),
        status: t.status,
      }));

      setTransactions(formattedTransactions);
    } catch (err: any) {
      console.error("Erro ao buscar transações:", err);
      setError(
        err.response?.status === 503
          ? "Banco de dados não configurado. Usando modo de demonstração."
          : "Erro ao carregar transações. Usando dados de demonstração."
      );
      // Em caso de erro, usar dados simulados
      setTransactions(generateSampleTransactions());
    } finally {
      setLoading(false);
    }
  };

  // Formatar data ISO para d/M/Y
  const formatDateFromISO = (isoDate: string): string => {
    const date = new Date(isoDate);
    return `${date.getDate()}/${date.getMonth() + 1}/${date.getFullYear()}`;
  };

  // Formatar hora de ISO
  const formatTimeFromISO = (isoDate: string): string => {
    const date = new Date(isoDate);
    return `${String(date.getHours()).padStart(2, "0")}:${String(
      date.getMinutes()
    ).padStart(2, "0")}`;
  };

  // Simulated transaction data (since we don't have a backend endpoint for this)
  const generateSampleTransactions = (): Transaction[] => {
    const currencies = ["USD", "EUR", "GBP", "JPY", "BRL"];
    const types = ["Compra", "Venda"];
    const sampleTransactions: Transaction[] = [];

    for (let i = 0; i < 20; i++) {
      const date = new Date();
      date.setDate(date.getDate() - Math.floor(Math.random() * 30));

      // Format date as d/M/Y
      const formattedDate = `${date.getDate()}/${
        date.getMonth() + 1
      }/${date.getFullYear()}`;

      sampleTransactions.push({
        id: i + 1,
        date: formattedDate,
        time: `${String(Math.floor(Math.random() * 24)).padStart(
          2,
          "0"
        )}:${String(Math.floor(Math.random() * 60)).padStart(2, "0")}`,
        type: types[Math.floor(Math.random() * types.length)],
        fromCurrency: currencies[Math.floor(Math.random() * currencies.length)],
        toCurrency: currencies[Math.floor(Math.random() * currencies.length)],
        amount: (Math.random() * 10000 + 100).toFixed(2),
        convertedAmount: (Math.random() * 10000 + 100).toFixed(2),
        rate: (Math.random() * 10 + 0.1).toFixed(4),
        status: Math.random() > 0.1 ? "Concluído" : "Pendente",
      });
    }

    return sampleTransactions.sort((a, b) => {
      const dateA = new Date(a.date.split("/").reverse().join("-"));
      const dateB = new Date(b.date.split("/").reverse().join("-"));
      return dateB.getTime() - dateA.getTime();
    });
  };

  useEffect(() => {
    fetchTransactions();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [filters]); // Recarregar quando os filtros mudarem

  const filteredTransactions = transactions.filter((transaction) => {
    // Convert d/M/Y format to Date for comparison
    const transactionDate = new Date(
      transaction.date.split("/").reverse().join("-")
    );

    if (filters.dateFrom) {
      const fromDate = new Date(filters.dateFrom);
      if (transactionDate < fromDate) return false;
    }

    if (filters.dateTo) {
      const toDate = new Date(filters.dateTo);
      if (transactionDate > toDate) return false;
    }

    if (
      filters.currency &&
      !transaction.fromCurrency.includes(filters.currency) &&
      !transaction.toCurrency.includes(filters.currency)
    )
      return false;
    if (filters.type && transaction.type !== filters.type) return false;
    return true;
  });

  const generatePDF = () => {
    alert("Funcionalidade de exportação PDF será implementada em breve!");
  };

  const exportCSV = () => {
    const csvContent = [
      [
        "Data",
        "Hora",
        "Tipo",
        "Moeda Origem",
        "Moeda Destino",
        "Valor",
        "Valor Convertido",
        "Taxa",
        "Status",
      ],
      ...filteredTransactions.map((t) => [
        t.date,
        t.time,
        t.type,
        t.fromCurrency,
        t.toCurrency,
        t.amount,
        t.convertedAmount,
        t.rate,
        t.status,
      ]),
    ]
      .map((row) => row.join(","))
      .join("\n");

    const blob = new Blob([csvContent], { type: "text/csv" });
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = `extrato-cambio-${new Date().toISOString().split("T")[0]}.csv`;
    a.click();
    window.URL.revokeObjectURL(url);
  };

  const clearFilters = () => {
    setFilters({ dateFrom: "", dateTo: "", currency: "", type: "" });
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 to-blue-50 p-6">
      <div className="max-w-7xl mx-auto space-y-6">
        {/* Header */}
        <div className="bg-white rounded-lg p-8 shadow-lg border border-blue-100">
          <div className="flex items-center gap-4">
            <BarChart3 className="w-8 h-8 text-blue-600" />
            <div className="flex-1">
              <h1 className="text-3xl font-bold text-gray-900">
                Extrato de Operações
              </h1>
              <p className="text-gray-600 mt-1">
                Histórico completo de transações cambiais
              </p>
            </div>
            {loading && (
              <div className="flex items-center gap-2 text-blue-600">
                <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-blue-600"></div>
                <span className="text-sm">Carregando...</span>
              </div>
            )}
          </div>
          {error && (
            <div className="mt-4 p-3 bg-yellow-50 border border-yellow-200 rounded-lg flex items-start gap-2">
              <AlertCircle className="w-5 h-5 text-yellow-600 flex-shrink-0 mt-0.5" />
              <div className="text-sm text-yellow-800">{error}</div>
            </div>
          )}
        </div>

        {/* Filters */}
        <div className="bg-white rounded-lg p-8 shadow-lg border border-blue-100">
          <div className="flex items-center gap-3 mb-6">
            <Filter className="w-6 h-6 text-blue-600" />
            <h3 className="text-xl font-semibold text-gray-900">Filtros</h3>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                <Calendar className="w-4 h-4 inline mr-1" />
                Data Inicial
              </label>
              <input
                type="date"
                value={filters.dateFrom}
                onChange={(e) =>
                  setFilters((prev) => ({ ...prev, dateFrom: e.target.value }))
                }
                className="input-field"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                <Calendar className="w-4 h-4 inline mr-1" />
                Data Final
              </label>
              <input
                type="date"
                value={filters.dateTo}
                onChange={(e) =>
                  setFilters((prev) => ({ ...prev, dateTo: e.target.value }))
                }
                className="input-field"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                <DollarSign className="w-4 h-4 inline mr-1" />
                Moeda
              </label>
              <select
                value={filters.currency}
                onChange={(e) =>
                  setFilters((prev) => ({ ...prev, currency: e.target.value }))
                }
                className="select-field"
              >
                <option value="">Todas</option>
                <option value="USD">USD</option>
                <option value="EUR">EUR</option>
                <option value="GBP">GBP</option>
                <option value="JPY">JPY</option>
                <option value="BRL">BRL</option>
              </select>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                <TrendingUp className="w-4 h-4 inline mr-1" />
                Tipo
              </label>
              <select
                value={filters.type}
                onChange={(e) =>
                  setFilters((prev) => ({ ...prev, type: e.target.value }))
                }
                className="select-field"
              >
                <option value="">Todos</option>
                <option value="Compra">Compra</option>
                <option value="Venda">Venda</option>
              </select>
            </div>
          </div>

          <div className="flex justify-end mt-6">
            <button
              onClick={clearFilters}
              className="flex items-center gap-2 px-4 py-2 text-gray-600 bg-gray-100 hover:bg-gray-200 rounded transition-colors duration-200"
            >
              <X className="w-4 h-4" />
              Limpar Filtros
            </button>
          </div>
        </div>

        {/* Summary Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <div className="bg-white rounded-lg p-6 shadow-lg border border-blue-100">
            <div className="flex items-center gap-3">
              <TrendingUp className="w-8 h-8 text-blue-600" />
              <div>
                <h4 className="text-sm font-medium text-gray-600">
                  Total de Operações
                </h4>
                <span className="text-2xl font-bold text-blue-600">
                  {filteredTransactions.length}
                </span>
              </div>
            </div>
          </div>

          <div className="bg-white rounded-lg p-6 shadow-lg border border-blue-100">
            <div className="flex items-center gap-3">
              <CheckCircle className="w-8 h-8 text-blue-600" />
              <div>
                <h4 className="text-sm font-medium text-gray-600">
                  Concluídas
                </h4>
                <span className="text-2xl font-bold text-blue-600">
                  {
                    filteredTransactions.filter((t) => t.status === "Concluído")
                      .length
                  }
                </span>
              </div>
            </div>
          </div>

          <div className="bg-white rounded-lg p-6 shadow-lg border border-blue-100">
            <div className="flex items-center gap-3">
              <Clock className="w-8 h-8 text-blue-600" />
              <div>
                <h4 className="text-sm font-medium text-gray-600">Pendentes</h4>
                <span className="text-2xl font-bold text-blue-600">
                  {
                    filteredTransactions.filter((t) => t.status === "Pendente")
                      .length
                  }
                </span>
              </div>
            </div>
          </div>

          <div className="bg-white rounded-lg p-6 shadow-lg border border-blue-100">
            <div className="flex items-center gap-3">
              <DollarSign className="w-8 h-8 text-blue-600" />
              <div>
                <h4 className="text-sm font-medium text-gray-600">
                  Volume Total
                </h4>
                <span className="text-2xl font-bold text-blue-600">
                  $
                  {filteredTransactions
                    .reduce((sum, t) => sum + parseFloat(t.amount), 0)
                    .toLocaleString("pt-BR", { minimumFractionDigits: 2 })}
                </span>
              </div>
            </div>
          </div>
        </div>

        {/* Export Actions */}
        <div className="flex gap-4 justify-end">
          <button
            onClick={generatePDF}
            className="btn-primary flex items-center justify-center gap-2 rounded"
          >
            <FileText className="w-4 h-4" />
            Exportar PDF
          </button>
          <button
            onClick={exportCSV}
            className="btn-primary flex items-center justify-center gap-2 rounded"
          >
            <Download className="w-4 h-4" />
            Exportar CSV
          </button>
        </div>

        {/* Transactions Table */}
        <div className="bg-white rounded-lg shadow-lg border border-blue-100 overflow-hidden">
          <div className="p-8 border-b border-gray-200">
            <div className="flex items-center justify-between">
              <h3 className="text-xl font-semibold text-gray-900">
                Histórico de Transações
              </h3>
              <span className="text-sm text-gray-600">
                {filteredTransactions.length} transação(ões) encontrada(s)
              </span>
            </div>
          </div>

          <div className="overflow-x-auto">
            <div className="min-w-full">
              {/* Table Header */}
              <div className="bg-gray-50 grid grid-cols-6 gap-4 p-4 font-semibold text-gray-700 text-sm">
                <div>Data/Hora</div>
                <div>Tipo</div>
                <div>Conversão</div>
                <div>Valores</div>
                <div>Taxa</div>
                <div>Status</div>
              </div>

              {/* Table Body */}
              <div className="divide-y divide-gray-200">
                {filteredTransactions.map((transaction) => (
                  <div
                    key={transaction.id}
                    className="grid grid-cols-6 gap-4 p-4 hover:bg-gray-50 transition-colors"
                  >
                    <div>
                      <div className="font-medium text-gray-900">
                        {transaction.date}
                      </div>
                      <div className="text-sm text-gray-500">
                        {transaction.time}
                      </div>
                    </div>

                    <div>
                      <span
                        className={`inline-flex px-3 py-1 rounded text-xs font-medium ${
                          transaction.type === "Compra"
                            ? "bg-green-100 text-green-800"
                            : "bg-red-100 text-red-800"
                        }`}
                      >
                        {transaction.type}
                      </span>
                    </div>

                    <div className="font-mono text-sm">
                      {transaction.fromCurrency} → {transaction.toCurrency}
                    </div>

                    <div>
                      <div className="font-medium text-gray-900">
                        {transaction.fromCurrency}{" "}
                        {parseFloat(transaction.amount).toLocaleString(
                          "pt-BR",
                          { minimumFractionDigits: 2 }
                        )}
                      </div>
                      <div className="text-sm text-gray-600">
                        {transaction.toCurrency}{" "}
                        {parseFloat(transaction.convertedAmount).toLocaleString(
                          "pt-BR",
                          { minimumFractionDigits: 2 }
                        )}
                      </div>
                    </div>

                    <div className="font-mono text-sm font-medium">
                      {transaction.rate}
                    </div>

                    <div>
                      <span
                        className={`inline-flex px-3 py-1 rounded text-xs font-medium ${
                          transaction.status === "Concluído"
                            ? "bg-blue-100 text-blue-800"
                            : "bg-gray-100 text-gray-800"
                        }`}
                      >
                        {transaction.status}
                      </span>
                    </div>
                  </div>
                ))}

                {filteredTransactions.length === 0 && (
                  <div className="text-center py-16">
                    <Search className="w-12 h-12 text-gray-400 mx-auto mb-4" />
                    <div className="text-xl font-medium text-gray-900 mb-2">
                      Nenhuma transação encontrada
                    </div>
                    <div className="text-gray-600">
                      Ajuste os filtros ou realize uma nova operação
                    </div>
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Extract;

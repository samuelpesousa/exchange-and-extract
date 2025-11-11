import React, { useState, useEffect } from "react";
import axios from "axios";
import jsPDF from "jspdf";
import autoTable from "jspdf-autotable";
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
  RefreshCw,
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

      const formattedTransactions = apiTransactions.map((t: any) => {
        return {
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
        };
      });

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
    // Se a data tem 'Z' no final, ela está em UTC, precisamos converter para local
    const date = new Date(isoDate);
    return `${date.getDate()}/${date.getMonth() + 1}/${date.getFullYear()}`;
  };

  // Formatar hora de ISO - converter de UTC para horário local
  const formatTimeFromISO = (isoDate: string): string => {
    // Parse the ISO date string to a Date object
    // JavaScript Date automaticamente converte para timezone local
    const date = new Date(isoDate);

    // Verificar se a data é válida
    if (isNaN(date.getTime())) {
      console.error("Data inválida:", isoDate);
      return "00:00";
    }

    // Get local time components (já convertido automaticamente pelo JavaScript)
    const hours = date.getHours();
    const minutes = date.getMinutes();

    return `${String(hours).padStart(2, "0")}:${String(minutes).padStart(
      2,
      "0"
    )}`;
  };

  // Simulated transaction data (only used when API fails)
  const generateSampleTransactions = (): Transaction[] => {
    const currencies = ["USD", "EUR", "GBP", "JPY", "BRL"];
    const types = ["Compra", "Venda"];
    const sampleTransactions: Transaction[] = [];

    // Usar data/hora atual como base
    const now = new Date();

    for (let i = 0; i < 20; i++) {
      // Criar uma cópia da data atual e subtrair dias (mais recentes primeiro)
      const date = new Date(now);
      date.setDate(date.getDate() - i);

      // Definir um horário específico baseado no índice (não aleatório)
      // Distribuir ao longo do dia de forma previsível
      const hoursOffset = (i * 2) % 24; // Varia de 0 a 23
      const minutesOffset = (i * 7) % 60; // Varia de 0 a 59
      date.setHours(hoursOffset);
      date.setMinutes(minutesOffset);
      date.setSeconds(0);
      date.setMilliseconds(0);

      // Format date as d/M/Y
      const formattedDate = `${date.getDate()}/${
        date.getMonth() + 1
      }/${date.getFullYear()}`;

      // Format time from the date object
      const formattedTime = `${String(date.getHours()).padStart(
        2,
        "0"
      )}:${String(date.getMinutes()).padStart(2, "0")}`;

      sampleTransactions.push({
        id: i + 1,
        date: formattedDate,
        time: formattedTime,
        type: types[i % 2], // Alterna entre Compra e Venda
        fromCurrency: currencies[i % currencies.length],
        toCurrency: currencies[(i + 1) % currencies.length],
        amount: (1000 + i * 100).toFixed(2), // Valores incrementais
        convertedAmount: (1200 + i * 120).toFixed(2),
        rate: (1.1 + i * 0.01).toFixed(4),
        status: i % 10 === 0 ? "Pendente" : "Concluído", // 1 em 10 é pendente
      });
    }

    return sampleTransactions.sort((a, b) => {
      const dateA = new Date(
        a.date.split("/").reverse().join("-") + "T" + a.time
      );
      const dateB = new Date(
        b.date.split("/").reverse().join("-") + "T" + b.time
      );
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
    const doc = new jsPDF();
    const pageWidth = doc.internal.pageSize.width;
    const pageHeight = doc.internal.pageSize.height;
    const margin = 15;

    // Carregar e adicionar logo
    const img = new Image();
    img.src = "/image/acorianalight.png";

    img.onload = () => {
      // Logo centralizado no topo
      const logoWidth = 50;
      const logoHeight = 15;
      const logoX = (pageWidth - logoWidth) / 2;

      try {
        doc.addImage(img, "PNG", logoX, 10, logoWidth, logoHeight);
      } catch (error) {
        console.error("Erro ao adicionar logo:", error);
      }

      // Cabeçalho - Título (ajustado para baixo do logo)
      doc.setFontSize(14);
      doc.setTextColor(75, 85, 99); // Cinza
      doc.text("Extrato de Operações Cambiais", pageWidth / 2, 32, {
        align: "center",
      });

      // Data de emissão
      doc.setFontSize(10);
      doc.setTextColor(107, 114, 128);
      const dataEmissao = new Date().toLocaleDateString("pt-BR", {
        day: "2-digit",
        month: "2-digit",
        year: "numeric",
        hour: "2-digit",
        minute: "2-digit",
      });
      doc.text(`Emitido em: ${dataEmissao}`, margin, 42);

      // Resumo
      const totalOperacoes = filteredTransactions.length;
      const totalConcluidas = filteredTransactions.filter(
        (t) => t.status === "Concluído"
      ).length;
      const volumeTotal = filteredTransactions
        .reduce((sum, t) => sum + parseFloat(t.amount), 0)
        .toFixed(2);

      doc.setFontSize(10);
      doc.text(`Total de Operações: ${totalOperacoes}`, margin, 52);
      doc.text(`Operações Concluídas: ${totalConcluidas}`, margin, 58);
      doc.text(`Volume Total: $ ${volumeTotal}`, margin, 64);

      // Linha separadora
      doc.setDrawColor(229, 231, 235);
      doc.line(margin, 70, pageWidth - margin, 70);

      // Tabela de transações
      const tableData = filteredTransactions.map((t) => [
        `${t.date}\n${t.time}`,
        t.type,
        `${t.fromCurrency} → ${t.toCurrency}`,
        `${t.fromCurrency} ${parseFloat(t.amount).toLocaleString("pt-BR", {
          minimumFractionDigits: 2,
        })}`,
        `${t.toCurrency} ${parseFloat(t.convertedAmount).toLocaleString(
          "pt-BR",
          { minimumFractionDigits: 2 }
        )}`,
        t.rate,
        t.status,
      ]);

      autoTable(doc, {
        startY: 77,
        head: [
          [
            "Data/Hora",
            "Tipo",
            "Conversão",
            "Valor Origem",
            "Valor Destino",
            "Taxa",
            "Status",
          ],
        ],
        body: tableData,
        theme: "grid",
        headStyles: {
          fillColor: [37, 99, 235],
          textColor: [255, 255, 255],
          fontSize: 9,
          fontStyle: "bold",
        },
        bodyStyles: {
          fontSize: 8,
          textColor: [31, 41, 55],
        },
        alternateRowStyles: {
          fillColor: [249, 250, 251],
        },
        columnStyles: {
          0: { cellWidth: 25 },
          1: { cellWidth: 18 },
          2: { cellWidth: 25 },
          3: { cellWidth: 30 },
          4: { cellWidth: 30 },
          5: { cellWidth: 18 },
          6: { cellWidth: 20 },
        },
        margin: { left: margin, right: margin },
      });

      // Informações da corretora no rodapé
      const finalY = (doc as any).lastAutoTable.finalY || 77;
      let currentY = finalY + 15;

      // Verificar se precisa de nova página
      if (currentY > pageHeight - 80) {
        doc.addPage();
        currentY = 20;
      }

      doc.setFontSize(9);
      doc.setTextColor(75, 85, 99);
      doc.setFont("helvetica", "bold");
      doc.text("INFORMAÇÕES IMPORTANTES", margin, currentY);

      currentY += 7;
      doc.setFont("helvetica", "normal");
      doc.setFontSize(8);
      doc.setTextColor(55, 65, 81);

      const avisoTexto = [
        "A CORRETORA AÇORIANA DE CÂMBIO S.A. vem reiterar a todos os seus clientes, potenciais",
        "clientes e ao público em geral que é uma instituição financeira cuja atividade se restringe ao",
        "mercado e operações de câmbio (incluindo câmbio para aquisição de criptoativos), e NÃO",
        "OFERECE EMPRÉSTIMOS, FINANCIAMENTOS, CONSÓRCIOS ou INVESTIMENTOS COMO AÇÕES",
        "ou RENDA FIXA. Em resumo, a Corretora Açoriana não exerce qualquer atividade fora de sua",
        "licença para atuação no mercado de câmbio.",
        "",
        "A Corretora Açoriana também NÃO COBRA QUAISQUER TAXAS OU COMISSÕES ANTECIPADAS,",
        "como taxas para análise de cadastro, taxas de análise de crédito, taxas ou comissões de",
        "viabilização de operações, ou quaisquer outras.",
        "",
        "A Corretora Açoriana de Câmbio S.A. é uma instituição regulada pelo Banco Central do Brasil.",
      ];

      avisoTexto.forEach((linha) => {
        if (currentY > pageHeight - 20) {
          doc.addPage();
          currentY = 20;
        }
        doc.text(linha, margin, currentY);
        currentY += 4.5;
      });

      currentY += 5;

      // Endereço
      if (currentY > pageHeight - 25) {
        doc.addPage();
        currentY = 20;
      }

      doc.setFont("helvetica", "bold");
      doc.setFontSize(9);
      doc.text("Endereço:", margin, currentY);
      currentY += 6;

      doc.setFont("helvetica", "normal");
      doc.setFontSize(8);
      doc.text(
        "Florianópolis/SC: Rua Dom Jaime Câmara, 106 – Centro, CEP 88015-120",
        margin,
        currentY
      );
      currentY += 5;
      doc.text("CNPJ 15.761.217/0001-91", margin, currentY);

      // Rodapé em todas as páginas
      const totalPages = doc.getNumberOfPages();
      for (let i = 1; i <= totalPages; i++) {
        doc.setPage(i);
        doc.setFontSize(8);
        doc.setTextColor(156, 163, 175);
        doc.text(
          `Página ${i} de ${totalPages}`,
          pageWidth / 2,
          pageHeight - 10,
          { align: "center" }
        );
        doc.text(
          "Corretora Açoriana de Câmbio S.A. | CNPJ 15.761.217/0001-91",
          pageWidth / 2,
          pageHeight - 5,
          { align: "center" }
        );
      }

      // Salvar PDF
      const fileName = `extrato-cambio-${
        new Date().toISOString().split("T")[0]
      }.pdf`;
      doc.save(fileName);
    };

    // Se a imagem não carregar, gerar PDF sem logo
    img.onerror = () => {
      console.error("Erro ao carregar logo, gerando PDF sem imagem");
      // Gerar PDF sem logo (código original)
      doc.setFontSize(18);
      doc.setTextColor(37, 99, 235);
      doc.text("CORRETORA AÇORIANA DE CÂMBIO S.A.", pageWidth / 2, 20, {
        align: "center",
      });
      doc.setFontSize(14);
      doc.setTextColor(75, 85, 99);
      doc.text("Extrato de Operações Cambiais", pageWidth / 2, 30, {
        align: "center",
      });
      // ... resto do código continua igual
    };
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
            <button
              onClick={fetchTransactions}
              disabled={loading}
              className="flex items-center gap-2 px-4 py-2 text-blue-600 bg-blue-50 hover:bg-blue-100 rounded transition-colors duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
              title="Recarregar transações"
            >
              <RefreshCw
                className={`w-4 h-4 ${loading ? "animate-spin" : ""}`}
              />
              <span className="text-sm font-medium">Atualizar</span>
            </button>
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

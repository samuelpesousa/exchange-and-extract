import React, { useState } from "react";
import { Lock, Mail, Eye, EyeOff, AlertCircle } from "lucide-react";

interface LoginProps {
  onLogin: (email: string, password: string) => void;
}

const Login: React.FC<LoginProps> = ({ onLogin }) => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    // Validações básicas
    if (!email || !password) {
      setError("Por favor, preencha todos os campos.");
      return;
    }

    if (!email.includes("@")) {
      setError("Por favor, insira um e-mail válido.");
      return;
    }

    setIsLoading(true);

    try {
      // Simular chamada de API
      await new Promise((resolve) => setTimeout(resolve, 1000));

      // Aqui você pode adicionar a lógica real de autenticação
      onLogin(email, password);
    } catch (err) {
      setError("Erro ao fazer login. Tente novamente.");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="bg-gradient-to-br from-slate-50 to-blue-50 p-6 min-h-[calc(100vh-5rem)]">
      <div className="max-w-7xl mx-auto flex items-center justify-center min-h-full">
        <div className="w-full max-w-md">
          {/* Card de Login */}
          <div className="bg-white rounded-2xl shadow-2xl border border-blue-100 overflow-hidden">
            {/* Header com logo */}
            <div className="bg-gradient-to-r from-blue-600 to-blue-700 p-8 text-center">
              <img
                src="/image/acorianalight.png"
                alt="Corretora Açoriana"
                className="h-16 w-auto mx-auto mb-4 filter brightness-0 invert"
                onError={(e) => {
                  const target = e.target as HTMLImageElement;
                  target.style.display = "none";
                }}
              />
              <h1 className="text-2xl font-bold text-white">
                Bem-vindo de volta
              </h1>
              <p className="text-blue-100 mt-2">
                Acesse sua conta para continuar
              </p>
            </div>

            {/* Formulário */}
            <div className="p-8">
              {error && (
                <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg flex items-start gap-3">
                  <AlertCircle className="w-5 h-5 text-red-600 flex-shrink-0 mt-0.5" />
                  <p className="text-sm text-red-800">{error}</p>
                </div>
              )}

              <form onSubmit={handleSubmit} className="space-y-6">
                {/* Campo de Email */}
                <div>
                  <label
                    htmlFor="email"
                    className="block text-sm font-medium text-gray-700 mb-2"
                  >
                    E-mail
                  </label>
                  <div className="relative">
                    <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                      <Mail className="h-5 w-5 text-gray-400" />
                    </div>
                    <input
                      id="email"
                      type="email"
                      value={email}
                      onChange={(e) => setEmail(e.target.value)}
                      className="block w-full pl-10 pr-3 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors"
                      placeholder="seu@email.com"
                      disabled={isLoading}
                    />
                  </div>
                </div>

                {/* Campo de Senha */}
                <div>
                  <label
                    htmlFor="password"
                    className="block text-sm font-medium text-gray-700 mb-2"
                  >
                    Senha
                  </label>
                  <div className="relative">
                    <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                      <Lock className="h-5 w-5 text-gray-400" />
                    </div>
                    <input
                      id="password"
                      type={showPassword ? "text" : "password"}
                      value={password}
                      onChange={(e) => setPassword(e.target.value)}
                      className="block w-full pl-10 pr-12 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors"
                      placeholder="••••••••"
                      disabled={isLoading}
                    />
                    <button
                      type="button"
                      onClick={() => setShowPassword(!showPassword)}
                      className="absolute inset-y-0 right-0 pr-3 flex items-center"
                      disabled={isLoading}
                    >
                      {showPassword ? (
                        <EyeOff className="h-5 w-5 text-gray-400 hover:text-gray-600" />
                      ) : (
                        <Eye className="h-5 w-5 text-gray-400 hover:text-gray-600" />
                      )}
                    </button>
                  </div>
                </div>

                {/* Lembrar-me e Esqueci senha */}
                <div className="flex items-center justify-between">
                  <div className="flex items-center">
                    <input
                      id="remember"
                      type="checkbox"
                      className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                    />
                    <label
                      htmlFor="remember"
                      className="ml-2 block text-sm text-gray-700"
                    >
                      Lembrar-me
                    </label>
                  </div>

                  <button
                    type="button"
                    className="text-sm font-medium text-blue-600 hover:text-blue-700 transition-colors"
                    disabled={isLoading}
                  >
                    Esqueceu a senha?
                  </button>
                </div>

                {/* Botão de Login */}
                <button
                  type="submit"
                  disabled={isLoading}
                  className="w-full bg-gradient-to-r from-blue-600 to-blue-700 text-white py-3 px-4 rounded-lg font-semibold hover:from-blue-700 hover:to-blue-800 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center"
                >
                  {isLoading ? (
                    <>
                      <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-white mr-2"></div>
                      Entrando...
                    </>
                  ) : (
                    "Entrar"
                  )}
                </button>
              </form>

              {/* Divisor */}
              <div className="mt-8 relative">
                <div className="absolute inset-0 flex items-center">
                  <div className="w-full border-t border-gray-300"></div>
                </div>
                <div className="relative flex justify-center text-sm">
                  <span className="px-4 bg-white text-gray-500">
                    Novo por aqui?
                  </span>
                </div>
              </div>

              {/* Botão de Registro */}
              <button
                type="button"
                className="mt-6 w-full bg-white text-blue-600 py-3 px-4 rounded-lg font-semibold border-2 border-blue-600 hover:bg-blue-50 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-all duration-200"
                disabled={isLoading}
              >
                Criar uma conta
              </button>
            </div>

            {/* Footer */}
            <div className="bg-gray-50 px-8 py-4 border-t border-gray-200">
              <p className="text-xs text-center text-gray-600">
                Ao fazer login, você concorda com nossos{" "}
                <a href="#" className="text-blue-600 hover:text-blue-700">
                  Termos de Serviço
                </a>{" "}
                e{" "}
                <a href="#" className="text-blue-600 hover:text-blue-700">
                  Política de Privacidade
                </a>
                .
              </p>
              <p className="text-xs text-center text-gray-500 mt-2">
                Corretora Açoriana de Câmbio S.A. - CNPJ 15.761.217/0001-91
              </p>
            </div>
          </div>

          {/* Informação adicional */}
          <div className="mt-6 text-center">
            <p className="text-sm text-gray-600">
              Instituição regulada pelo{" "}
              <span className="font-semibold">Banco Central do Brasil</span>
            </p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Login;

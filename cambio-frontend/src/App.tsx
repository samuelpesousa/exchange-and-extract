import React, { useState, useEffect } from "react";
import Navbar from "./components/Navbar";
import ExchangeRate from "./components/ExchangeRate";
import Extract from "./components/Extract";
import LoginSimple from "./components/LoginSimple";
import Register from "./components/Register";
import authService from "./services/authService";
import "./App.css";

const App: React.FC = () => {
  const [activeTab, setActiveTab] = useState<string>("exchange");
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);
  const [showRegister, setShowRegister] = useState<boolean>(false);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [currentUser, setCurrentUser] = useState<any>(null);

  // Verificar se já está autenticado ao carregar o app
  useEffect(() => {
    const checkAuth = async () => {
      const isAuth = authService.isAuthenticated();
      console.log("Verificando autenticação:", isAuth);

      if (isAuth) {
        // Buscar dados do usuário
        const userData = authService.getCurrentUser();
        setCurrentUser(userData);
      }

      setIsAuthenticated(isAuth);
      setIsLoading(false);
    };

    checkAuth();
  }, []);

  const handleLogin = (email: string, password: string) => {
    console.log("Login bem-sucedido:", email);
    const userData = authService.getCurrentUser();
    setCurrentUser(userData);
    setIsAuthenticated(true);
    setActiveTab("exchange");
    setShowRegister(false);
  };

  const handleLogout = async () => {
    console.log("=== LOGOUT CLICKED ===");
    await authService.logout();
    setIsAuthenticated(false);
    setActiveTab("login");
    setShowRegister(false);
  };

  const handleRegisterSuccess = () => {
    console.log("Registro bem-sucedido");
    setShowRegister(false);
  };

  const handleShowRegister = () => {
    setShowRegister(true);
  };

  const handleBackToLogin = () => {
    setShowRegister(false);
  };

  const renderActiveTab = () => {
    switch (activeTab) {
      case "exchange":
        return <ExchangeRate />;
      case "extract":
        return <Extract />;
      case "login":
        return <LoginSimple onLogin={handleLogin} />;
      default:
        return <ExchangeRate />;
    }
  };

  // Tela de carregamento
  if (isLoading) {
    return (
      <div className="App">
        <div className="min-h-screen bg-gradient-to-br from-slate-50 to-blue-50 flex items-center justify-center">
          <div className="text-center">
            <div className="animate-spin rounded-full h-16 w-16 border-b-2 border-blue-900 mx-auto mb-4"></div>
            <p className="text-gray-600">Carregando...</p>
          </div>
        </div>
      </div>
    );
  }

  // Se não estiver autenticado, mostra login ou registro
  if (!isAuthenticated) {
    return (
      <div className="App">
        {showRegister ? (
          <Register
            onRegisterSuccess={handleRegisterSuccess}
            onBackToLogin={handleBackToLogin}
          />
        ) : (
          <LoginSimple
            onLogin={handleLogin}
            onShowRegister={handleShowRegister}
          />
        )}
      </div>
    );
  }

  return (
    <div className="App">
      {activeTab !== "login" && (
        <Navbar
          activeTab={activeTab}
          setActiveTab={setActiveTab}
          onLogout={handleLogout}
          currentUser={currentUser}
        />
      )}
      <main className="main-content">{renderActiveTab()}</main>
    </div>
  );
};

export default App;

import React, { useState, useRef, useEffect } from "react";
import {
  TrendingUp,
  FileText,
  User,
  Settings,
  LogOut,
  ChevronDown,
} from "lucide-react";
import logoImage from "../image/acorianalight.png";

interface NavbarProps {
  activeTab: string;
  setActiveTab: (tab: string) => void;
  onLogout?: () => void;
  currentUser?: {
    nome: string;
    email: string;
  } | null;
}

interface Tab {
  id: string;
  label: string;
  icon: React.ComponentType<{ className?: string }>;
}

const Navbar: React.FC<NavbarProps> = ({
  activeTab,
  setActiveTab,
  onLogout,
  currentUser,
}) => {
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);
  const dropdownRefMobile = useRef<HTMLDivElement>(null);

  const tabs: Tab[] = [
    { id: "exchange", label: "Cotação e Câmbio", icon: TrendingUp },
    { id: "extract", label: "Extrato", icon: FileText },
  ];

  // Função para pegar iniciais do nome
  const getInitials = (name: string) => {
    if (!name) return "U";
    const names = name.trim().split(" ");
    if (names.length === 1) return names[0][0].toUpperCase();
    return (names[0][0] + names[names.length - 1][0]).toUpperCase();
  };

  // Função para gerar cor do avatar (sempre azul)
  const getAvatarColor = () => {
    return "from-blue-500 to-blue-600";
  };

  const userName = currentUser?.nome || "Usuário";
  const userEmail = currentUser?.email || "usuario@email.com";
  const userInitials = getInitials(userName);
  const avatarGradient = getAvatarColor();

  // Fechar dropdown ao clicar fora
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        dropdownRef.current &&
        !dropdownRef.current.contains(event.target as Node) &&
        dropdownRefMobile.current &&
        !dropdownRefMobile.current.contains(event.target as Node)
      ) {
        setIsDropdownOpen(false);
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, []);

  const handleProfile = () => {
    // Implementar navegação para perfil
    console.log("Ver perfil");
    setIsDropdownOpen(false);
  };

  const handleEditProfile = () => {
    // Implementar navegação para editar perfil
    console.log("Editar perfil");
    setIsDropdownOpen(false);
  };

  const handleLogoutClick = () => {
    // Implementar lógica de logout
    console.log("=== handleLogoutClick chamado ===");
    console.log("onLogout existe?", !!onLogout);
    setIsDropdownOpen(false);
    if (onLogout) {
      console.log("Chamando onLogout do App");
      onLogout();
    } else {
      console.log("ERRO: onLogout não está definido!");
    }
  };

  return (
    <nav className="bg-white shadow-lg border-b-2 border-blue-100">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-20">
          {/* Logo e título */}
          <div className="flex items-center space-x-3">
            <img
              src="/image/acorianalight.png"
              alt="Acoriana Logo"
              className="h-16 w-auto object-contain"
              onError={(e) => {
                const target = e.target as HTMLImageElement;
                target.src = logoImage;
              }}
            />
          </div>

          {/* Navegação Desktop */}
          <div className="hidden md:flex items-center space-x-4">
            <div className="flex space-x-1 p-1 rounded-xl">
              {tabs.map((tab) => {
                const IconComponent = tab.icon;
                const isActive = activeTab === tab.id;

                return (
                  <button
                    key={tab.id}
                    className={`
                      relative flex items-center space-x-2 px-6 py-3 rounded-lg transition-all duration-300 ease-in-out
                      ${
                        isActive
                          ? "bg-white text-blue-700  font-semibold"
                          : "text-gray-600 hover:text-gray-800"
                      }
                    `}
                    onClick={() => {
                      console.log("Clicking tab:", tab.id);
                      setActiveTab(tab.id);
                    }}
                  >
                    <IconComponent className="w-4 h-4 transition-transform duration-300" />
                    <span>{tab.label}</span>

                    {/* Underline fixo para cada aba */}
                    <div
                      className={`
                        absolute bottom-0 left-1/2 transform -translate-x-1/2 h-0.5 bg-gradient-to-r from-blue-500 to-blue-600 transition-all duration-300
                        ${isActive ? "w-full opacity-100" : "w-0 opacity-0"}
                      `}
                    />
                  </button>
                );
              })}
            </div>

            {/* Avatar e Dropdown Menu */}
            <div className="relative" ref={dropdownRef}>
              <button
                onClick={() => setIsDropdownOpen(!isDropdownOpen)}
                className="flex items-center space-x-3 px-3 py-2 rounded-lg hover:bg-gray-100 transition-colors duration-200"
              >
                <div
                  className={`w-10 h-10 rounded-full bg-gradient-to-br ${avatarGradient} flex items-center justify-center text-white font-semibold shadow-lg ring-2 ring-white`}
                >
                  {userInitials}
                </div>
                <div className="hidden lg:block text-left">
                  <p className="text-sm font-semibold text-gray-900">
                    {userName}
                  </p>
                  <p className="text-xs text-gray-500">{userEmail}</p>
                </div>
                <ChevronDown
                  className={`w-4 h-4 text-gray-600 transition-transform duration-200 ${
                    isDropdownOpen ? "rotate-180" : ""
                  }`}
                />
              </button>

              {/* Dropdown Menu */}
              {isDropdownOpen && (
                <div className="absolute right-0 mt-2 w-64 bg-white rounded-lg shadow-xl border border-gray-200 py-2 z-50">
                  <div className="px-4 py-3 border-b border-gray-100">
                    <div className="flex items-center space-x-3">
                      <div
                        className={`w-12 h-12 rounded-full bg-gradient-to-br ${avatarGradient} flex items-center justify-center text-white font-bold shadow-lg`}
                      >
                        {userInitials}
                      </div>
                      <div className="flex-1 min-w-0">
                        <p className="text-sm font-semibold text-gray-900 truncate">
                          {userName}
                        </p>
                        <p className="text-xs text-gray-500 truncate">
                          {userEmail}
                        </p>
                      </div>
                    </div>
                  </div>

                  <button
                    onClick={handleProfile}
                    className="w-full flex items-center space-x-3 px-4 py-3 hover:bg-gray-50 transition-colors duration-200"
                  >
                    <User className="w-4 h-4 text-gray-600" />
                    <span className="text-sm text-gray-700">
                      Informações pessoais
                    </span>
                  </button>

                  <button
                    onClick={handleEditProfile}
                    className="w-full flex items-center space-x-3 px-4 py-3 hover:bg-gray-50 transition-colors duration-200"
                  >
                    <Settings className="w-4 h-4 text-gray-600" />
                    <span className="text-sm text-gray-700">Editar perfil</span>
                  </button>

                  <div className="border-t border-gray-100 mt-2 pt-2">
                    <button
                      onClick={handleLogoutClick}
                      className="w-full flex items-center space-x-3 px-4 py-3 hover:bg-red-50 transition-colors duration-200 text-red-600"
                    >
                      <LogOut className="w-4 h-4" />
                      <span className="text-sm font-medium">Sair</span>
                    </button>
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>

        {/* Navegação Mobile */}
        <div className="md:hidden">
          <div className="flex items-center justify-between mb-4">
            <div className="flex space-x-1 bg-gray-100 p-1 rounded-xl flex-1 mr-3">
              {tabs.map((tab) => {
                const IconComponent = tab.icon;
                const isActive = activeTab === tab.id;

                return (
                  <button
                    key={tab.id}
                    className={`
                      relative flex-1 flex flex-col items-center space-y-1 py-3 px-4 rounded-lg transition-all duration-300 ease-in-out
                      ${
                        isActive
                          ? "text-gray-900 bg-white border border-gray-300 font-semibold"
                          : "text-gray-600 bg-gray-50 border border-gray-200 hover:bg-gray-100"
                      }
                    `}
                    onClick={() => setActiveTab(tab.id)}
                  >
                    <IconComponent className="w-4 h-4 transition-transform duration-300" />
                    <span className="text-xs">{tab.label}</span>

                    {/* Underline para mobile */}
                    <div
                      className={`
                        absolute bottom-0 left-1/2 transform -translate-x-1/2 h-0.5 bg-gradient-to-r from-blue-500 to-blue-600 transition-all duration-300
                        ${isActive ? "w-full opacity-100" : "w-0 opacity-0"}
                      `}
                    />
                  </button>
                );
              })}
            </div>

            {/* Avatar Mobile */}
            <div className="relative" ref={dropdownRefMobile}>
              <button
                onClick={() => setIsDropdownOpen(!isDropdownOpen)}
                className={`w-10 h-10 rounded-full bg-gradient-to-br ${avatarGradient} flex items-center justify-center text-white font-semibold shadow-lg ring-2 ring-white`}
              >
                {userInitials}
              </button>

              {/* Dropdown Menu Mobile */}
              {isDropdownOpen && (
                <div className="absolute right-0 mt-2 w-64 bg-white rounded-lg shadow-xl border border-gray-200 py-2 z-50">
                  <div className="px-4 py-3 border-b border-gray-100">
                    <div className="flex items-center space-x-3">
                      <div
                        className={`w-12 h-12 rounded-full bg-gradient-to-br ${avatarGradient} flex items-center justify-center text-white font-bold shadow-lg`}
                      >
                        {userInitials}
                      </div>
                      <div className="flex-1 min-w-0">
                        <p className="text-sm font-semibold text-gray-900 truncate">
                          {userName}
                        </p>
                        <p className="text-xs text-gray-500 truncate">
                          {userEmail}
                        </p>
                      </div>
                    </div>
                  </div>

                  <button
                    onClick={handleProfile}
                    className="w-full flex items-center space-x-3 px-4 py-3 hover:bg-gray-50 transition-colors duration-200"
                  >
                    <User className="w-4 h-4 text-gray-600" />
                    <span className="text-sm text-gray-700">
                      Informações pessoais
                    </span>
                  </button>

                  <button
                    onClick={handleEditProfile}
                    className="w-full flex items-center space-x-3 px-4 py-3 hover:bg-gray-50 transition-colors duration-200"
                  >
                    <Settings className="w-4 h-4 text-gray-600" />
                    <span className="text-sm text-gray-700">Editar perfil</span>
                  </button>

                  <div className="border-t border-gray-100 mt-2 pt-2">
                    <button
                      onClick={handleLogoutClick}
                      className="w-full flex items-center space-x-3 px-4 py-3 hover:bg-red-50 transition-colors duration-200 text-red-600"
                    >
                      <LogOut className="w-4 h-4" />
                      <span className="text-sm font-medium">Sair</span>
                    </button>
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </nav>
  );
};

export default Navbar;

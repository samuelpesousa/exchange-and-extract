import React from "react";
import { TrendingUp, FileText } from "lucide-react";
import logoImage from "../image/acorianalight.png";

interface NavbarProps {
  activeTab: string;
  setActiveTab: (tab: string) => void;
}

interface Tab {
  id: string;
  label: string;
  icon: React.ComponentType<{ className?: string }>;
}

const Navbar: React.FC<NavbarProps> = ({ activeTab, setActiveTab }) => {
  const tabs: Tab[] = [
    { id: "exchange", label: "Cotação e Câmbio", icon: TrendingUp },
    { id: "extract", label: "Extrato", icon: FileText },
  ];

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
          <div className="hidden md:flex space-x-1 p-1 rounded-xl">
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
                    ${tab.id === "exchange" ? "rounded-tl rounded-bl" : ""}
                    ${tab.id === "extract" ? "rounded-tr rounded-br" : ""}
                  `}
                  onClick={() => setActiveTab(tab.id)}
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
        </div>

        {/* Navegação Mobile */}
        <div className="md:hidden">
          <div className="flex space-x-1 bg-gray-100 p-1 rounded-xl mb-4">
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
                  <span>{tab.label}</span>

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
        </div>
      </div>
    </nav>
  );
};

export default Navbar;

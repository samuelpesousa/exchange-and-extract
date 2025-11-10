import React, { useState } from "react";
import Navbar from "./components/Navbar";
import ExchangeRate from "./components/ExchangeRate";
import Extract from "./components/Extract";
import "./App.css";

const App: React.FC = () => {
  const [activeTab, setActiveTab] = useState<string>("exchange");

  const renderActiveTab = () => {
    switch (activeTab) {
      case "exchange":
        return <ExchangeRate />;
      case "extract":
        return <Extract />;
      default:
        return <ExchangeRate />;
    }
  };

  return (
    <div className="App">
      <Navbar activeTab={activeTab} setActiveTab={setActiveTab} />
      <main className="main-content">{renderActiveTab()}</main>
    </div>
  );
};

export default App;

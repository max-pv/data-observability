import React from "react";
import { BrowserRouter as Router, Route, Routes, Link, NavLink } from "react-router-dom";
import Dashboard from "./pages/Dashboard";
import HistoricalData from "./pages/HistoricalData";
import "./App.css";

const App: React.FC = () => {
  return (
    <Router>
      <div className="app">
        <nav>
          <NavLink to="/" end>
            Live Dashboard
          </NavLink>
          <NavLink to="/historical">
            Historical Data
          </NavLink>
        </nav>

        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/historical" element={<HistoricalData />} />
        </Routes>
      </div>
    </Router>
  );
};

export default App;

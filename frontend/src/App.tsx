import { useState } from "react";
import { AnimatePresence, motion } from "motion/react";
import Sidebar, { type TabKey } from "./components/Sidebar";
import Header from "./components/Header";
import DashboardPage from "./pages/DashboardPage";
import ProductsPage from "./pages/ProductsPage";
import OrdersPage from "./pages/OrdersPage";
import PaymentsPage from "./pages/PaymentsPage";
import UsersPage from "./pages/UsersPage";
import SettingsPage from "./pages/SettingsPage";
import { useTheme } from "./hooks/useTheme";

const tabTitles: Record<TabKey, string> = {
  dashboard: "Dashboard",
  products: "Produits",
  orders: "Commandes",
  payments: "Paiements",
  users: "Utilisateurs",
  settings: "Paramètres",
};

const searchPages: TabKey[] = ["products", "orders", "payments"];

export default function App() {
  useTheme();
  const [activeTab, setActiveTab] = useState<TabKey>("dashboard");
  const [search, setSearch] = useState("");

  // Reset search on tab change
  const handleTabChange = (tab: TabKey) => {
    setActiveTab(tab);
    setSearch("");
  };

  const renderPage = () => {
    const pageKey = activeTab;
    switch (pageKey) {
      case "dashboard":
        return <DashboardPage />;
      case "products":
        return <ProductsPage />;
      case "orders":
        return <OrdersPage />;
      case "payments":
        return <PaymentsPage />;
      case "users":
        return <UsersPage />;
      case "settings":
        return <SettingsPage />;
      default:
        return <DashboardPage />;
    }
  };

  return (
    <div className="flex h-screen bg-surface text-text-primary overflow-hidden">
      {/* Sidebar */}
      <Sidebar activeTab={activeTab} onTabChange={handleTabChange} />

      {/* Main area */}
      <div className="flex-1 flex flex-col ml-[var(--sidebar-width)] min-w-0">
        <Header
          title={tabTitles[activeTab]}
          search={searchPages.includes(activeTab) ? search : undefined}
          onSearch={searchPages.includes(activeTab) ? setSearch : undefined}
        />

        <main className="flex-1 overflow-y-auto">
          <AnimatePresence mode="wait">
            <motion.div
              key={activeTab}
              initial={{ opacity: 0, y: 8 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -8 }}
              transition={{ duration: 0.15, ease: "easeOut" }}
            >
              {renderPage()}
            </motion.div>
          </AnimatePresence>
        </main>
      </div>
    </div>
  );
}

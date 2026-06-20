import {
  ChartPieSlice,
  Package,
  ShoppingCart,
  CreditCard,
  Users,
  Gear,
  ArrowSquareOut,
} from "@phosphor-icons/react";

export type TabKey = "dashboard" | "products" | "orders" | "payments" | "users" | "settings";

interface NavItem {
  key: TabKey;
  label: string;
  icon: React.ReactNode;
}

const navItems: NavItem[] = [
  { key: "dashboard", label: "Dashboard",      icon: <ChartPieSlice size={20} weight="duotone" /> },
  { key: "products",  label: "Produits",       icon: <Package size={20} weight="duotone" /> },
  { key: "orders",    label: "Commandes",      icon: <ShoppingCart size={20} weight="duotone" /> },
  { key: "payments",  label: "Paiements",      icon: <CreditCard size={20} weight="duotone" /> },
  { key: "users",     label: "Utilisateurs",   icon: <Users size={20} weight="duotone" /> },
  { key: "settings",  label: "Paramètres",     icon: <Gear size={20} weight="duotone" /> },
];

interface SidebarProps {
  activeTab: TabKey;
  onTabChange: (tab: TabKey) => void;
}

export default function Sidebar({ activeTab, onTabChange }: SidebarProps) {
  return (
    <aside className="fixed top-0 left-0 z-40 h-screen w-[var(--sidebar-width)] bg-surface-sidebar border-r border-border flex flex-col">
      {/* Logo */}
      <div className="flex items-center gap-3 px-5 h-16 border-b border-border shrink-0">
        <div className="w-8 h-8 rounded-lg bg-accent/20 flex items-center justify-center">
          <svg width="18" height="18" viewBox="0 0 32 32" fill="none">
            <path d="M8 16L14 10L18 14L24 8" stroke="#06B6D4" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round" />
            <path d="M8 22L14 16L18 20L24 14" stroke="#4f46e5" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round" opacity="0.6" />
          </svg>
        </div>
        <div>
          <span className="text-sm font-semibold text-text-primary tracking-tight">NexusFlow</span>
          <p className="text-[10px] text-text-muted leading-tight">Dashboard</p>
        </div>
      </div>

      {/* Navigation */}
      <nav className="flex-1 px-3 py-4 space-y-1 overflow-y-auto">
        {navItems.map((item) => {
          const isActive = activeTab === item.key;
          return (
            <button
              key={item.key}
              onClick={() => onTabChange(item.key)}
              className={`w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-all duration-150 ${
                isActive
                  ? "bg-accent/10 text-accent shadow-[inset_2px_0_0_var(--color-accent)]"
                  : "text-text-secondary hover:text-text-primary hover:bg-surface-hover"
              }`}
            >
              <span className={`shrink-0 ${isActive ? "text-accent" : "text-text-muted"}`}>
                {item.icon}
              </span>
              {item.label}
            </button>
          );
        })}
      </nav>

      {/* Footer */}
      <div className="px-3 py-4 border-t border-border shrink-0">
        <a
          href="#"
          className="flex items-center gap-3 px-3 py-2 rounded-lg text-xs text-text-muted hover:text-text-secondary hover:bg-surface-hover transition-colors"
        >
          <ArrowSquareOut size={16} />
          Documentation
        </a>
      </div>
    </aside>
  );
}

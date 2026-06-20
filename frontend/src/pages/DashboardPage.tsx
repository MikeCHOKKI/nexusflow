import { useState, useEffect } from "react";
import type { ReactNode } from "react";
import { motion } from "motion/react";
import {
  Package,
  ShoppingCart,
  CurrencyDollar,
  Users,
  Circle,
  WifiHigh,
  WarningCircle,
} from "@phosphor-icons/react";
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts";
import StatsCard from "../components/ui/StatsCard";
import StatusBadge from "../components/ui/StatusBadge";
import type { DashboardStats, ServiceHealth } from "../lib/types";
import { dashboard } from "../lib/api";

function ServiceDot({ status }: { status: ServiceHealth["status"] }) {
  const colors = {
    healthy: "text-success",
    degraded: "text-warning",
    down: "text-danger",
  };
  return <Circle size={10} weight="fill" className={`${colors[status]}`} />;
}

const serviceIcons: Record<string, ReactNode> = {
  "api-gateway": <WifiHigh size={18} />,
  "service-catalog": <Package size={18} />,
  "service-order": <ShoppingCart size={18} />,
  "service-payment": <CurrencyDollar size={18} />,
  "service-user": <Users size={18} />,
  "service-notify": <WarningCircle size={18} />,
};

const defaultServices: ServiceHealth[] = [
  { name: "api-gateway",    status: "healthy",  uptime: "99.9%", latency: 12 },
  { name: "service-catalog", status: "healthy", uptime: "99.8%", latency: 8 },
  { name: "service-order",  status: "healthy",  uptime: "99.7%", latency: 15 },
  { name: "service-payment",status: "healthy",  uptime: "99.9%", latency: 10 },
  { name: "service-user",   status: "healthy",  uptime: "99.9%", latency: 6 },
  { name: "service-notify", status: "degraded", uptime: "97.2%", latency: 45 },
];

const defaultStats: DashboardStats = {
  totalProducts: 1_284,
  totalOrders: 3_571,
  totalRevenue: 48_290_000,
  totalUsers: 892,
  ordersTrend: [
    { date: "Lun", count: 120, revenue: 1_800_000 },
    { date: "Mar", count: 98,  revenue: 1_450_000 },
    { date: "Mer", count: 145, revenue: 2_100_000 },
    { date: "Jeu", count: 132, revenue: 1_950_000 },
    { date: "Ven", count: 168, revenue: 2_450_000 },
    { date: "Sam", count: 110, revenue: 1_600_000 },
    { date: "Dim", count: 85,  revenue: 1_200_000 },
  ],
  recentOrders: [
    { id: "CMD-001", userId: "1", userName: "Fatou Diallo", items: [], total: 45_000, status: "delivered", paymentMethod: "wave", shippingAddress: "", createdAt: "2026-06-19T14:30:00Z", updatedAt: "2026-06-19T14:30:00Z" },
    { id: "CMD-002", userId: "2", userName: "Mamadou Ndiaye", items: [], total: 128_500, status: "shipped", paymentMethod: "orange_money", shippingAddress: "", createdAt: "2026-06-19T12:15:00Z", updatedAt: "2026-06-19T12:15:00Z" },
    { id: "CMD-003", userId: "3", userName: "Aïcha Ba", items: [], total: 23_900, status: "processing", paymentMethod: "visa", shippingAddress: "", createdAt: "2026-06-19T10:00:00Z", updatedAt: "2026-06-19T10:00:00Z" },
    { id: "CMD-004", userId: "4", userName: "Oumar Fall", items: [], total: 312_000, status: "pending", paymentMethod: "mtn_money", shippingAddress: "", createdAt: "2026-06-19T09:45:00Z", updatedAt: "2026-06-19T09:45:00Z" },
    { id: "CMD-005", userId: "5", userName: "Khadija Sow", items: [], total: 67_800, status: "delivered", paymentMethod: "wave", shippingAddress: "", createdAt: "2026-06-18T16:20:00Z", updatedAt: "2026-06-18T16:20:00Z" },
  ],
  services: defaultServices,
};

function formatCFA(amount: number): string {
  return new Intl.NumberFormat("fr-FR").format(amount) + " FCFA";
}

export default function DashboardPage() {
  const [stats, setStats] = useState<DashboardStats>(defaultStats);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    dashboard
      .stats()
      .then((res) => {
        if (res.data) setStats(res.data);
      })
      .catch(() => {
        // Use default data on error
      })
      .finally(() => setLoading(false));
  }, []);

  if (loading) {
    return (
      <div className="p-8 space-y-6">
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          {Array.from({ length: 4 }).map((_, i) => (
            <div key={i} className="bg-surface-elevated border border-border rounded-xl p-5">
              <div className="skeleton h-10 w-10 rounded-lg mb-3" />
              <div className="skeleton h-7 w-24 mb-1" />
              <div className="skeleton h-4 w-20" />
            </div>
          ))}
        </div>
        <div className="skeleton h-80 rounded-xl" />
      </div>
    );
  }

  return (
    <div className="p-8 space-y-6">
      {/* Stats Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <StatsCard
          icon={<Package size={22} />}
          label="Total Produits"
          value={stats.totalProducts.toLocaleString("fr-FR")}
          trend={12}
          index={0}
        />
        <StatsCard
          icon={<ShoppingCart size={22} />}
          label="Total Commandes"
          value={stats.totalOrders.toLocaleString("fr-FR")}
          trend={8}
          index={1}
        />
        <StatsCard
          icon={<CurrencyDollar size={22} />}
          label="Revenu Total"
          value={formatCFA(stats.totalRevenue)}
          trend={15}
          index={2}
        />
        <StatsCard
          icon={<Users size={22} />}
          label="Utilisateurs"
          value={stats.totalUsers.toLocaleString("fr-FR")}
          trend={-3}
          index={3}
        />
      </div>

      {/* Chart + Services */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Orders Trend Chart */}
        <motion.div
          initial={{ opacity: 0, y: 16 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.35, delay: 0.25 }}
          className="lg:col-span-2 bg-surface-elevated border border-border rounded-xl p-6"
        >
          <h3 className="text-sm font-semibold text-text-primary mb-4">
            Évolution des commandes
          </h3>
          <div className="h-64">
            <ResponsiveContainer width="100%" height="100%">
              <LineChart data={stats.ordersTrend}>
                <CartesianGrid strokeDasharray="3 3" stroke="rgba(255,255,255,0.04)" />
                <XAxis
                  dataKey="date"
                  tick={{ fill: "#94a3b8", fontSize: 12 }}
                  axisLine={{ stroke: "rgba(255,255,255,0.06)" }}
                  tickLine={false}
                />
                <YAxis
                  tick={{ fill: "#94a3b8", fontSize: 12 }}
                  axisLine={false}
                  tickLine={false}
                />
                <Tooltip
                  contentStyle={{
                    backgroundColor: "#1a1a3e",
                    border: "1px solid rgba(255,255,255,0.1)",
                    borderRadius: 8,
                    color: "#e2e8f0",
                    fontSize: 12,
                  }}
                />
                <Line
                  type="monotone"
                  dataKey="count"
                  stroke="#06B6D4"
                  strokeWidth={2}
                  dot={{ fill: "#06B6D4", r: 3 }}
                  activeDot={{ r: 5, fill: "#06B6D4" }}
                />
              </LineChart>
            </ResponsiveContainer>
          </div>
        </motion.div>

        {/* Microservices Status */}
        <motion.div
          initial={{ opacity: 0, y: 16 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.35, delay: 0.3 }}
          className="bg-surface-elevated border border-border rounded-xl p-6"
        >
          <h3 className="text-sm font-semibold text-text-primary mb-4">
            Statut des services
          </h3>
          <div className="space-y-3">
            {(stats.services ?? defaultServices).map((svc) => (
              <div
                key={svc.name}
                className="flex items-center justify-between py-2 px-3 rounded-lg bg-surface/50"
              >
                <div className="flex items-center gap-3">
                  <span className="text-text-muted">
                    {serviceIcons[svc.name] ?? <WifiHigh size={18} />}
                  </span>
                  <div>
                    <p className="text-xs font-medium text-text-primary capitalize">
                      {svc.name.replace(/-/g, " ")}
                    </p>
                    <p className="text-[10px] text-text-muted">
                      {svc.latency}ms · {svc.uptime}
                    </p>
                  </div>
                </div>
                <ServiceDot status={svc.status} />
              </div>
            ))}
          </div>
        </motion.div>
      </div>

      {/* Recent Orders */}
      <motion.div
        initial={{ opacity: 0, y: 16 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.35, delay: 0.35 }}
        className="bg-surface-elevated border border-border rounded-xl overflow-hidden"
      >
        <div className="px-6 py-4 border-b border-border">
          <h3 className="text-sm font-semibold text-text-primary">
            Dernières commandes
          </h3>
        </div>
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b border-border">
                <th className="px-6 py-3 text-left text-xs font-medium text-text-muted uppercase">ID</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-text-muted uppercase">Client</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-text-muted uppercase">Total</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-text-muted uppercase">Statut</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-text-muted uppercase">Date</th>
              </tr>
            </thead>
            <tbody>
              {stats.recentOrders.map((order) => (
                <tr key={order.id} className="table-row border-b border-border/30">
                  <td className="px-6 py-3 text-sm font-mono text-accent">{order.id}</td>
                  <td className="px-6 py-3 text-sm text-text-primary">{order.userName}</td>
                  <td className="px-6 py-3 text-sm text-text-primary font-medium">
                    {formatCFA(order.total)}
                  </td>
                  <td className="px-6 py-3">
                    <StatusBadge status={order.status} />
                  </td>
                  <td className="px-6 py-3 text-sm text-text-muted">
                    {new Date(order.createdAt).toLocaleDateString("fr-FR", {
                      day: "numeric",
                      month: "short",
                    })}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </motion.div>
    </div>
  );
}

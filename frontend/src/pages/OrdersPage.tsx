import { useState, useCallback } from "react";
import { Eye } from "@phosphor-icons/react";
import DataTable, { type Column } from "../components/ui/DataTable";
import StatusBadge from "../components/ui/StatusBadge";
import Modal from "../components/ui/Modal";
import PageContainer from "../components/PageContainer";
import type { Order, OrderStatus } from "../lib/types";
import { orders as ordersApi } from "../lib/api";

const statusFilters: (OrderStatus | "")[] = [
  "", "pending", "confirmed", "processing", "shipped", "delivered", "cancelled",
];

const defaultOrders: Order[] = [
  { id: "CMD-001", userId: "1", userName: "Fatou Diallo", items: [{ productId: "P1", productName: "Smartphone X200", quantity: 1, unitPrice: 450000 }], total: 450000, status: "delivered", paymentMethod: "wave", shippingAddress: "Dakar, Sénégal", createdAt: "2026-06-19T14:30:00Z", updatedAt: "2026-06-19T14:30:00Z" },
  { id: "CMD-002", userId: "2", userName: "Mamadou Ndiaye", items: [{ productId: "P2", productName: "Casque Bluetooth Pro", quantity: 2, unitPrice: 85000 }], total: 170000, status: "shipped", paymentMethod: "orange_money", shippingAddress: "Thiès, Sénégal", createdAt: "2026-06-19T12:15:00Z", updatedAt: "2026-06-19T12:15:00Z" },
  { id: "CMD-003", userId: "3", userName: "Aïcha Ba", items: [{ productId: "P4", productName: "Robot Aspirateur", quantity: 1, unitPrice: 295000 }], total: 295000, status: "processing", paymentMethod: "visa", shippingAddress: "Saint-Louis, Sénégal", createdAt: "2026-06-19T10:00:00Z", updatedAt: "2026-06-19T10:00:00Z" },
  { id: "CMD-004", userId: "4", userName: "Oumar Fall", items: [{ productId: "P5", productName: "Huile d'Olive Bio 1L", quantity: 10, unitPrice: 8500 }], total: 85000, status: "pending", paymentMethod: "mtn_money", shippingAddress: "Dakar, Sénégal", createdAt: "2026-06-19T09:45:00Z", updatedAt: "2026-06-19T09:45:00Z" },
  { id: "CMD-005", userId: "5", userName: "Khadija Sow", items: [{ productId: "P6", productName: "Montre Connectée S3", quantity: 1, unitPrice: 175000 }], total: 175000, status: "delivered", paymentMethod: "wave", shippingAddress: "Touba, Sénégal", createdAt: "2026-06-18T16:20:00Z", updatedAt: "2026-06-18T16:20:00Z" },
  { id: "CMD-006", userId: "1", userName: "Fatou Diallo", items: [{ productId: "P7", productName: "Sac à Main Cuir", quantity: 1, unitPrice: 95000 }], total: 95000, status: "cancelled", paymentMethod: "visa", shippingAddress: "Dakar, Sénégal", createdAt: "2026-06-18T14:00:00Z", updatedAt: "2026-06-18T14:00:00Z" },
  { id: "CMD-007", userId: "6", userName: "Ibrahima Gueye", items: [{ productId: "P10", productName: "Ballon de Foot Pro", quantity: 3, unitPrice: 32000 }], total: 96000, status: "pending", paymentMethod: "orange_money", shippingAddress: "Ziguinchor, Sénégal", createdAt: "2026-06-18T11:30:00Z", updatedAt: "2026-06-18T11:30:00Z" },
];

const statusLabels: Record<OrderStatus, string> = {
  pending: "pending",
  confirmed: "confirmed",
  processing: "processing",
  shipped: "shipped",
  delivered: "delivered",
  cancelled: "cancelled",
};

function formatCFA(amount: number): string {
  return new Intl.NumberFormat("fr-FR").format(amount) + " FCFA";
}

export default function OrdersPage() {
  const [orderList, setOrderList] = useState<Order[]>(defaultOrders);
  const [page, setPage] = useState(1);
  const [totalPages] = useState(2);
  const [statusFilter, setStatusFilter] = useState<string>("");
  const [selectedOrder, setSelectedOrder] = useState<Order | null>(null);

  const handleStatusUpdate = useCallback(
    async (id: string, newStatus: OrderStatus) => {
      try {
        await ordersApi.updateStatus(id, newStatus);
      } catch {
        // fallback
      }
      setOrderList((prev) =>
        prev.map((o) => (o.id === id ? { ...o, status: newStatus } : o)),
      );
      if (selectedOrder?.id === id) {
        setSelectedOrder((prev) => prev ? { ...prev, status: newStatus } : null);
      }
    },
    [selectedOrder],
  );

  const filtered = statusFilter
    ? orderList.filter((o) => o.status === statusFilter)
    : orderList;

  const columns: Column<Order>[] = [
    { key: "id", header: "ID", className: "font-mono text-accent text-xs", sortable: true },
    { key: "userName", header: "Client", sortable: true },
    {
      key: "total",
      header: "Total",
      render: (o) => (
        <span className="font-medium">{formatCFA(o.total)}</span>
      ),
      sortable: true,
    },
    {
      key: "paymentMethod",
      header: "Paiement",
      render: (o) => <StatusBadge status={o.paymentMethod} />,
    },
    {
      key: "status",
      header: "Statut",
      render: (o) => <StatusBadge status={o.status} />,
    },
    {
      key: "createdAt",
      header: "Date",
      render: (o) => (
        <span className="text-text-muted">
          {new Date(o.createdAt).toLocaleDateString("fr-FR")}
        </span>
      ),
    },
    {
      key: "actions",
      header: "",
      render: (o) => (
        <button
          onClick={() => setSelectedOrder(o)}
          className="p-1.5 rounded-md text-text-muted hover:text-accent hover:bg-accent/10 transition-colors"
        >
          <Eye size={16} />
        </button>
      ),
    },
  ];

  return (
    <PageContainer title="Commandes">
      {/* Filters */}
      <div className="flex items-center gap-2 mb-4 flex-wrap">
        {statusFilters.map((s) => (
          <button
            key={s}
            onClick={() => setStatusFilter(s)}
            className={`px-3 py-1.5 text-xs font-medium rounded-lg transition-colors ${
              statusFilter === s
                ? "bg-accent/10 text-accent border border-accent/30"
                : "bg-surface border border-border text-text-secondary hover:text-text-primary"
            }`}
          >
            {s === "" ? "Tous" : statusLabels[s as OrderStatus] ?? s}
          </button>
        ))}
      </div>

      {/* Table */}
      <div className="bg-surface-elevated border border-border rounded-xl overflow-hidden">
        <DataTable
          columns={columns}
          data={filtered}
          keyExtractor={(o) => o.id}
          page={page}
          totalPages={totalPages}
          onPageChange={setPage}
        />
      </div>

      {/* Detail Modal */}
      <Modal
        open={selectedOrder !== null}
        onClose={() => setSelectedOrder(null)}
        title={`Commande ${selectedOrder?.id ?? ""}`}
        size="lg"
      >
        {selectedOrder && (
          <div className="space-y-5">
            {/* Info */}
            <div className="grid grid-cols-2 gap-4">
              <div>
                <p className="text-xs text-text-muted">Client</p>
                <p className="text-sm font-medium text-text-primary">{selectedOrder.userName}</p>
              </div>
              <div>
                <p className="text-xs text-text-muted">Date</p>
                <p className="text-sm font-medium text-text-primary">
                  {new Date(selectedOrder.createdAt).toLocaleDateString("fr-FR", {
                    day: "numeric",
                    month: "long",
                    year: "numeric",
                    hour: "2-digit",
                    minute: "2-digit",
                  })}
                </p>
              </div>
              <div>
                <p className="text-xs text-text-muted">Statut</p>
                <StatusBadge status={selectedOrder.status} />
              </div>
              <div>
                <p className="text-xs text-text-muted">Paiement</p>
                <StatusBadge status={selectedOrder.paymentMethod} />
              </div>
              <div className="col-span-2">
                <p className="text-xs text-text-muted">Adresse de livraison</p>
                <p className="text-sm text-text-primary">{selectedOrder.shippingAddress}</p>
              </div>
            </div>

            {/* Items */}
            <div>
              <p className="text-xs font-medium text-text-muted uppercase tracking-wider mb-2">
                Articles
              </p>
              <div className="space-y-2">
                {selectedOrder.items.map((item, i) => (
                  <div
                    key={i}
                    className="flex items-center justify-between px-3 py-2 rounded-lg bg-surface/50"
                  >
                    <div>
                      <p className="text-sm text-text-primary">{item.productName}</p>
                      <p className="text-xs text-text-muted">
                        Qté: {item.quantity} × {formatCFA(item.unitPrice)}
                      </p>
                    </div>
                    <p className="text-sm font-medium text-text-primary">
                      {formatCFA(item.quantity * item.unitPrice)}
                    </p>
                  </div>
                ))}
              </div>
            </div>

            {/* Total */}
            <div className="flex justify-between items-center pt-3 border-t border-border">
              <span className="text-sm font-semibold text-text-primary">Total</span>
              <span className="text-lg font-bold text-text-primary">
                {formatCFA(selectedOrder.total)}
              </span>
            </div>

            {/* Status update */}
            <div className="flex items-center gap-2 pt-2">
              <span className="text-xs text-text-muted">Mettre à jour :</span>
              {(["pending", "confirmed", "processing", "shipped", "delivered", "cancelled"] as OrderStatus[]).map(
                (s) => (
                  <button
                    key={s}
                    onClick={() => handleStatusUpdate(selectedOrder.id, s)}
                    className={`px-2.5 py-1 text-xs rounded-lg transition-colors ${
                      selectedOrder.status === s
                        ? "bg-accent/10 text-accent"
                        : "bg-surface border border-border text-text-secondary hover:text-text-primary"
                    }`}
                  >
                    {statusLabels[s]}
                  </button>
                ),
              )}
            </div>
          </div>
        )}
      </Modal>
    </PageContainer>
  );
}

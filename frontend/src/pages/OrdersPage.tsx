import { useState, useEffect, useCallback } from "react";
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
  const [orderList, setOrderList] = useState<Order[]>([]);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [loading, setLoading] = useState(true);
  const [statusFilter, setStatusFilter] = useState<string>("");
  const [selectedOrder, setSelectedOrder] = useState<Order | null>(null);

  useEffect(() => {
    setLoading(true);
    ordersApi
      .list({ page, status: statusFilter ? (statusFilter as OrderStatus) : undefined })
      .then((res) => {
        setOrderList(res.data);
        setTotalPages(res.totalPages);
      })
      .catch(() => {})
      .finally(() => setLoading(false));
  }, [page, statusFilter]);

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

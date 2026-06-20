import { useState, useEffect, useCallback } from "react";
import { ArrowsClockwise } from "@phosphor-icons/react";
import DataTable, { type Column } from "../components/ui/DataTable";
import StatusBadge from "../components/ui/StatusBadge";
import PageContainer from "../components/PageContainer";
import type { Payment, PaymentStatus } from "../lib/types";
import { payments as paymentsApi } from "../lib/api";

const statusFilters: (PaymentStatus | "")[] = ["", "completed", "pending", "failed", "refunded"];

function formatCFA(amount: number): string {
  return new Intl.NumberFormat("fr-FR").format(amount) + " FCFA";
}

export default function PaymentsPage() {
  const [paymentList, setPaymentList] = useState<Payment[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [statusFilter, setStatusFilter] = useState<string>("");
  const [selectedPayment, setSelectedPayment] = useState<Payment | null>(null);
  const [refunding, setRefunding] = useState(false);

  useEffect(() => {
    setLoading(true);
    paymentsApi
      .list({ page, status: statusFilter ? (statusFilter as PaymentStatus) : undefined })
      .then((res) => {
        setPaymentList(res.data);
        setTotalPages(res.totalPages);
      })
      .catch(() => {})
      .finally(() => setLoading(false));
  }, [page, statusFilter]);

  const handleRefund = useCallback(async (id: string) => {
    setRefunding(true);
    try {
      await paymentsApi.refund(id);
    } catch {
      // fallback
    }
    setPaymentList((prev) =>
      prev.map((p) =>
        p.id === id ? { ...p, status: "refunded" as PaymentStatus } : p,
      ),
    );
    if (selectedPayment?.id === id) {
      setSelectedPayment((prev) => prev ? { ...prev, status: "refunded" } : null);
    }
    setRefunding(false);
  }, [selectedPayment]);

  const filtered = statusFilter
    ? paymentList.filter((p) => p.status === statusFilter)
    : paymentList;

  const columns: Column<Payment>[] = [
    { key: "id", header: "ID", className: "font-mono text-accent text-xs", sortable: true },
    { key: "orderId", header: "Commande", className: "font-mono text-xs text-text-muted" },
    {
      key: "amount",
      header: "Montant",
      render: (p) => <span className="font-medium">{formatCFA(p.amount)}</span>,
      sortable: true,
    },
    {
      key: "method",
      header: "Méthode",
      render: (p) => <StatusBadge status={p.method} />,
    },
    {
      key: "status",
      header: "Statut",
      render: (p) => <StatusBadge status={p.status} />,
    },
    {
      key: "transactionRef",
      header: "Réf. Transaction",
      className: "font-mono text-xs text-text-muted",
    },
    {
      key: "actions",
      header: "",
      render: (p) =>
        p.status === "completed" && (
          <button
            onClick={() => handleRefund(p.id)}
            disabled={refunding}
            className="flex items-center gap-1 px-2 py-1 text-xs text-warning hover:bg-warning/10 rounded-md transition-colors disabled:opacity-50"
          >
            <ArrowsClockwise size={14} />
            Rembourser
          </button>
        ),
    },
  ];

  return (
    <PageContainer title="Paiements">
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
            {s === "" ? "Tous" : s}
          </button>
        ))}
      </div>

      {/* Table */}
      <div className="bg-surface-elevated border border-border rounded-xl overflow-hidden">
        <DataTable
          columns={columns}
          data={filtered}
          keyExtractor={(p) => p.id}
          page={page}
          totalPages={totalPages}
          onPageChange={setPage}
        />
      </div>
    </PageContainer>
  );
}

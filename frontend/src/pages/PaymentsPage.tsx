import { useState, useCallback } from "react";
import { ArrowsClockwise } from "@phosphor-icons/react";
import DataTable, { type Column } from "../components/ui/DataTable";
import StatusBadge from "../components/ui/StatusBadge";
import PageContainer from "../components/PageContainer";
import type { Payment, PaymentStatus } from "../lib/types";
import { payments as paymentsApi } from "../lib/api";

const statusFilters: (PaymentStatus | "")[] = ["", "completed", "pending", "failed", "refunded"];

const defaultPayments: Payment[] = [
  { id: "PAY-001", orderId: "CMD-001", userId: "1", amount: 450000, method: "wave", status: "completed", transactionRef: "WAV-6A3F2B", createdAt: "2026-06-19T14:31:00Z", updatedAt: "2026-06-19T14:31:00Z" },
  { id: "PAY-002", orderId: "CMD-002", userId: "2", amount: 170000, method: "orange_money", status: "completed", transactionRef: "OM-8D1E4C", createdAt: "2026-06-19T12:16:00Z", updatedAt: "2026-06-19T12:16:00Z" },
  { id: "PAY-003", orderId: "CMD-003", userId: "3", amount: 295000, method: "visa", status: "completed", transactionRef: "VISA-4F7A2E", createdAt: "2026-06-19T10:01:00Z", updatedAt: "2026-06-19T10:01:00Z" },
  { id: "PAY-004", orderId: "CMD-004", userId: "4", amount: 85000, method: "mtn_money", status: "pending", transactionRef: "MTN-2B5C8D", createdAt: "2026-06-19T09:46:00Z", updatedAt: "2026-06-19T09:46:00Z" },
  { id: "PAY-005", orderId: "CMD-005", userId: "5", amount: 175000, method: "wave", status: "completed", transactionRef: "WAV-1E3F7A", createdAt: "2026-06-18T16:21:00Z", updatedAt: "2026-06-18T16:21:00Z" },
  { id: "PAY-006", orderId: "CMD-006", userId: "1", amount: 95000, method: "visa", status: "refunded", transactionRef: "VISA-9C4D2B", createdAt: "2026-06-18T14:01:00Z", updatedAt: "2026-06-19T08:00:00Z" },
  { id: "PAY-007", orderId: "CMD-007", userId: "6", amount: 96000, method: "orange_money", status: "failed", transactionRef: "OM-5E8A1F", createdAt: "2026-06-18T11:31:00Z", updatedAt: "2026-06-18T11:31:00Z" },
  { id: "PAY-008", orderId: "CMD-002", userId: "2", amount: 45000, method: "wave", status: "completed", transactionRef: "WAV-7D2E4B", createdAt: "2026-06-18T09:15:00Z", updatedAt: "2026-06-18T09:15:00Z" },
];

function formatCFA(amount: number): string {
  return new Intl.NumberFormat("fr-FR").format(amount) + " FCFA";
}

export default function PaymentsPage() {
  const [paymentList, setPaymentList] = useState<Payment[]>(defaultPayments);
  const [page, setPage] = useState(1);
  const [totalPages] = useState(2);
  const [statusFilter, setStatusFilter] = useState<string>("");
  const [selectedPayment, setSelectedPayment] = useState<Payment | null>(null);
  const [refunding, setRefunding] = useState(false);

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

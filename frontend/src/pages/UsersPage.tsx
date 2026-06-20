import { useState, useEffect, useCallback } from "react";
import DataTable, { type Column } from "../components/ui/DataTable";
import StatusBadge from "../components/ui/StatusBadge";
import PageContainer from "../components/PageContainer";
import type { User } from "../lib/types";
import { users as usersApi } from "../lib/api";

export default function UsersPage() {
  const [userList, setUserList] = useState<User[]>([]);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);

  useEffect(() => {
    usersApi
      .list({ page })
      .then((res) => {
        setUserList(res.data);
        setTotalPages(res.totalPages);
      })
      .catch(() => {});
  }, [page]);

  const toggleStatus = useCallback(async (id: string) => {
    try {
      await usersApi.toggleStatus(id);
    } catch {
      // fallback
    }
    setUserList((prev) =>
      prev.map((u) =>
        u.id === id
          ? { ...u, status: u.status === "active" ? "inactive" : "active" }
          : u,
      ),
    );
  }, []);

  const columns: Column<User>[] = [
    { key: "id", header: "ID", className: "font-mono text-accent text-xs", sortable: true },
    { key: "name", header: "Nom", sortable: true },
    { key: "email", header: "Email" },
    {
      key: "role",
      header: "Rôle",
      render: (u) => <StatusBadge status={u.role} />,
    },
    {
      key: "status",
      header: "Statut",
      render: (u) => <StatusBadge status={u.status} />,
    },
    {
      key: "createdAt",
      header: "Inscription",
      render: (u) => (
        <span className="text-text-muted">
          {new Date(u.createdAt).toLocaleDateString("fr-FR")}
        </span>
      ),
    },
    {
      key: "actions",
      header: "",
      render: (u) => (
        <button
          onClick={() => toggleStatus(u.id)}
          className={`px-2.5 py-1 text-xs rounded-lg transition-colors ${
            u.status === "active"
              ? "text-danger bg-danger/10 hover:bg-danger/20"
              : "text-success bg-success/10 hover:bg-success/20"
          }`}
        >
          {u.status === "active" ? "Désactiver" : "Activer"}
        </button>
      ),
    },
  ];

  return (
    <PageContainer title="Utilisateurs">
      <div className="bg-surface-elevated border border-border rounded-xl overflow-hidden">
        <DataTable
          columns={columns}
          data={userList}
          keyExtractor={(u) => u.id}
          page={page}
          totalPages={totalPages}
          onPageChange={setPage}
        />
      </div>
    </PageContainer>
  );
}

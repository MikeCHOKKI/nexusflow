import { useState, useCallback } from "react";
import DataTable, { type Column } from "../components/ui/DataTable";
import StatusBadge from "../components/ui/StatusBadge";
import PageContainer from "../components/PageContainer";
import type { User } from "../lib/types";
import { users as usersApi } from "../lib/api";

const defaultUsers: User[] = [
  { id: "USR-001", email: "admin@nexusflow.io", name: "Admin", role: "admin", status: "active", createdAt: "2026-01-01T00:00:00Z", updatedAt: "2026-06-19T00:00:00Z" },
  { id: "USR-002", email: "fatou.diallo@email.com", name: "Fatou Diallo", role: "manager", status: "active", createdAt: "2026-02-15T00:00:00Z", updatedAt: "2026-06-18T00:00:00Z" },
  { id: "USR-003", email: "mamadou.ndiaye@email.com", name: "Mamadou Ndiaye", role: "viewer", status: "active", createdAt: "2026-03-01T00:00:00Z", updatedAt: "2026-06-17T00:00:00Z" },
  { id: "USR-004", email: "aicha.ba@email.com", name: "Aïcha Ba", role: "viewer", status: "active", createdAt: "2026-03-10T00:00:00Z", updatedAt: "2026-06-15T00:00:00Z" },
  { id: "USR-005", email: "oumar.fall@email.com", name: "Oumar Fall", role: "manager", status: "active", createdAt: "2026-03-20T00:00:00Z", updatedAt: "2026-06-14T00:00:00Z" },
  { id: "USR-006", email: "khadija.sow@email.com", name: "Khadija Sow", role: "viewer", status: "inactive", createdAt: "2026-04-01T00:00:00Z", updatedAt: "2026-06-10T00:00:00Z" },
  { id: "USR-007", email: "ibrahima.gueye@email.com", name: "Ibrahima Gueye", role: "viewer", status: "active", createdAt: "2026-04-15T00:00:00Z", updatedAt: "2026-06-08T00:00:00Z" },
  { id: "USR-008", email: "ndeye.sarr@email.com", name: "Ndeye Sarr", role: "manager", status: "active", createdAt: "2026-05-01T00:00:00Z", updatedAt: "2026-06-05T00:00:00Z" },
];

export default function UsersPage() {
  const [userList, setUserList] = useState<User[]>(defaultUsers);
  const [page, setPage] = useState(1);
  const [totalPages] = useState(1);

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

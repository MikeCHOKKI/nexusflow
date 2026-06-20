import { useState, useEffect, useCallback } from "react";
import { Plus, PencilSimple, TrashSimple } from "@phosphor-icons/react";
import DataTable, { type Column } from "../components/ui/DataTable";
import StatusBadge from "../components/ui/StatusBadge";
import Modal from "../components/ui/Modal";
import PageContainer from "../components/PageContainer";
import SearchInput from "../components/ui/SearchInput";
import type { Product, ProductInput, ProductStatus } from "../lib/types";
import { products as productsApi } from "../lib/api";

const categories = ["Électronique", "Mode", "Maison", "Alimentation", "Sport", "Beauté"];

const columns: Column<Product>[] = [
  { key: "id", header: "ID", className: "font-mono text-accent text-xs", sortable: true },
  { key: "name", header: "Nom", sortable: true },
  {
    key: "price",
    header: "Prix",
    render: (p) => `${p.price.toLocaleString("fr-FR")} FCFA`,
    sortable: true,
  },
  {
    key: "stock",
    header: "Stock",
    render: (p) => (
      <span className={p.stock <= 5 ? "text-danger font-medium" : ""}>
        {p.stock}
      </span>
    ),
    sortable: true,
  },
  {
    key: "category",
    header: "Catégorie",
    render: (p) => (
      <span className="text-xs bg-accent/10 text-accent px-2 py-0.5 rounded-full">
        {p.category}
      </span>
    ),
  },
  {
    key: "status",
    header: "Statut",
    render: (p) => <StatusBadge status={p.status} />,
  },
];

const defaultProducts: Product[] = [
  { id: "PROD-001", name: "Smartphone X200", description: "", price: 450_000, stock: 34, category: "Électronique", status: "active", createdAt: "", updatedAt: "" },
  { id: "PROD-002", name: "Casque Bluetooth Pro", description: "", price: 85_000, stock: 120, category: "Électronique", status: "active", createdAt: "", updatedAt: "" },
  { id: "PROD-003", name: "T-shirt Coton Bio", description: "", price: 12_500, stock: 0, category: "Mode", status: "inactive", createdAt: "", updatedAt: "" },
  { id: "PROD-004", name: "Robot Aspirateur", description: "", price: 295_000, stock: 18, category: "Maison", status: "active", createdAt: "", updatedAt: "" },
  { id: "PROD-005", name: "Huile d'Olive Bio 1L", description: "", price: 8_500, stock: 250, category: "Alimentation", status: "active", createdAt: "", updatedAt: "" },
  { id: "PROD-006", name: "Montre Connectée S3", description: "", price: 175_000, stock: 45, category: "Électronique", status: "active", createdAt: "", updatedAt: "" },
  { id: "PROD-007", name: "Sac à Main Cuir", description: "", price: 95_000, stock: 2, category: "Mode", status: "active", createdAt: "", updatedAt: "" },
  { id: "PROD-008", name: "Set de Rangement", description: "", price: 22_000, stock: 0, category: "Maison", status: "draft", createdAt: "", updatedAt: "" },
  { id: "PROD-009", name: "Crème Hydratante Premium", description: "", price: 18_500, stock: 78, category: "Beauté", status: "active", createdAt: "", updatedAt: "" },
  { id: "PROD-010", name: "Ballon de Foot Pro", description: "", price: 32_000, stock: 56, category: "Sport", status: "active", createdAt: "", updatedAt: "" },
];

const emptyForm: ProductInput = {
  name: "",
  description: "",
  price: 0,
  stock: 0,
  category: "Électronique",
  status: "active",
};

export default function ProductsPage() {
  const [productList, setProductList] = useState<Product[]>(defaultProducts);
  const [page, setPage] = useState(1);
  const [totalPages] = useState(3);
  const [search, setSearch] = useState("");
  const [loading, setLoading] = useState(true);
  const [modalOpen, setModalOpen] = useState(false);
  const [editingProduct, setEditingProduct] = useState<Product | null>(null);
  const [form, setForm] = useState<ProductInput>(emptyForm);
  const [confirmDelete, setConfirmDelete] = useState<string | null>(null);

  useEffect(() => {
    productsApi
      .list({ page, search })
      .then((res) => {
        if (res.data.length > 0) setProductList(res.data);
      })
      .catch(() => {})
      .finally(() => setLoading(false));
  }, [page, search]);

  const openCreate = useCallback(() => {
    setEditingProduct(null);
    setForm(emptyForm);
    setModalOpen(true);
  }, []);

  const openEdit = useCallback((product: Product) => {
    setEditingProduct(product);
    setForm({
      name: product.name,
      description: product.description,
      price: product.price,
      stock: product.stock,
      category: product.category,
      status: product.status,
    });
    setModalOpen(true);
  }, []);

  const handleSave = useCallback(async () => {
    try {
      if (editingProduct) {
        const res = await productsApi.update(editingProduct.id, form);
        setProductList((prev) =>
          prev.map((p) => (p.id === editingProduct.id ? res.data : p)),
        );
      } else {
        const res = await productsApi.create(form);
        setProductList((prev) => [res.data, ...prev]);
      }
      setModalOpen(false);
    } catch {
      // Offline fallback: update local state
      if (editingProduct) {
        setProductList((prev) =>
          prev.map((p) =>
            p.id === editingProduct.id ? { ...p, ...form } : p,
          ),
        );
      } else {
        const fakeProduct: Product = {
          id: `PROD-${String(Date.now()).slice(-4)}`,
          ...form,
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString(),
        };
        setProductList((prev) => [fakeProduct, ...prev]);
      }
      setModalOpen(false);
    }
  }, [editingProduct, form]);

  const handleDelete = useCallback(async (id: string) => {
    try {
      await productsApi.delete(id);
    } catch {
      // fallback
    }
    setProductList((prev) => prev.filter((p) => p.id !== id));
    setConfirmDelete(null);
  }, []);

  const filtered = search
    ? productList.filter(
        (p) =>
          p.name.toLowerCase().includes(search.toLowerCase()) ||
          p.id.toLowerCase().includes(search.toLowerCase()),
      )
    : productList;

  return (
    <PageContainer
      title="Produits"
      actions={
        <button
          onClick={openCreate}
          className="flex items-center gap-1.5 px-4 py-2 bg-accent hover:bg-accent/90 text-white text-sm font-medium rounded-lg transition-colors"
        >
          <Plus size={16} weight="bold" />
          Ajouter
        </button>
      }
    >
      {/* Filters */}
      <div className="flex items-center gap-3 mb-4">
        <div className="w-64">
          <SearchInput value={search} onChange={setSearch} />
        </div>
        <select className="px-3 py-2 text-sm bg-surface border border-border rounded-lg text-text-primary focus:outline-none focus:ring-2 focus:ring-accent/30">
          <option value="">Toutes les catégories</option>
          {categories.map((cat) => (
            <option key={cat} value={cat}>{cat}</option>
          ))}
        </select>
      </div>

      {/* Table */}
      <div className="bg-surface-elevated border border-border rounded-xl overflow-hidden">
        <DataTable
          columns={[
            ...columns,
            {
              key: "actions",
              header: "",
              render: (p) => (
                <div className="flex items-center gap-1">
                  <button
                    onClick={() => openEdit(p)}
                    className="p-1.5 rounded-md text-text-muted hover:text-accent hover:bg-accent/10 transition-colors"
                  >
                    <PencilSimple size={16} />
                  </button>
                  <button
                    onClick={() => setConfirmDelete(p.id)}
                    className="p-1.5 rounded-md text-text-muted hover:text-danger hover:bg-danger/10 transition-colors"
                  >
                    <TrashSimple size={16} />
                  </button>
                </div>
              ),
            },
          ]}
          data={filtered}
          keyExtractor={(p) => p.id}
          page={page}
          totalPages={totalPages}
          onPageChange={setPage}
          loading={loading}
        />
      </div>

      {/* Create/Edit Modal */}
      <Modal
        open={modalOpen}
        onClose={() => setModalOpen(false)}
        title={editingProduct ? "Modifier le produit" : "Nouveau produit"}
      >
        <div className="space-y-4">
          <div>
            <label className="block text-xs font-medium text-text-secondary mb-1">Nom</label>
            <input
              type="text"
              value={form.name}
              onChange={(e) => setForm({ ...form, name: e.target.value })}
              className="w-full px-3 py-2 text-sm bg-surface border border-border rounded-lg text-text-primary focus:outline-none focus:ring-2 focus:ring-accent/30"
              placeholder="Nom du produit"
            />
          </div>
          <div>
            <label className="block text-xs font-medium text-text-secondary mb-1">Description</label>
            <textarea
              value={form.description}
              onChange={(e) => setForm({ ...form, description: e.target.value })}
              className="w-full px-3 py-2 text-sm bg-surface border border-border rounded-lg text-text-primary focus:outline-none focus:ring-2 focus:ring-accent/30 resize-none"
              rows={3}
              placeholder="Description..."
            />
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-xs font-medium text-text-secondary mb-1">Prix (FCFA)</label>
              <input
                type="number"
                value={form.price}
                onChange={(e) => setForm({ ...form, price: Number(e.target.value) })}
                className="w-full px-3 py-2 text-sm bg-surface border border-border rounded-lg text-text-primary focus:outline-none focus:ring-2 focus:ring-accent/30"
              />
            </div>
            <div>
              <label className="block text-xs font-medium text-text-secondary mb-1">Stock</label>
              <input
                type="number"
                value={form.stock}
                onChange={(e) => setForm({ ...form, stock: Number(e.target.value) })}
                className="w-full px-3 py-2 text-sm bg-surface border border-border rounded-lg text-text-primary focus:outline-none focus:ring-2 focus:ring-accent/30"
              />
            </div>
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-xs font-medium text-text-secondary mb-1">Catégorie</label>
              <select
                value={form.category}
                onChange={(e) => setForm({ ...form, category: e.target.value })}
                className="w-full px-3 py-2 text-sm bg-surface border border-border rounded-lg text-text-primary focus:outline-none focus:ring-2 focus:ring-accent/30"
              >
                {categories.map((cat) => (
                  <option key={cat} value={cat}>{cat}</option>
                ))}
              </select>
            </div>
            <div>
              <label className="block text-xs font-medium text-text-secondary mb-1">Statut</label>
              <select
                value={form.status}
                onChange={(e) => setForm({ ...form, status: e.target.value as ProductStatus })}
                className="w-full px-3 py-2 text-sm bg-surface border border-border rounded-lg text-text-primary focus:outline-none focus:ring-2 focus:ring-accent/30"
              >
                <option value="active">Actif</option>
                <option value="inactive">Inactif</option>
                <option value="draft">Brouillon</option>
              </select>
            </div>
          </div>
          <div className="flex justify-end gap-2 pt-2">
            <button
              onClick={() => setModalOpen(false)}
              className="px-4 py-2 text-sm text-text-secondary hover:text-text-primary transition-colors"
            >
              Annuler
            </button>
            <button
              onClick={handleSave}
              className="px-4 py-2 bg-accent hover:bg-accent/90 text-white text-sm font-medium rounded-lg transition-colors"
            >
              {editingProduct ? "Enregistrer" : "Créer"}
            </button>
          </div>
        </div>
      </Modal>

      {/* Delete confirmation */}
      <Modal
        open={confirmDelete !== null}
        onClose={() => setConfirmDelete(null)}
        title="Confirmer la suppression"
        size="sm"
      >
        <p className="text-sm text-text-secondary mb-4">
          Êtes-vous sûr de vouloir supprimer ce produit ? Cette action est irréversible.
        </p>
        <div className="flex justify-end gap-2">
          <button
            onClick={() => setConfirmDelete(null)}
            className="px-4 py-2 text-sm text-text-secondary hover:text-text-primary transition-colors"
          >
            Annuler
          </button>
          <button
            onClick={() => confirmDelete && handleDelete(confirmDelete)}
            className="px-4 py-2 bg-danger hover:bg-danger/90 text-white text-sm font-medium rounded-lg transition-colors"
          >
            Supprimer
          </button>
        </div>
      </Modal>
    </PageContainer>
  );
}

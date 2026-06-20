import type {
  ApiResponse,
  AuthResponse,
  DashboardStats,
  LoginRequest,
  Order,
  OrderStatus,
  PaginatedResponse,
  Payment,
  PaymentStatus,
  Product,
  ProductInput,
  RegisterRequest,
  User,
} from "./types";

/* ─── Config ─── */
const BASE_URL = "/api/v1";

/* ─── Token management ─── */
const TOKEN_KEY = "nexusflow_token";

export function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY);
}

export function setToken(token: string): void {
  localStorage.setItem(TOKEN_KEY, token);
}

export function clearToken(): void {
  localStorage.removeItem(TOKEN_KEY);
}

/* ─── HTTP helpers ─── */
class ApiError extends Error {
  status: number;

  constructor(message: string, status: number) {
    super(message);
    this.name = "ApiError";
    this.status = status;
  }
}

async function request<T>(
  path: string,
  options: RequestInit = {},
): Promise<T> {
  const token = getToken();
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...(options.headers as Record<string, string>),
  };

  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  const response = await fetch(`${BASE_URL}${path}`, {
    ...options,
    headers,
  });

  if (!response.ok) {
    const body = await response.json().catch(() => ({}));
    const message = (body as { error?: string }).error ?? response.statusText;
    throw new ApiError(message, response.status);
  }

  if (response.status === 204) return undefined as T;

  const raw = await response.json() as { success: boolean; data: unknown; error?: string };

  // Unwrap the {success, data, error} envelope
  if (!raw.success) {
    throw new ApiError(raw.error ?? "Unknown error", response.status);
  }

  return raw.data as T;
}

/* ─── Auth API ─── */
export const auth = {
  login: (data: LoginRequest) =>
    request<ApiResponse<AuthResponse>>("/auth/login", {
      method: "POST",
      body: JSON.stringify(data),
    }),

  register: (data: RegisterRequest) =>
    request<ApiResponse<AuthResponse>>("/auth/register", {
      method: "POST",
      body: JSON.stringify(data),
    }),

  getMe: () => request<ApiResponse<User>>("/auth/me"),
};

/* ─── Products API ─── */
export const products = {
  list: async (params?: { page?: number; limit?: number; search?: string; category?: string }) => {
    const q = new URLSearchParams();
    if (params?.page) q.set("page", String(params.page));
    if (params?.limit) q.set("limit", String(params.limit));
    if (params?.search) q.set("search", params.search);
    if (params?.category) q.set("category", params.category);

    const raw = await request<{ products: Product[]; pagination: { total: number; page: number; limit: number; total_pages: number } }>(`/products?${q.toString()}`);

    return {
      data: raw.products.map(adaptProduct),
      total: raw.pagination.total,
      page: raw.pagination.page,
      limit: raw.pagination.limit,
      totalPages: raw.pagination.total_pages,
    } as PaginatedResponse<Product>;
  },

  get: async (id: string) => {
    const raw = await request<Product>(`/products/${id}`);
    return { data: adaptProduct(raw) } as ApiResponse<Product>;
  },

  create: async (data: ProductInput) => {
    const raw = await request<Product>("/products", {
      method: "POST",
      body: JSON.stringify(data),
    });
    return { data: adaptProduct(raw) } as ApiResponse<Product>;
  },

  update: async (id: string, data: Partial<ProductInput>) => {
    const raw = await request<Product>(`/products/${id}`, {
      method: "PUT",
      body: JSON.stringify(data),
    });
    return { data: adaptProduct(raw) } as ApiResponse<Product>;
  },

  delete: (id: string) =>
    request<null>(`/products/${id}`, { method: "DELETE" }),
};

// eslint-disable-next-line @typescript-eslint/no-explicit-any
function adaptProduct(p: any): Product {
  return {
    id: String(p.id ?? ""),
    name: String(p.name ?? ""),
    description: String(p.description ?? ""),
    price: Number(p.price ?? 0),
    stock: Number(p.stock ?? 0),
    category: String(p.category ?? ""),
    status: (p.is_active === false || p.is_active === null) ? "inactive" : "active",
    image: String(p.image_url ?? ""),
    createdAt: String(p.created_at ?? p.createdAt ?? ""),
    updatedAt: String(p.updated_at ?? p.updatedAt ?? ""),
  };
}

/* ─── Orders API ─── */
export const orders = {
  list: async (params?: { page?: number; limit?: number; status?: OrderStatus }) => {
    const q = new URLSearchParams();
    if (params?.page) q.set("page", String(params.page));
    if (params?.limit) q.set("limit", String(params.limit));
    if (params?.status) q.set("status", params.status);

    const raw = await request<{ orders: Record<string, unknown>[]; pagination: { total: number; page: number; limit: number; total_pages: number } }>(`/orders?${q.toString()}`);

    return {
      data: raw.orders.map(adaptOrder),
      total: raw.pagination.total,
      page: raw.pagination.page,
      limit: raw.pagination.limit,
      totalPages: raw.pagination.total_pages,
    } as PaginatedResponse<Order>;
  },

  get: async (id: string) => {
    const raw = await request<Record<string, unknown>>(`/orders/${id}`);
    return { data: adaptOrder(raw) } as ApiResponse<Order>;
  },

  updateStatus: async (id: string, status: OrderStatus) => {
    const raw = await request<Record<string, unknown>>(`/orders/${id}/status`, {
      method: "PUT",
      body: JSON.stringify({ status }),
    });
    return { data: adaptOrder(raw) } as ApiResponse<Order>;
  },
};

// eslint-disable-next-line @typescript-eslint/no-explicit-any
function adaptOrder(o: any): Order {
  return {
    id: String(o.id ?? ""),
    userId: String(o.user_id ?? o.userId ?? ""),
    userName: String(o.user_name ?? o.userName ?? o.user_id ?? "").slice(0, 8),
    items: Array.isArray(o.items) ? o.items.map(adaptItem) : [],
    total: Number(o.total_amount ?? o.total ?? 0),
    status: String(o.status ?? "pending") as OrderStatus,
    paymentMethod: String(o.payment_method ?? o.paymentMethod ?? "wave") as Order["paymentMethod"],
    shippingAddress: String(o.shipping_address ?? o.shippingAddress ?? ""),
    createdAt: String(o.created_at ?? o.createdAt ?? ""),
    updatedAt: String(o.updated_at ?? o.updatedAt ?? ""),
  };
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
function adaptItem(i: any): Order["items"][0] {
  return {
    productId: String(i.product_id ?? i.productId ?? ""),
    productName: String(i.product_name ?? i.productName ?? ""),
    quantity: Number(i.quantity ?? 0),
    unitPrice: Number(i.unit_price ?? i.unitPrice ?? 0),
  };
}

/* ─── Payments API ─── */
export const payments = {
  list: async (params?: { page?: number; limit?: number; status?: PaymentStatus }) => {
    const q = new URLSearchParams();
    if (params?.page) q.set("page", String(params.page));
    if (params?.limit) q.set("limit", String(params.limit));
    if (params?.status) q.set("status", params.status);

    const raw = await request<{ payments: Record<string, unknown>[]; pagination: { total: number; page: number; limit: number; total_pages: number } }>(`/payments?${q.toString()}`);
    // Payments endpoint doesn't exist on the PHP service — fallback to empty
    return {
      data: raw?.payments?.map(adaptPayment) ?? [],
      total: raw?.pagination?.total ?? 0,
      page: raw?.pagination?.page ?? 1,
      limit: raw?.pagination?.limit ?? 20,
      totalPages: raw?.pagination?.total_pages ?? 0,
    } as PaginatedResponse<Payment>;
  },

  process: async (orderId: string, method: string) => {
    const raw = await request<Record<string, unknown>>("/payments", {
      method: "POST",
      body: JSON.stringify({ order_id: orderId, method }),
    });
    return { data: adaptPayment(raw) } as ApiResponse<Payment>;
  },

  refund: async (id: string) => {
    const raw = await request<Record<string, unknown>>(`/payments/${id}/refund`, {
      method: "POST",
    });
    return { data: adaptPayment(raw) } as ApiResponse<Payment>;
  },
};

// eslint-disable-next-line @typescript-eslint/no-explicit-any
function adaptPayment(p: any): Payment {
  return {
    id: String(p.id ?? ""),
    orderId: String(p.order_id ?? p.orderId ?? ""),
    userId: String(p.user_id ?? p.userId ?? ""),
    amount: Number(p.amount ?? 0),
    method: String(p.method ?? p.payment_method ?? "wave") as Payment["method"],
    status: String(p.status ?? "pending") as PaymentStatus,
    transactionRef: String(p.transaction_ref ?? p.transactionRef ?? ""),
    createdAt: String(p.created_at ?? p.createdAt ?? ""),
    updatedAt: String(p.updated_at ?? p.updatedAt ?? ""),
  };
}

/* ─── Users API (admin) ─── */
export const users = {
  list: async (params?: { page?: number; limit?: number }) => {
    const q = new URLSearchParams();
    if (params?.page) q.set("page", String(params.page));
    if (params?.limit) q.set("limit", String(params.limit));

    const raw = await request<{ users: Record<string, unknown>[]; pagination: { total: number; page: number; limit: number; total_pages: number } }>(`/users?${q.toString()}`);

    return {
      data: raw.users.map(adaptUser),
      total: raw.pagination.total,
      page: raw.pagination.page,
      limit: raw.pagination.limit,
      totalPages: raw.pagination.total_pages,
    } as PaginatedResponse<User>;
  },

  toggleStatus: async (id: string) => {
    const raw = await request<Record<string, unknown>>(`/users/${id}/toggle-status`, {
      method: "PATCH",
    });
    return { data: adaptUser(raw) } as ApiResponse<User>;
  },
};

// eslint-disable-next-line @typescript-eslint/no-explicit-any
function adaptUser(u: any): User {
  return {
    id: String(u.id ?? ""),
    email: String(u.email ?? ""),
    name: String(u.name ?? ""),
    role: String(u.role ?? "viewer") as User["role"],
    status: (u.is_active === false || u.is_active === null) ? "inactive" : "active",
    createdAt: String(u.created_at ?? u.createdAt ?? ""),
    updatedAt: String(u.updated_at ?? u.updatedAt ?? ""),
  };
}

/* ─── Dashboard API ─── */
export const dashboard = {
  stats: () => request<DashboardStats>("/dashboard/stats"),
};

/* ─── Health check ─── */
export const health = {
  all: () => request<{ services: { name: string; status: string; uptime: string; latency: number }[] }>("/health"),
};

export { ApiError };

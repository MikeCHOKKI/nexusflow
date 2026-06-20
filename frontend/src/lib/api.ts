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
const BASE_URL = "/api";

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

  return response.json() as Promise<T>;
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
  list: (params?: { page?: number; limit?: number; search?: string; category?: string }) => {
    const q = new URLSearchParams();
    if (params?.page) q.set("page", String(params.page));
    if (params?.limit) q.set("limit", String(params.limit));
    if (params?.search) q.set("search", params.search);
    if (params?.category) q.set("category", params.category);
    return request<PaginatedResponse<Product>>(`/products?${q.toString()}`);
  },

  get: (id: string) => request<ApiResponse<Product>>(`/products/${id}`),

  create: (data: ProductInput) =>
    request<ApiResponse<Product>>("/products", {
      method: "POST",
      body: JSON.stringify(data),
    }),

  update: (id: string, data: Partial<ProductInput>) =>
    request<ApiResponse<Product>>(`/products/${id}`, {
      method: "PUT",
      body: JSON.stringify(data),
    }),

  delete: (id: string) =>
    request<ApiResponse<null>>(`/products/${id}`, { method: "DELETE" }),
};

/* ─── Orders API ─── */
export const orders = {
  list: (params?: { page?: number; limit?: number; status?: OrderStatus }) => {
    const q = new URLSearchParams();
    if (params?.page) q.set("page", String(params.page));
    if (params?.limit) q.set("limit", String(params.limit));
    if (params?.status) q.set("status", params.status);
    return request<PaginatedResponse<Order>>(`/orders?${q.toString()}`);
  },

  get: (id: string) => request<ApiResponse<Order>>(`/orders/${id}`),

  updateStatus: (id: string, status: OrderStatus) =>
    request<ApiResponse<Order>>(`/orders/${id}/status`, {
      method: "PATCH",
      body: JSON.stringify({ status }),
    }),
};

/* ─── Payments API ─── */
export const payments = {
  list: (params?: { page?: number; limit?: number; status?: PaymentStatus }) => {
    const q = new URLSearchParams();
    if (params?.page) q.set("page", String(params.page));
    if (params?.limit) q.set("limit", String(params.limit));
    if (params?.status) q.set("status", params.status);
    return request<PaginatedResponse<Payment>>(`/payments?${q.toString()}`);
  },

  process: (orderId: string, method: string) =>
    request<ApiResponse<Payment>>("/payments/process", {
      method: "POST",
      body: JSON.stringify({ orderId, method }),
    }),

  refund: (id: string) =>
    request<ApiResponse<Payment>>(`/payments/${id}/refund`, {
      method: "POST",
    }),
};

/* ─── Users API (admin) ─── */
export const users = {
  list: (params?: { page?: number; limit?: number }) => {
    const q = new URLSearchParams();
    if (params?.page) q.set("page", String(params.page));
    if (params?.limit) q.set("limit", String(params.limit));
    return request<PaginatedResponse<User>>(`/users?${q.toString()}`);
  },

  toggleStatus: (id: string) =>
    request<ApiResponse<User>>(`/users/${id}/toggle-status`, {
      method: "PATCH",
    }),
};

/* ─── Dashboard API ─── */
export const dashboard = {
  stats: () => request<ApiResponse<DashboardStats>>("/dashboard/stats"),
};

/* ─── Health check ─── */
export const health = {
  all: () => request<ApiResponse<{ services: { name: string; status: string; uptime: string; latency: number }[] }>>("/health"),
};

export { ApiError };

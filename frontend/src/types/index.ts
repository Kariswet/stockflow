export interface Product {
  id: number;
  sku: string;
  name: string;
  price: number;
  stock: number;
  is_active: boolean;
  created_at: string;
  last_movement_at: string | null;
}

export interface BackendResponse<T> {
  metadata: {
    status: boolean;
    message: string;
    time_execution: string;
  };
  data: {
    data: T[]; // Ini array produknya
    page?: number;
    limit?: number;
    total?: number;
  };
}

export interface ApiError {
  error: {
    code: string;
    available?: number;
    sku?: string;
  }
}

export interface OrderItemRequest {
  sku: string;
  qty: number;
}

export interface OrderRequest {
  customer: string;
  items: OrderItemRequest[];
}

export interface OrderItem {
  id: number;
  order_id: number;
  sku: string;
  qty: number;
  unit_price: number;
}

export interface OrderResponse {
  id: number;
  customer: string;
  status: string;
  total: number;
  created_at: string;
  items: OrderItem[];
}
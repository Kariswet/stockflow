import axios from 'axios';
import { type Product, type BackendResponse, type OrderRequest} from '../types';

const client = axios.create({
  baseURL: 'http://localhost:8000/api',
});

export const api = {
  getProducts: async (search: string = '') => {
  
    const res = await client.get<BackendResponse<Product>>('/products', {
      params: { search, page: 1, limit: 100 }
    });
    return res.data; 
  },
  
  adjustStock: async (id: number, payload: { type: 'in' | 'out'; qty: number; note: string }) => {
    try {
      
      const res = await client.post(`/products/${id}/stock`, payload);
      return res.data;
    } catch (error: any) {
      
      if (error.response && error.response.data) {
        throw error.response.data;
      }
      throw error;
    }
  },

  createOrder: async (payload: OrderRequest) => {
    try {

      const res = await client.post('/orders', payload);
      return res.data;
    } catch (error: any) {
      if (error.response && error.response.data) {
        throw error.response.data;
      }
      throw error;
    }
  },
  getOrder: async (id: number) => {
    // Sesuai spesifikasi: GET /api/orders/{id}
    const res = await client.get(`/orders/${id}`);
    return res.data;
  },
  
  cancelOrder: async (id: number) => {
    try {
      // Sesuai spesifikasi: POST /api/orders/{id}/cancel[cite: 7]
      const res = await client.post(`/orders/${id}/cancel`);
      return res.data;
    } catch (error: any) {
      if (error.response && error.response.data) throw error.response.data;
      throw error;
    }
  }
};
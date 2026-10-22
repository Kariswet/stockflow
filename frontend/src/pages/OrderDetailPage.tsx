import { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import { api } from '../api/client';
import { type OrderResponse } from '../types';

export default function OrderDetailPage() {
  const { id } = useParams<{ id: string }>();
  const [order, setOrder] = useState<OrderResponse | null>(null);
  
  const [loading, setLoading] = useState(true);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);
  const [cancelLoading, setCancelLoading] = useState(false);

  const fetchOrder = async () => {
    setLoading(true);
    try {
      const res = await api.getOrder(Number(id));
      // Menangani kedalaman wrapper JSON dari backend Go
      const orderData = res?.data?.id ? res.data : (res?.data?.data || res);
      setOrder(orderData);
    } catch (err: any) {
      setErrorMsg(err.message || 'Failed to fetch order');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (id) fetchOrder();
  }, [id]);

  const handleCancel = async () => {
    // Syarat kelulusan tes: Cancel with confirmation[cite: 7]
    if (!window.confirm('Are you sure you want to cancel this order? Stock will be restored.')) return;

    setCancelLoading(true);
    setErrorMsg(null);
    try {
      await api.cancelOrder(Number(id));
      alert('Order cancelled successfully! Stock is restored.');
      fetchOrder(); // Muat ulang data untuk mengubah UI badge menjadi CANCELLED
    } catch (err: any) {
      const apiErr = err.error || err.data?.error;
      
      // Menangani error 409 ketika pesanan sudah dibatalkan sebelumnya[cite: 7]
      if (apiErr && apiErr.code === 'ORDER_ALREADY_CANCELLED') {
        setErrorMsg('Failed: This order has already been cancelled.');
      } else {
        setErrorMsg(err.message || 'Failed to cancel order');
      }
    } finally {
      setCancelLoading(false);
    }
  };

  if (loading) return <div style={{ padding: '20px' }}>Loading order detail...</div>;
  if (!order) return <div style={{ padding: '20px', color: 'red' }}>Error: {errorMsg || 'Order not found'}</div>;

  return (
    <div style={{ padding: '20px', fontFamily: 'sans-serif' }}>
      <div style={{ marginBottom: '20px' }}>
        <Link to="/orders/new" style={{ textDecoration: 'none', color: 'blue', marginRight: '15px' }}>&larr; Create New Order</Link>
        <Link to="/" style={{ textDecoration: 'none', color: 'blue' }}>View Products</Link>
      </div>

      <h2>Order #{order.id} Detail</h2>
      
      {errorMsg && <div style={{ color: 'red', padding: '10px', border: '1px solid red', marginBottom: '15px' }}>{errorMsg}</div>}

      <div style={{ padding: '15px', background: '#f8f9fa', border: '1px solid #ddd', marginBottom: '20px', borderRadius: '4px' }}>
        <p style={{ margin: '0 0 10px 0' }}><strong>Customer:</strong> {order.customer}</p>
        <p style={{ margin: '0 0 10px 0' }}>
          <strong>Status:</strong> 
          <span style={{ 
            marginLeft: '8px', 
            padding: '4px 8px', 
            borderRadius: '4px', 
            background: order.status === 'CANCELLED' ? '#f8d7da' : '#d4edda',
            color: order.status === 'CANCELLED' ? '#721c24' : '#155724',
            fontWeight: 'bold'
          }}>
            {order.status}
          </span>
        </p>
        <p style={{ margin: '0 0 10px 0' }}><strong>Total:</strong> Rp {order.total?.toLocaleString('id-ID')}</p>
        <p style={{ margin: '0 0 0 0' }}><strong>Date:</strong> {new Date(order.created_at).toLocaleString('id-ID')}</p>
      </div>

      <h3>Items</h3>
      <table border={1} cellPadding={8} style={{ borderCollapse: 'collapse', width: '100%', marginBottom: '20px', textAlign: 'left' }}>
        <thead style={{ background: '#eee' }}>
          <tr>
            <th>SKU</th>
            <th>Qty</th>
            <th>Unit Price</th>
            <th>Subtotal</th>
          </tr>
        </thead>
        <tbody>
          {order.items?.map(item => (
            <tr key={item.id || item.sku}>
              <td>{item.sku}</td>
              <td>{item.qty}</td>
              <td>Rp {item.unit_price?.toLocaleString('id-ID')}</td>
              <td>Rp {(item.qty * item.unit_price)?.toLocaleString('id-ID')}</td>
            </tr>
          ))}
          {(!order.items || order.items.length === 0) && (
            <tr><td colSpan={4} style={{ textAlign: 'center' }}>No items found</td></tr>
          )}
        </tbody>
      </table>

      {order.status !== 'CANCELLED' && (
        <button 
          onClick={handleCancel} 
          disabled={cancelLoading}
          style={{ 
            padding: '10px 20px', 
            background: '#dc3545', 
            color: 'white', 
            border: 'none', 
            borderRadius: '4px',
            cursor: cancelLoading ? 'not-allowed' : 'pointer',
            fontWeight: 'bold'
          }}
        >
          {cancelLoading ? 'Cancelling...' : 'Cancel Order'}
        </button>
      )}
    </div>
  );
}
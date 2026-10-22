import React, { useState } from 'react';
import { api } from '../api/client';
import { type OrderResponse } from '../types';
import { Link, useNavigate } from 'react-router-dom';

export default function CreateOrderPage() {
  const navigate = useNavigate();
  const [customer, setCustomer] = useState('');
  const [items, setItems] = useState([{ sku: '', qty: 1 }]);
  const [lookupId, setLookupId] = useState('');
  
  const [loading, setLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);
  const [successOrder, setSuccessOrder] = useState<OrderResponse | null>(null);

  const handleAddItem = () => setItems([...items, { sku: '', qty: 1 }]);
  
  const handleRemoveItem = (index: number) => {
    const newItems = [...items];
    newItems.splice(index, 1);
    setItems(newItems);
  };

  const handleChangeItem = (index: number, field: 'sku' | 'qty', value: string) => {
    const newItems = [...items];
    if (field === 'qty') {
      newItems[index].qty = value ? parseInt(value, 10) : 0;
    } else {
      newItems[index].sku = value;
    }
    setItems(newItems);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setErrorMsg(null);
    setSuccessOrder(null);
    setLoading(true);

    try {
      const res = await api.createOrder({ customer, items });
      const createdOrder = res?.data?.id ? res.data : (res?.data?.data || res);
      
      setSuccessOrder(createdOrder);
      setCustomer('');
      setItems([{ sku: '', qty: 1 }]);
    } catch (err: any) {
      const apiErr = err.error || err.data?.error;
      
      if (apiErr && apiErr.code === 'INSUFFICIENT_STOCK') {
        setErrorMsg(`Failed: Insufficient stock for SKU "${apiErr.sku}". Only ${apiErr.available} left.`);
      } else if (err.metadata && err.metadata.message) {
        setErrorMsg(`Validation Error: ${err.metadata.message}`);
      } else {
        setErrorMsg(err.message || "Failed to create order");
      }
    } finally {
      setLoading(false);
    }
  };

  const handleLookup = (e: React.FormEvent) => {
    e.preventDefault();
    if (lookupId.trim()) {
      navigate(`/orders/${lookupId.trim()}`);
    }
  };

  return (
    <div style={{ padding: '20px', fontFamily: 'sans-serif' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px' }}>
        <Link to="/" style={{ textDecoration: 'none', color: 'blue' }}>&larr; Back to Products</Link>
        
        <form onSubmit={handleLookup} style={{ display: 'flex', gap: '5px' }}>
          <input 
            type="number" 
            placeholder="Direct to Order ID..." 
            value={lookupId} 
            onChange={e => setLookupId(e.target.value)}
            style={{ padding: '5px' }}
          />
          <button type="submit" style={{ padding: '5px 10px', cursor: 'pointer' }}>Go</button>
        </form>
      </div>

      <h2>Create Sales Order</h2>

      {errorMsg && <div style={{ color: 'red', padding: '10px', border: '1px solid red', marginBottom: '15px' }}>{errorMsg}</div>}
      
      {/* --- INI ADALAH BAGIAN YANG DIPERBAIKI --- */}
      {successOrder && successOrder.id && (
        <div style={{ color: '#155724', backgroundColor: '#d4edda', borderColor: '#c3e6cb', padding: '15px', border: '1px solid', marginBottom: '20px', borderRadius: '5px' }}>
          <h4 style={{ margin: '0 0 10px 0' }}>✅ Order Created Successfully!</h4>
          <p style={{ margin: '0 0 5px 0' }}>Order ID: {successOrder.id}</p>
          <p style={{ margin: '0 0 15px 0' }}>Total: Rp {successOrder.total?.toLocaleString('id-ID')}</p>
          
          <Link 
            to={`/orders/${successOrder.id}`} 
            style={{ 
              display: 'inline-block', 
              padding: '10px 15px', 
              background: '#28a745', 
              color: 'white', 
              textDecoration: 'none', 
              borderRadius: '4px',
              fontWeight: 'bold',
              marginTop: '10px'
            }}
          >
            View Detail & Cancel Order #{successOrder.id} &rarr;
          </Link>
        </div>
      )}
      {/* ---------------------------------------- */}

      <form onSubmit={handleSubmit} style={{ maxWidth: '500px' }}>
        <div style={{ marginBottom: '15px' }}>
          <label style={{ display: 'block', fontWeight: 'bold' }}>Customer Name:</label>
          <input 
            type="text" 
            value={customer} 
            onChange={e => setCustomer(e.target.value)} 
            required 
            style={{ width: '100%', padding: '8px' }}
          />
        </div>

        <div style={{ marginBottom: '15px' }}>
          <label style={{ display: 'block', fontWeight: 'bold' }}>Items:</label>
          {items.map((item, index) => (
            <div key={index} style={{ display: 'flex', gap: '10px', marginBottom: '10px' }}>
              <input 
                type="text" 
                placeholder="SKU" 
                value={item.sku} 
                onChange={e => handleChangeItem(index, 'sku', e.target.value)} 
                required 
                style={{ flex: 2, padding: '8px' }}
              />
              <input 
                type="number" 
                placeholder="Qty" 
                min="1" 
                value={item.qty} 
                onChange={e => handleChangeItem(index, 'qty', e.target.value)} 
                required 
                style={{ flex: 1, padding: '8px' }}
              />
              {items.length > 1 && (
                <button type="button" onClick={() => handleRemoveItem(index)} style={{ padding: '8px', color: 'red' }}>X</button>
              )}
            </div>
          ))}
          <button type="button" onClick={handleAddItem} style={{ padding: '8px' }}>+ Add Another Item</button>
        </div>

        <button type="submit" disabled={loading} style={{ padding: '10px 20px', background: 'blue', color: 'white', border: 'none', cursor: loading ? 'not-allowed' : 'pointer' }}>
          {loading ? 'Processing...' : 'Submit Order'}
        </button>
      </form>
    </div>
  );
}
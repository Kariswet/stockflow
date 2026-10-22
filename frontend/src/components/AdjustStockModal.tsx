import React, { useState } from 'react';
import { api } from '../api/client';
import { type Product } from '../types';

interface ModalProps {
  product: Product;
  onClose: () => void;
  onSuccess: () => void;
}

export default function AdjustStockModal({ product, onClose, onSuccess }: ModalProps) {
  const [type, setType] = useState<'in' | 'out'>('in');
  const [qty, setQty] = useState<number | ''>('');
  const [note, setNote] = useState<string>('');
  
  const [loading, setLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setErrorMsg(null);

    const parsedQty = parseInt(String(qty), 10);
    
    // 1. Validasi: integer > 0[cite: 5]
    if (isNaN(parsedQty) || parsedQty <= 0) {
      setErrorMsg("Quantity must be an integer greater than 0");
      return;
    }

    setLoading(true);
    try {
      await api.adjustStock(product.id, { type, qty: parsedQty, note });
      onSuccess();
      onClose();
    } catch (err: any) {
      // 3. Menampilkan 409 INSUFFICIENT_STOCK di UI[cite: 5]
      // Antisipasi letak objek error (bisa langsung err.error atau di dalam err.data.error)
      const apiErr = err.error || err.data?.error; 
      
      if (apiErr && apiErr.code === 'INSUFFICIENT_STOCK') {
        setErrorMsg(`Failed: Insufficient stock. Only ${apiErr.available} available.`);
      } else {
        setErrorMsg(err.message || "Failed to adjust stock");
      }
    } finally {
      setLoading(false);
    }
  };

  const overlayStyle: React.CSSProperties = {
    position: 'fixed', top: 0, left: 0, right: 0, bottom: 0,
    backgroundColor: 'rgba(0,0,0,0.5)', display: 'flex', alignItems: 'center', justifyContent: 'center'
  };

  return (
    <div style={overlayStyle}>
      <div style={{ background: 'white', padding: '20px', borderRadius: '8px', minWidth: '300px' }}>
        <h3>Adjust Stock: {product.sku}</h3>
        <p>Current Stock: <strong>{product.stock}</strong></p>

        {errorMsg && <div style={{ color: 'red', marginBottom: '10px', padding: '10px', border: '1px solid red' }}>{errorMsg}</div>}

        <form onSubmit={handleSubmit}>
          <div style={{ marginBottom: '10px' }}>
            <label>Direction: </label>
            <select value={type} onChange={e => setType(e.target.value as 'in' | 'out')}>
              <option value="in">IN (Receive)</option>
              <option value="out">OUT (Consume)</option>
            </select>
          </div>
          
          <div style={{ marginBottom: '10px' }}>
            <label>Quantity: </label>
            <input type="number" min="1" value={qty} onChange={e => setQty(e.target.value ? Number(e.target.value) : '')} required />
          </div>

          <div style={{ marginBottom: '10px' }}>
            <label>Note: </label>
            <input type="text" value={note} onChange={e => setNote(e.target.value)} />
          </div>

          <div style={{ display: 'flex', gap: '10px', marginTop: '20px' }}>
            {/* 2. Submit disabled while saving[cite: 5] */}
            <button type="submit" disabled={loading}>
              {loading ? 'Saving...' : 'Save'}
            </button>
            <button type="button" onClick={onClose} disabled={loading}>Cancel</button>
          </div>
        </form>
      </div>
    </div>
  );
}
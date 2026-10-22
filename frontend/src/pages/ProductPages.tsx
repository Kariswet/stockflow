import { useEffect, useState } from 'react';
import { api } from '../api/client';
import { type Product } from '../types';
import AdjustStockModal from '../components/AdjustStockModal';

// Tambahkan di atas export default function ProductsPage() { ... }
function SummaryCards({ products }: { products: Product[] }) {
  // 1. Total Produk
  const totalProducts = products.length;
  
  // 2. Produk Aktif[cite: 5]
  const activeProducts = products.filter(p => p.is_active).length;
  
  // 3. Stok Menipis (<= 10)[cite: 5]
  const lowStockCount = products.filter(p => p.stock <= 10).length;
  
  // 4. Total Nilai Stok (Harga * Stok)[cite: 5]
  const totalStockValue = products.reduce((sum, p) => sum + (p.price * p.stock), 0);

  const cardStyle: React.CSSProperties = {
    padding: '15px',
    background: '#f8f9fa',
    border: '1px solid #dee2e6',
    borderRadius: '6px',
    flex: '1',
    minWidth: '150px',
    boxShadow: '0 2px 4px rgba(0,0,0,0.05)'
  };

  const titleStyle: React.CSSProperties = { margin: '0 0 10px 0', fontSize: '14px', color: '#6c757d' };
  const valueStyle: React.CSSProperties = { margin: '0', fontSize: '24px', fontWeight: 'bold', color: '#212529' };

  return (
    <div style={{ display: 'flex', gap: '15px', marginBottom: '25px', flexWrap: 'wrap' }}>
      <div style={cardStyle}>
        <h4 style={titleStyle}>Total Products</h4>
        <p style={valueStyle}>{totalProducts}</p>
      </div>
      <div style={cardStyle}>
        <h4 style={titleStyle}>Active Products</h4>
        <p style={valueStyle}>{activeProducts}</p>
      </div>
      <div style={cardStyle}>
        <h4 style={titleStyle}>Low Stock (≤ 10)</h4>
        <p style={{ ...valueStyle, color: lowStockCount > 0 ? '#dc3545' : '#28a745' }}>
          {lowStockCount}
        </p>
      </div>
      <div style={cardStyle}>
        <h4 style={titleStyle}>Total Stock Value</h4>
        <p style={{ ...valueStyle, fontSize: '20px' }}>
          Rp {totalStockValue.toLocaleString('id-ID')}
        </p>
      </div>
    </div>
  );
}

export default function ProductsPage() {
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [search, setSearch] = useState('');
  const [filter, setFilter] = useState('ALL');

  const [selectedProduct, setSelectedProduct] = useState<Product | null>(null);

  const fetchProducts = async () => {
    setLoading(true);
    setError(null);
    try {
      const res = await api.getProducts(search); 
      
      // PERBAIKAN: Menarik array produk dari struktur JSON backend yang sebenarnya
      let mappedData = res?.data?.data || [];

      // Filter manual di frontend sesuai spesifikasi
      if (filter === 'ACTIVE') mappedData = mappedData.filter((p: Product) => p.is_active);
      if (filter === 'INACTIVE') mappedData = mappedData.filter((p: Product) => !p.is_active);
      if (filter === 'LOW_STOCK') mappedData = mappedData.filter((p: Product) => p.stock <= 10);

      setProducts(mappedData);
    } catch (err: any) {
      setError(err.message || 'Gagal mengambil data');
    } finally {
      setLoading(false);
    }
  };

  // Efek debounce sederhana bawaan react
  useEffect(() => {
    const timer = setTimeout(() => {
      fetchProducts();
    }, 500);
    return () => clearTimeout(timer);
  }, [search, filter]);

  return (
    <div style={{ padding: '20px', fontFamily: 'sans-serif' }}>
      <h2>Produk (Mini ERP)</h2>

      {!loading && !error && products.length > 0 && (
        <SummaryCards products={products} />
      )}
      
      <div style={{ marginBottom: '15px' }}>
        <input 
          placeholder="Cari SKU / Nama..." 
          value={search} onChange={e => setSearch(e.target.value)} 
          style={{ marginRight: '10px', padding: '5px' }}
        />
        <select value={filter} onChange={e => setFilter(e.target.value)} style={{ padding: '5px' }}>
          <option value="ALL">All</option>
          <option value="ACTIVE">Active</option>
          <option value="INACTIVE">Inactive</option>
          <option value="LOW_STOCK">Low Stock (≤ 10)</option>
        </select>
      </div>

      {error && <div><p style={{ color: 'red' }}>Error: {error}</p><button onClick={fetchProducts}>Retry</button></div>}
      {loading && <p>Loading...</p>}
      {!loading && !error && products.length === 0 && <p>Data Kosong.</p>}

      {!loading && !error && products.length > 0 && (
        <table border={1} cellPadding={8} style={{ borderCollapse: 'collapse', width: '100%' }}>
          <thead style={{ background: '#eee' }}>
            <tr>
              <th>SKU</th><th>Nama</th><th>Harga</th><th>Stok</th><th>Status</th><th>Last Mvmt</th><th>Aksi</th>
            </tr>
          </thead>
          <tbody>
            {products.map(p => (
              <tr key={p.id}>
                <td>{p.sku}</td>
                <td>{p.name}</td>
                <td>Rp {p.price.toLocaleString('id-ID')}</td>
                <td>{p.stock}</td>
                <td>{p.is_active ? '✅ Active' : '❌ Inactive'}</td>
                <td>{p.last_movement_at ? new Date(p.last_movement_at).toLocaleDateString() : '-'}</td>
                <td><button onClick={() => setSelectedProduct(p)}>Adjust Stock</button></td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
      {selectedProduct && (
        <AdjustStockModal 
          product={selectedProduct} 
          onClose={() => setSelectedProduct(null)} 
          onSuccess={() => {
            fetchProducts(); // Refresh tabel otomatis jika sukses
          }}
        />
      )}
    </div>
  );
}
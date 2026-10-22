import { BrowserRouter, Routes, Route, Link } from 'react-router-dom';
import ProductsPage from './pages/ProductPages';
import CreateOrderPage from './pages/CreateOrderPAge';
import OrderDetailPage from './pages/OrderDetailPage'; 

export default function App() {
  return (
    <BrowserRouter>
      <nav style={{ padding: '15px 20px', background: '#333', color: 'white' }}>
        <Link to="/" style={{ color: 'white', marginRight: '20px', textDecoration: 'none' }}>Products</Link>
        <Link to="/orders/new" style={{ color: 'white', textDecoration: 'none' }}>Create Order</Link>
      </nav>

      <Routes>
        <Route path="/" element={<ProductsPage />} />
        <Route path="/orders/new" element={<CreateOrderPage />} />
        {/* Tambahkan rute dinamis untuk detail pesanan */}
        <Route path="/orders/:id" element={<OrderDetailPage />} /> 
      </Routes>
    </BrowserRouter>
  );
}
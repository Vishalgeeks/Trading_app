import { createBrowserRouter } from 'react-router-dom';
import MainLayout from '../layouts/MainLayout';
import Home from '../pages/Home';
import Login from '../pages/Login';
import Register from '../pages/Register';
import Designs from '../pages/Designs';
import Bookings from '../pages/Bookings';
import AdminDashboard from '../pages/AdminDashboard';

const router = createBrowserRouter([
  {
    path: '/',
    element: <MainLayout />,
    children: [
      { index: true, element: <Home /> },
      { path: 'login', element: <Login /> },
      { path: 'register', element: <Register /> },
      { path: 'designs', element: <Designs /> },
      { path: 'bookings', element: <Bookings /> },
      { path: 'admin', element: <AdminDashboard /> },
    ],
  },
]);

export default router;

import { createBrowserRouter } from 'react-router-dom';
import ClientLayout from "../layouts/ClientLayout";
import AdminLayout from "../layouts/AdminLayout";
import Home from "../pages/Home";
import Login from "../pages/Login";
import Register from "../pages/Register";
import Browse from "../pages/Browse";
import DesignDetails from "../pages/DesignDetails";
import Favorites from "../pages/Favorites";
import Bookings from "../pages/Bookings";
import BookingDetails from "../pages/BookingDetails";
import Booking from "../pages/Booking";
import Profile from "../pages/Profile";
import Settings from "../pages/Settings";
import AdminBookings from "../pages/AdminBookings";
import AdminBookingDetails from "../pages/AdminBookingDetails";
import AdminDesigns from "../pages/AdminDesigns";
import AdminCategories from "../pages/AdminCategories";
import AdminAvailability from "../pages/AdminAvailability";
import AdminDashboard from "../pages/AdminDashboard";
import ProtectedRoute from "./ProtectedRoute";
import AdminRoute from "./AdminRoute";

const router = createBrowserRouter([
  {
    path: "/",
    element: <ClientLayout />,
    children: [
      { index: true, element: <Home /> },
      { path: "login", element: <Login /> },
      { path: "register", element: <Register /> },
      { path: "browse", element: <Browse /> },
      { path: "designs", element: <Browse /> },
      { path: "designs/:id", element: <DesignDetails /> },
      {
        path: "booking/:id",
        element: (
          <ProtectedRoute>
            <Booking />
          </ProtectedRoute>
        ),
      },
      { path: "favorites", element: <Favorites /> },
      {
        path: "bookings",
        element: (
          <ProtectedRoute>
            <Bookings />
          </ProtectedRoute>
        ),
      },
      {
        path: "bookings/:id",
        element: (
          <ProtectedRoute>
            <BookingDetails />
          </ProtectedRoute>
        ),
      },
      {
        path: "profile",
        element: (
          <ProtectedRoute>
            <Profile />
          </ProtectedRoute>
        ),
      },
      {
        path: "settings",
        element: (
          <ProtectedRoute>
            <Settings />
          </ProtectedRoute>
        ),
      },
    ],
  },
  {
    path: "/admin",
    element: (
      <AdminRoute>
        <AdminLayout />
      </AdminRoute>
    ),
    children: [
      { index: true, element: <AdminDashboard /> },
      { path: "bookings", element: <AdminBookings /> },
      { path: "bookings/:id", element: <AdminBookingDetails /> },
      { path: "designs", element: <AdminDesigns /> },
      { path: "categories", element: <AdminCategories /> },
      { path: "availability", element: <AdminAvailability /> },
    ],
  },
]);

export default router;
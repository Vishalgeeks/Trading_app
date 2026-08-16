import { Link, useLocation } from 'react-router-dom';

export default function Navbar() {
  const location = useLocation();
  
  const isActive = (path) => location.pathname === path ? 'text-rose-600 font-semibold' : 'text-gray-600 hover:text-rose-500';

  return (
    <nav className="bg-white shadow-sm border-b border-gray-100">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between h-16">
          <div className="flex items-center">
            <Link to="/" className="text-2xl font-bold text-rose-500">
              Henna Booking
            </Link>
          </div>
          
          <div className="hidden md:flex items-center space-x-8">
            <Link to="/" className={isActive('/')}>Home</Link>
            <Link to="/designs" className={isActive('/designs')}>Designs</Link>
            <Link to="/bookings" className={isActive('/bookings')}>Bookings</Link>
            <Link to="/login" className={isActive('/login')}>Login</Link>
          </div>

          <div className="md:hidden flex items-center">
            <button className="text-gray-600 hover:text-rose-500">
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
              </svg>
            </button>
          </div>
        </div>
      </div>
    </nav>
  );
}

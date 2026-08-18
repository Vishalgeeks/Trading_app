import { useEffect, useState, useMemo } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { bookingService } from '../services/bookingService';
import BookingCard from '../components/bookings/BookingCard';
import { Plus, Filter, ChevronDown, Loader2 } from 'lucide-react';

const STATUS_OPTIONS = [
  { value: '', label: 'All' },
  { value: 'PENDING', label: 'Pending' },
  { value: 'CONFIRMED', label: 'Confirmed' },
  { value: 'CANCELLED', label: 'Cancelled' },
  { value: 'COMPLETED', label: 'Completed' }
];

export default function Bookings() {
  const { user } = useAuth();
  const navigate = useNavigate();
  const [bookings, setBookings] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [statusFilter, setStatusFilter] = useState('');
  const [total, setTotal] = useState(0);

  useEffect(() => {
    let cancelled = false;

    async function loadBookings() {
      setLoading(true);
      setError('');

      const result = await bookingService.getMyBookings({ status: statusFilter });
      if (cancelled) return;

      if (result.error) {
        if (result.status === 401) {
          navigate('/login');
          return;
        }
        setError(result.message || 'Failed to load bookings');
        setBookings([]);
        setTotal(0);
      } else {
        setBookings(result.data);
        setTotal(result.total || result.data.length);
      }
      setLoading(false);
    }

    loadBookings();

    return () => {
      cancelled = true;
    };
  }, [statusFilter, navigate]);

  const handleFilterChange = (value) => {
    setStatusFilter(value);
  };

  const filteredBookings = useMemo(() => {
    return bookings;
  }, [bookings]);

  if (loading) {
    return (
      <div className="pb-24">
        <div className="px-5 pt-6">
          <div className="flex items-center justify-between mb-6">
            <div>
              <h1 className="text-2xl font-bold text-gray-900 dark:text-white mb-1">My Bookings</h1>
              <p className="text-sm text-gray-500 dark:text-neutral-400">Loading bookings...</p>
            </div>
          </div>
          <div className="space-y-3">
            {[1, 2, 3].map((i) => (
              <div key={i} className="h-24 bg-gray-200 dark:bg-neutral-800 rounded-2xl animate-pulse" />
            ))}
          </div>
        </div>
      </div>
    );
  }

  const upcomingBookings = bookings.filter(b => 
    b.status === 'PENDING' || b.status === 'CONFIRMED'
  );
  const pastBookings = bookings.filter(b => 
    b.status === 'COMPLETED' || b.status === 'CANCELLED'
  );

  return (
    <div className="pb-24">
      <div className="px-5 pt-6">
        <div className="flex items-center justify-between mb-6">
          <div>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-white mb-1">My Bookings</h1>
            <p className="text-sm text-gray-500 dark:text-neutral-400">
              {total} booking{total !== 1 ? 's' : ''}
            </p>
          </div>
          <Link to="/booking" className="flex items-center gap-2 px-4 py-2 bg-rose-500 text-white rounded-xl text-sm font-medium hover:bg-rose-600 transition-colors">
            <Plus size={18} />
            Book New
          </Link>
        </div>

        <div className="flex items-center gap-3 mb-4 flex-wrap">
          <label className="flex items-center gap-2 text-sm text-gray-700 dark:text-neutral-300">
            <Filter size={16} className="text-gray-400" />
            <span>Filter by status:</span>
          </label>
          <div className="relative">
            <select
              value={statusFilter}
              onChange={(e) => handleFilterChange(e.target.value)}
              className="appearance-none bg-white dark:bg-neutral-800 border border-rose-200 dark:border-neutral-700 rounded-xl px-4 py-2 pr-10 text-sm text-gray-900 dark:text-white focus:ring-2 focus:ring-rose-500 focus:border-transparent outline-none cursor-pointer"
            >
              {STATUS_OPTIONS.map((opt) => (
                <option key={opt.value} value={opt.value}>{opt.label}</option>
              ))}
            </select>
            <ChevronDown className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 pointer-events-none" size={16} />
          </div>
        </div>

        {error && (
          <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 px-4 py-3 rounded-xl mb-4 text-sm">
            {error}
          </div>
        )}

        {filteredBookings.length > 0 ? (
          <>
            {upcomingBookings.length > 0 && (
              <div className="mb-6">
                <h2 className="text-sm font-semibold text-gray-900 dark:text-white uppercase tracking-wide mb-3">
                  Upcoming
                </h2>
                <div className="space-y-3">
                  {upcomingBookings.map((booking) => (
                    <Link key={booking.id} to={`/bookings/${booking.id}`}>
                      <BookingCard booking={booking} />
                    </Link>
                  ))}
                </div>
              </div>
            )}

            {pastBookings.length > 0 && (
              <div>
                <h2 className="text-sm font-semibold text-gray-900 dark:text-white uppercase tracking-wide mb-3">
                  Past
                </h2>
                <div className="space-y-3">
                  {pastBookings.map((booking) => (
                    <Link key={booking.id} to={`/bookings/${booking.id}`}>
                      <BookingCard booking={booking} />
                    </Link>
                  ))}
                </div>
              </div>
            )}
          </>
        ) : (
          <div className="text-center py-16">
            <div className="w-16 h-16 bg-rose-50 dark:bg-orange-950 rounded-full flex items-center justify-center mx-auto mb-4">
              <Plus size={28} className="text-rose-400 dark:text-orange-400" />
            </div>
            <p className="text-gray-500 dark:text-neutral-400 text-sm mb-4">
              {statusFilter ? 'No bookings match this filter.' : 'No bookings yet'}
            </p>
            <Link to="/browse" className="text-sm text-rose-500 dark:text-orange-400 font-medium">
              {statusFilter ? 'Clear filter' : 'Browse designs to book'}
            </Link>
          </div>
        )}
      </div>
    </div>
  );
}
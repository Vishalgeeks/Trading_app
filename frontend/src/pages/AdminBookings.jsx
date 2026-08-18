import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { bookingService } from '../services/bookingService';
import { Search } from 'lucide-react';
import BookingStatus from '../components/bookings/BookingStatus';

const STATUS_OPTIONS = [
  { value: '', label: 'All Statuses' },
  { value: 'PENDING', label: 'Pending' },
  { value: 'CONFIRMED', label: 'Confirmed' },
  { value: 'CANCELLED', label: 'Cancelled' },
  { value: 'COMPLETED', label: 'Completed' }
];

export default function AdminBookings() {
  const [bookings, setBookings] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [statusFilter, setStatusFilter] = useState('');
  const [dateFilter, setDateFilter] = useState('');
  const [searchFilter, setSearchFilter] = useState('');
  const [fromDate, setFromDate] = useState('');
  const [toDate, setToDate] = useState('');
  const [total, setTotal] = useState(0);

  useEffect(() => {
    let cancelled = false;

    async function loadBookings() {
      setLoading(true);
      setError('');

      const result = await bookingService.getAdminBookings({
        status: statusFilter,
        date: dateFilter,
        from: fromDate,
        to: toDate,
        search: searchFilter,
        limit: 50
      });
      if (cancelled) return;

      if (result.error) {
        if (result.status === 401) {
          // Redirect will be handled by ProtectedRoute/AdminRoute
          setError('Unauthorized');
        } else {
          setError(result.message || 'Failed to load bookings');
        }
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
  }, [statusFilter, dateFilter, fromDate, toDate, searchFilter]);

  const handleClearFilters = () => {
    setStatusFilter('');
    setDateFilter('');
    setFromDate('');
    setToDate('');
    setSearchFilter('');
  };

  const hasActiveFilters = statusFilter || dateFilter || fromDate || toDate || searchFilter;

  if (loading && bookings.length === 0) {
    return (
      <div className="p-5">
        <div className="flex items-center justify-between mb-6">
          <div>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-white mb-1">All Bookings</h1>
            <p className="text-sm text-gray-500 dark:text-neutral-400">Loading bookings...</p>
          </div>
        </div>
        <div className="space-y-3">
          {[1, 2, 3].map((i) => (
            <div key={i} className="h-24 bg-gray-200 dark:bg-neutral-800 rounded-2xl animate-pulse" />
          ))}
        </div>
      </div>
    );
  }

  return (
    <div className="p-5">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white mb-1">All Bookings</h1>
          <p className="text-sm text-gray-500 dark:text-neutral-400">{total} booking{total !== 1 ? 's' : ''}</p>
        </div>
      </div>

      {/* Filters */}
      <div className="bg-white dark:bg-neutral-900 rounded-2xl border border-rose-100 dark:border-neutral-800 p-5 shadow-sm mb-6">
        {/* Status Filter */}
        <div className="mb-4">
          <label className="block text-sm font-medium text-gray-700 dark:text-neutral-300 mb-2">
            Status Filter
          </label>
          <div className="flex flex-wrap gap-2">
            {STATUS_OPTIONS.map((opt) => (
              <button
                key={opt.value}
                onClick={() => setStatusFilter(opt.value)}
                className={`px-4 py-2 rounded-full text-sm font-medium whitespace-nowrap transition-all duration-200 ${
                  statusFilter === opt.value
                    ? 'bg-rose-500 text-white dark:bg-orange-500'
                    : 'bg-white dark:bg-neutral-800 text-gray-600 dark:text-neutral-400 border border-rose-100 dark:border-neutral-700'
                }`}
              >
                {opt.label}
              </button>
            ))}
          </div>
        </div>

        {/* Date Filter */}
        <div className="mb-4">
          <label className="block text-sm font-medium text-gray-700 dark:text-neutral-300 mb-2">
            Date Filter
          </label>
          <div className="flex flex-wrap gap-2">
            <input
              type="date"
              value={dateFilter}
              onChange={(e) => setDateFilter(e.target.value)}
              placeholder="Select date"
              className="px-4 py-2 bg-white dark:bg-neutral-800 border border-rose-200 dark:border-neutral-700 rounded-xl text-sm text-gray-900 dark:text-white focus:ring-2 focus:ring-rose-500 focus:border-transparent outline-none"
            />
            <input
              type="date"
              value={fromDate}
              onChange={(e) => setFromDate(e.target.value)}
              placeholder="From"
              className="px-4 py-2 bg-white dark:bg-neutral-800 border border-rose-200 dark:border-neutral-700 rounded-xl text-sm text-gray-900 dark:text-white focus:ring-2 focus:ring-rose-500 focus:border-transparent outline-none"
            />
            <input
              type="date"
              value={toDate}
              onChange={(e) => setToDate(e.target.value)}
              placeholder="To"
              className="px-4 py-2 bg-white dark:bg-neutral-800 border border-rose-200 dark:border-neutral-700 rounded-xl text-sm text-gray-900 dark:text-white focus:ring-2 focus:ring-rose-500 focus:border-transparent outline-none"
            />
          </div>
        </div>

        {/* Search Filter */}
        <div className="mb-4">
          <label className="block text-sm font-medium text-gray-700 dark:text-neutral-300 mb-2">
            Search
          </label>
          <div className="relative">
            <Search size={18} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
            <input
              type="text"
              value={searchFilter}
              onChange={(e) => setSearchFilter(e.target.value)}
              placeholder="Search by client name, email, or design..."
              className="w-full pl-10 pr-4 py-2 bg-white dark:bg-neutral-800 border border-rose-200 dark:border-neutral-700 rounded-xl text-sm text-gray-900 dark:text-white focus:ring-2 focus:ring-rose-500 focus:border-transparent outline-none"
            />
          </div>
        </div>

        {hasActiveFilters && (
          <button
            onClick={handleClearFilters}
            className="text-sm text-rose-500 dark:text-orange-400 font-medium hover:underline"
          >
            Clear all filters
          </button>
        )}
      </div>

      {error && (
        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 px-4 py-3 rounded-xl mb-4 text-sm">
          {error}
        </div>
      )}

      {bookings.length > 0 ? (
        <div className="space-y-3">
          {bookings.map((booking) => (
            <Link key={booking.id} to={`/admin/bookings/${booking.id}`} className="block">
              <div className="bg-white dark:bg-neutral-900 rounded-2xl border border-rose-100 dark:border-neutral-800 p-5 shadow-sm hover:shadow-md transition-shadow">
                <div className="flex items-start justify-between mb-3">
                  <div>
                    <h3 className="font-semibold text-gray-900 dark:text-white text-sm">
                      {booking.designName}
                    </h3>
                    <p className="text-xs text-gray-500 dark:text-neutral-400 mt-0.5">
                      {booking.clientName} &bull; {booking.clientEmail}
                    </p>
                  </div>
                  <BookingStatus status={booking.status} />
                </div>
                <div className="flex items-center gap-4 text-xs text-gray-600 dark:text-neutral-400">
                  <div className="flex items-center gap-1">
                    <span>📅</span>
                    <span>{booking.booking_date}</span>
                  </div>
                  <div className="flex items-center gap-1">
                    <span>🕐</span>
                    <span>{booking.start_time} - {booking.end_time}</span>
                  </div>
                  <div className="flex items-center gap-1">
                    <span>🆔</span>
                    <span className="font-mono">{booking.id.slice(0, 8)}...</span>
                  </div>
                </div>
              </div>
            </Link>
          ))}
        </div>
      ) : (
        <div className="text-center py-12">
          <p className="text-gray-500 dark:text-neutral-400">
            {hasActiveFilters ? 'No bookings match this filter.' : 'No bookings found.'}
          </p>
          {hasActiveFilters && (
            <button
              onClick={handleClearFilters}
              className="mt-2 text-sm text-rose-500 dark:text-orange-400 font-medium hover:underline"
            >
              Clear filters to see all bookings
            </button>
          )}
        </div>
      )}
    </div>
  );
}
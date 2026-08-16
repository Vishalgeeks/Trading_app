import { useState } from 'react';
import { mockAdminBookings } from '../mockData';
import AppointmentCard from '../components/admin/AppointmentCard';

export default function AdminBookings() {
  const [filter, setFilter] = useState('all');

  const filtered = filter === 'all'
    ? mockAdminBookings
    : mockAdminBookings.filter((b) => b.status === filter);

  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-900 dark:text-white mb-1">Bookings</h1>
      <p className="text-sm text-gray-500 dark:text-neutral-400 mb-6">
        Manage all client appointments
      </p>

      <div className="flex gap-2 overflow-x-auto pb-3 -mx-4 px-4 scrollbar-hide mb-4">
        {['all', 'PENDING', 'CONFIRMED', 'COMPLETED', 'CANCELLED'].map((f) => (
          <button
            key={f}
            onClick={() => setFilter(f)}
            className={`px-4 py-2 rounded-full text-sm font-medium whitespace-nowrap transition-all duration-200 ${
              filter === f
                ? 'bg-rose-500 text-white dark:bg-orange-500'
                : 'bg-white dark:bg-neutral-800 text-gray-600 dark:text-neutral-400 border border-rose-100 dark:border-neutral-700'
            }`}
          >
            {f === 'all' ? 'All' : f}
          </button>
        ))}
      </div>

      <div className="space-y-3">
        {filtered.map((booking) => (
          <AppointmentCard key={booking.id} appointment={booking} />
        ))}
      </div>

      {filtered.length === 0 && (
        <div className="text-center py-12">
          <p className="text-gray-500 dark:text-neutral-400">No bookings found.</p>
        </div>
      )}
    </div>
  );
}

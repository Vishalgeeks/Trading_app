import { Link } from 'react-router-dom';
import { mockBookings } from '../mockData';
import BookingCard from '../components/bookings/BookingCard';
import { Plus } from 'lucide-react';

export default function Bookings() {
  return (
    <div className="pb-24">
      <div className="px-5 pt-6">
        <h1 className="text-2xl font-bold text-gray-900 dark:text-white mb-1">My Bookings</h1>
        <p className="text-sm text-gray-500 dark:text-neutral-400 mb-6">
          {mockBookings.length} booking{mockBookings.length !== 1 ? 's' : ''}
        </p>

        <div className="space-y-3">
          {mockBookings.map((booking) => (
            <Link key={booking.id} to={`/bookings/${booking.id}`}>
              <BookingCard booking={booking} />
            </Link>
          ))}
        </div>

        {mockBookings.length === 0 && (
          <div className="text-center py-16">
            <div className="w-16 h-16 bg-rose-50 dark:bg-orange-950 rounded-full flex items-center justify-center mx-auto mb-4">
              <Plus size={28} className="text-rose-400 dark:text-orange-400" />
            </div>
            <p className="text-gray-500 dark:text-neutral-400 text-sm mb-4">No bookings yet</p>
            <Link to="/browse" className="text-sm text-rose-500 dark:text-orange-400 font-medium">
              Browse designs to book
            </Link>
          </div>
        )}
      </div>
    </div>
  );
}

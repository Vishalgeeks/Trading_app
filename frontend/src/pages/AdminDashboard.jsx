import { useEffect, useState } from 'react';
import { bookingService } from '../services/bookingService';
import StatCard from '../components/admin/StatCard';
import AppointmentCard from '../components/admin/AppointmentCard';
import {
  CalendarCheck,
  Clock,
  TrendingUp,
  Users,
} from 'lucide-react';

export default function AdminDashboard() {
  const [stats, setStats] = useState(null);
  const [upcoming, setUpcoming] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    async function load() {
      setLoading(true);
      setError('');

      const [statsResult, upcomingResult] = await Promise.all([
        bookingService.getAdminBookingStats(),
        bookingService.getAdminUpcomingBookings(),
      ]);

      if (statsResult.error || upcomingResult.error) {
        setError(statsResult.message || upcomingResult.message || 'Failed to load dashboard');
        setLoading(false);
        return;
      }

      setStats(statsResult.data || {});
      setUpcoming(upcomingResult.data || []);
      setLoading(false);
    }

    load();
  }, []);

  if (loading) {
    return (
      <div>
        <div className="mb-6">
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white mb-1">Dashboard</h1>
          <p className="text-sm text-gray-500 dark:text-neutral-400">Loading...</p>
        </div>
        <div className="grid grid-cols-2 gap-3 mb-8">
          {[1, 2, 3, 4].map((i) => (
            <div key={i} className="h-28 bg-gray-200 dark:bg-neutral-800 rounded-2xl animate-pulse" />
          ))}
        </div>
        <div className="space-y-3">
          {[1, 2, 3].map((i) => (
            <div key={i} className="h-24 bg-gray-200 dark:bg-neutral-800 rounded-2xl animate-pulse" />
          ))}
        </div>
      </div>
    );
  }

  const total = stats?.total_bookings || 0;
  const pending = stats?.pending_bookings || 0;
  const confirmed = stats?.confirmed_bookings || 0;
  const completed = stats?.completed_bookings || 0;
  const _cancelled = stats?.cancelled_bookings || 0;

  return (
    <div>
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-gray-900 dark:text-white mb-1">Dashboard</h1>
        <p className="text-sm text-gray-500 dark:text-neutral-400">
          {new Date().toLocaleDateString('en-US', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })}
        </p>
      </div>

      {error && (
        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 px-4 py-3 rounded-xl mb-4 text-sm">
          {error}
        </div>
      )}

      <div className="grid grid-cols-2 gap-3 mb-8">
        <StatCard label="Total Bookings" value={total} color="orange" icon={CalendarCheck} />
        <StatCard label="Upcoming" value={pending + confirmed} color="blue" icon={Clock} />
        <StatCard label="Confirmed" value={confirmed} color="green" icon={TrendingUp} />
        <StatCard label="Completed" value={completed} color="purple" icon={Users} />
      </div>

      <div>
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-bold text-gray-900 dark:text-white">Upcoming Appointments</h2>
          <span className="text-xs text-gray-500 dark:text-neutral-400">
            {upcoming.length} scheduled
          </span>
        </div>
        <div className="space-y-3">
          {upcoming.map((appointment) => (
            <AppointmentCard key={appointment.id} appointment={appointment} />
          ))}
        </div>
        {upcoming.length === 0 && (
          <div className="text-center py-8">
            <p className="text-gray-500 dark:text-neutral-400 text-sm">No upcoming appointments.</p>
          </div>
        )}
      </div>
    </div>
  );
}

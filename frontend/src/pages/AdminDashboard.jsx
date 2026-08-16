import { mockAdminBookings } from '../mockData';
import StatCard from '../components/admin/StatCard';
import AppointmentCard from '../components/admin/AppointmentCard';
import {
  CalendarCheck,
  Clock,
  TrendingUp,
} from 'lucide-react';

export default function AdminDashboard() {
  const stats = {
    total: mockAdminBookings.length,
    pending: mockAdminBookings.filter((b) => b.status === 'PENDING').length,
    confirmed: mockAdminBookings.filter((b) => b.status === 'CONFIRMED').length,
    completed: mockAdminBookings.filter((b) => b.status === 'COMPLETED').length,
    cancelled: mockAdminBookings.filter((b) => b.status === 'CANCELLED').length,
    upcoming: mockAdminBookings.filter((b) => b.status === 'CONFIRMED' || b.status === 'PENDING').length,
  };

  const upcomingAppointments = mockAdminBookings
    .filter((b) => b.status === 'CONFIRMED' || b.status === 'PENDING')
    .slice(0, 3);

  return (
    <div>
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-gray-900 dark:text-white mb-1">Dashboard</h1>
        <p className="text-sm text-gray-500 dark:text-neutral-400">
          {new Date().toLocaleDateString('en-US', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })}
        </p>
      </div>

      <div className="grid grid-cols-2 gap-3 mb-8">
        <StatCard label="Total Bookings" value={stats.total} color="orange" icon={CalendarCheck} />
        <StatCard label="Upcoming" value={stats.upcoming} color="blue" icon={Clock} />
        <StatCard label="Confirmed" value={stats.confirmed} color="green" icon={TrendingUp} />
        <StatCard label="Completed" value={stats.completed} color="purple" icon={CalendarCheck} />
      </div>

      <div>
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-bold text-gray-900 dark:text-white">Today's Appointments</h2>
          <span className="text-xs text-gray-500 dark:text-neutral-400">
            {upcomingAppointments.length} scheduled
          </span>
        </div>
        <div className="space-y-3">
          {upcomingAppointments.map((appointment) => (
            <AppointmentCard key={appointment.id} appointment={appointment} />
          ))}
        </div>
        {upcomingAppointments.length === 0 && (
          <div className="text-center py-8">
            <p className="text-gray-500 dark:text-neutral-400 text-sm">No upcoming appointments today.</p>
          </div>
        )}
      </div>
    </div>
  );
}

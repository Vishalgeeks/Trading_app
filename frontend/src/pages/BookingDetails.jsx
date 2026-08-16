import { useParams, useNavigate } from 'react-router-dom';
import { mockDesigns, mockBookings } from '../mockData';
import Button from '../components/ui/Button';
import Badge from '../components/ui/Badge';
import { ArrowLeft, Clock, Calendar, MapPin } from 'lucide-react';

export default function BookingDetails() {
  const { id } = useParams();
  const navigate = useNavigate();
  const booking = mockBookings.find((b) => b.id === Number(id));

  if (!booking) {
    return (
      <div className="flex items-center justify-center h-screen">
        <p className="text-gray-500">Booking not found</p>
      </div>
    );
  }

  const design = mockDesigns.find((d) => d.id === booking.designId);

  return (
    <div className="pb-24">
      <div className="relative">
        <div className="aspect-[16/9] bg-rose-50 dark:bg-neutral-800 flex items-center justify-center">
          <span className="text-7xl">🎨</span>
        </div>
        <button
          onClick={() => navigate(-1)}
          className="absolute top-4 left-4 w-10 h-10 bg-white/90 dark:bg-neutral-800/90 backdrop-blur-sm rounded-full flex items-center justify-center shadow-sm"
        >
          <ArrowLeft size={20} className="text-gray-700 dark:text-white" />
        </button>
      </div>

      <div className="px-5 pt-6">
        <div className="flex items-start justify-between mb-4">
          <div>
            <Badge status={booking.status}>{booking.status}</Badge>
            <h1 className="text-xl font-bold text-gray-900 dark:text-white mt-2">
              {design?.name || booking.designName}
            </h1>
            <p className="text-sm text-gray-500 dark:text-neutral-400 mt-1">
              {design?.category || booking.designCategory}
            </p>
          </div>
          <div className="text-right">
            <p className="text-xl font-bold text-rose-600 dark:text-orange-400">₹{design?.price || 0}</p>
          </div>
        </div>

        <div className="bg-rose-50 dark:bg-neutral-800 rounded-2xl p-5 mb-6 border border-rose-100 dark:border-neutral-700">
          <h3 className="text-sm font-semibold text-gray-900 dark:text-white mb-4">Booking Details</h3>
          <div className="space-y-3">
            <div className="flex items-center gap-3">
              <Calendar size={18} className="text-rose-500 dark:text-orange-400" />
              <div>
                <p className="text-xs text-gray-500 dark:text-neutral-400">Date</p>
                <p className="text-sm font-medium text-gray-900 dark:text-white">{booking.date}</p>
              </div>
            </div>
            <div className="flex items-center gap-3">
              <Clock size={18} className="text-rose-500 dark:text-orange-400" />
              <div>
                <p className="text-xs text-gray-500 dark:text-neutral-400">Time</p>
                <p className="text-sm font-medium text-gray-900 dark:text-white">
                  {booking.startTime} - {booking.endTime}
                </p>
              </div>
            </div>
            <div className="flex items-center gap-3">
              <MapPin size={18} className="text-rose-500 dark:text-orange-400" />
              <div>
                <p className="text-xs text-gray-500 dark:text-neutral-400">Location</p>
                <p className="text-sm font-medium text-gray-900 dark:text-white">Studio Address</p>
              </div>
            </div>
          </div>
        </div>

        {booking.status === 'PENDING' && (
          <div className="flex gap-3">
            <Button variant="outline" className="flex-1" onClick={() => {}}>Cancel</Button>
            <Button className="flex-1">Confirm</Button>
          </div>
        )}
        {booking.status === 'CONFIRMED' && (
          <Button variant="outline" className="w-full" onClick={() => {}}>Cancel Booking</Button>
        )}
      </div>
    </div>
  );
}

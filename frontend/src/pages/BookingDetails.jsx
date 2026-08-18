import { useEffect, useState } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { bookingService } from '../services/bookingService';
import Button from '../components/ui/Button';
import Badge from '../components/ui/Badge';
import BookingStatus from '../components/bookings/BookingStatus';
import { ArrowLeft, Clock, Calendar, MapPin, AlertCircle, XCircle, Loader2 } from 'lucide-react';

function formatTime(time24) {
  if (!time24) return '';
  const [hours, minutes] = time24.split(':').map(Number);
  const period = hours >= 12 ? 'PM' : 'AM';
  const hours12 = hours % 12 || 12;
  return `${hours12}:${minutes.toString().padStart(2, '0')} ${period}`;
}

export default function BookingDetails() {
  const { id } = useParams();
  const navigate = useNavigate();
  const { user } = useAuth();
  const [booking, setBooking] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [cancelling, setCancelling] = useState(false);

  useEffect(() => {
    let cancelled = false;

    async function loadBooking() {
      setLoading(true);
      setError('');

      const result = await bookingService.getMyBooking(id);
      if (cancelled) return;

      if (result.error) {
        if (result.status === 401) {
          navigate('/login');
          return;
        }
        if (result.status === 404) {
          setError('Booking not found');
        } else {
          setError(result.message || 'Failed to load booking');
        }
        setLoading(false);
        return;
      }

      setBooking(result.data);
      setLoading(false);
    }

    loadBooking();

    return () => {
      cancelled = true;
    };
  }, [id, navigate]);

  const handleCancel = async () => {
    if (!booking) return;
    
    const confirmed = window.confirm(
      `Are you sure you want to cancel this booking for ${booking.designName} on ${booking.booking_date} at ${formatTime(booking.start_time)}?`
    );
    
    if (!confirmed) return;

    setCancelling(true);
    setError('');

    const result = await bookingService.cancelBooking(booking.id);
    if (result.error) {
      if (result.status === 409) {
        setError('This booking cannot be cancelled: ' + (result.message || 'Invalid status'));
      } else {
        setError(result.message || 'Failed to cancel booking');
      }
      setCancelling(false);
      return;
    }

    setBooking(result.data);
    setCancelling(false);
  };

  if (loading) {
    return (
      <div className="pb-24">
        <div className="px-5 pt-6">
          <div className="h-8 bg-gray-200 dark:bg-neutral-800 rounded w-48 mb-4 animate-pulse" />
          <div className="h-64 bg-gray-200 dark:bg-neutral-800 rounded-2xl animate-pulse mb-4" />
          <div className="h-48 bg-gray-200 dark:bg-neutral-800 rounded-2xl animate-pulse" />
        </div>
      </div>
    );
  }

  if (error || !booking) {
    return (
      <div className="flex items-center justify-center h-screen">
        <div className="text-center px-5">
          <p className="text-gray-500 dark:text-neutral-400 mb-4">{error || 'Booking not found'}</p>
          <Link to="/bookings" className="text-rose-500 text-sm font-medium">
            Back to bookings
          </Link>
        </div>
      </div>
    );
  }

  const canCancel = booking.status === 'PENDING' || booking.status === 'CONFIRMED';

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
            <BookingStatus status={booking.status} />
            <h1 className="text-xl font-bold text-gray-900 dark:text-white mt-2">
              {booking.designName}
            </h1>
            <p className="text-sm text-gray-500 dark:text-neutral-400 mt-1">
              {booking.design_category || booking.designCategory}
            </p>
          </div>
        </div>

        <div className="bg-white dark:bg-neutral-900 rounded-2xl border border-rose-100 dark:border-neutral-800 p-5 shadow-sm mb-6">
          <div className="space-y-4">
            <div className="flex items-center justify-between py-2 border-b border-rose-100 dark:border-neutral-800">
              <div className="flex items-center gap-3 text-gray-600 dark:text-neutral-400">
                <Calendar size={18} className="text-rose-500" />
                <span className="font-medium text-gray-900 dark:text-white">{booking.booking_date}</span>
              </div>
            </div>
            <div className="flex items-center justify-between py-2 border-b border-rose-100 dark:border-neutral-800">
              <div className="flex items-center gap-3 text-gray-600 dark:text-neutral-400">
                <Clock size={18} className="text-rose-500" />
                <span className="font-medium text-gray-900 dark:text-white">
                  {formatTime(booking.start_time)} - {formatTime(booking.end_time)}
                </span>
              </div>
            </div>
            {booking.notes && (
              <div className="flex items-start justify-between py-2">
                <div className="flex items-start gap-3 text-gray-600 dark:text-neutral-400">
                  <MapPin size={18} className="text-rose-500 mt-0.5" />
                  <div>
                    <p className="text-xs text-gray-500 dark:text-neutral-400">Notes</p>
                    <p className="font-medium text-gray-900 dark:text-white">{booking.notes}</p>
                  </div>
                </div>
              </div>
            )}
            <div className="flex items-center justify-between py-2">
              <div className="flex items-center gap-3 text-gray-600 dark:text-neutral-400">
                <span>🆔</span>
                <span className="font-mono text-xs text-gray-500 dark:text-neutral-400">{booking.id}</span>
              </div>
            </div>
          </div>
        </div>

        {error && (
          <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 px-4 py-3 rounded-xl mb-4 text-sm flex items-start gap-2">
            <AlertCircle size={18} className="shrink-0 mt-0.5" />
            <span>{error}</span>
          </div>
        )}

        {canCancel && (
          <Button
            variant="outline"
            className="w-full"
            onClick={handleCancel}
            disabled={cancelling}
          >
            {cancelling ? (
              <>
                <Loader2 size={18} className="animate-spin mr-2" />
                Cancelling...
              </>
            ) : (
              <>
                <XCircle size={18} className="mr-2" />
                Cancel Booking
              </>
            )}
          </Button>
        )}

        {!canCancel && booking.status === 'CANCELLED' && (
          <div className="w-full text-center text-sm text-gray-500 dark:text-neutral-400 py-2">
            This booking has been cancelled
          </div>
        )}

        {!canCancel && booking.status === 'COMPLETED' && (
          <div className="w-full text-center text-sm text-gray-500 dark:text-neutral-400 py-2">
            This booking has been completed
          </div>
        )}
      </div>
    </div>
  );
}
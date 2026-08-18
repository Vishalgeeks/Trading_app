import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { bookingService } from '../services/bookingService';
import Button from '../components/ui/Button';
import BookingStatus from '../components/bookings/BookingStatus';
import { ArrowLeft, Clock, Calendar, MapPin, User, Mail, Phone, Loader2, AlertCircle } from 'lucide-react';

function formatTime(time24) {
  if (!time24) return '';
  const [hours, minutes] = time24.split(':').map(Number);
  const period = hours >= 12 ? 'PM' : 'AM';
  const hours12 = hours % 12 || 12;
  return `${hours12}:${minutes.toString().padStart(2, '0')} ${period}`;
}

const VALID_TRANSITIONS = {
  PENDING: ['CONFIRMED', 'CANCELLED'],
  CONFIRMED: ['COMPLETED', 'CANCELLED'],
  CANCELLED: [],
  COMPLETED: []
};

export default function AdminBookingDetails() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [booking, setBooking] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [updating, setUpdating] = useState(false);

  useEffect(() => {
    let cancelled = false;

    async function loadBooking() {
      setLoading(true);
      setError('');

      const result = await bookingService.getAdminBooking(id);
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

  const handleStatusUpdate = async (newStatus) => {
    if (!booking || updating) return;

    setUpdating(true);
    setError('');

    const result = await bookingService.updateBookingStatus(booking.id, newStatus);
    if (result.error) {
      if (result.status === 409) {
        setError('Invalid status transition: ' + (result.message || 'Not allowed'));
      } else {
        setError(result.message || 'Failed to update status');
      }
      setUpdating(false);
      return;
    }

    setBooking(result.data);
    setUpdating(false);
  };

  const availableActions = booking ? VALID_TRANSITIONS[booking.status] || [] : [];

  if (loading) {
    return (
      <div className="p-5">
        <div className="h-8 bg-gray-200 dark:bg-neutral-800 rounded w-48 mb-4 animate-pulse" />
        <div className="h-64 bg-gray-200 dark:bg-neutral-800 rounded-2xl animate-pulse mb-4" />
        <div className="h-48 bg-gray-200 dark:bg-neutral-800 rounded-2xl animate-pulse" />
      </div>
    );
  }

  if (error || !booking) {
    return (
      <div className="flex items-center justify-center h-screen">
        <div className="text-center px-5">
          <p className="text-gray-500 dark:text-neutral-400 mb-4">{error || 'Booking not found'}</p>
          <button onClick={() => navigate('/admin/bookings')} className="text-rose-500 text-sm font-medium">
            Back to bookings
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="p-5">
      <div className="flex items-center gap-4 mb-6">
        <button
          onClick={() => navigate('/admin/bookings')}
          className="w-10 h-10 bg-white dark:bg-neutral-800 border border-rose-100 dark:border-neutral-700 rounded-xl flex items-center justify-center shadow-sm hover:bg-gray-50 dark:hover:bg-neutral-700 transition-colors"
        >
          <ArrowLeft size={20} className="text-gray-700 dark:text-white" />
        </button>
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Booking Details</h1>
          <p className="text-sm text-gray-500 dark:text-neutral-400">Manage booking status and details</p>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Main Details */}
        <div className="lg:col-span-2 space-y-6">
          <div className="bg-white dark:bg-neutral-900 rounded-2xl border border-rose-100 dark:border-neutral-800 p-5 shadow-sm">
            <div className="flex items-start justify-between mb-4">
              <div>
                <BookingStatus status={booking.status} />
                <h2 className="text-xl font-bold text-gray-900 dark:text-white mt-2">{booking.designName}</h2>
              </div>
            </div>

            <div className="space-y-4">
              <div className="flex items-center gap-3 py-3 border-b border-rose-100 dark:border-neutral-800">
                <div className="w-10 h-10 bg-rose-50 dark:bg-orange-950 rounded-xl flex items-center justify-center text-rose-500 dark:text-orange-400">
                  <Calendar size={20} />
                </div>
                <div>
                  <p className="text-xs text-gray-500 dark:text-neutral-400">Date</p>
                  <p className="font-medium text-gray-900 dark:text-white">{booking.booking_date}</p>
                </div>
              </div>

              <div className="flex items-center gap-3 py-3 border-b border-rose-100 dark:border-neutral-800">
                <div className="w-10 h-10 bg-rose-50 dark:bg-orange-950 rounded-xl flex items-center justify-center text-rose-500 dark:text-orange-400">
                  <Clock size={20} />
                </div>
                <div>
                  <p className="text-xs text-gray-500 dark:text-neutral-400">Time</p>
                  <p className="font-medium text-gray-900 dark:text-white">
                    {formatTime(booking.start_time)} - {formatTime(booking.end_time)}
                  </p>
                </div>
              </div>

              <div className="flex items-center gap-3 py-3 border-b border-rose-100 dark:border-neutral-800">
                <div className="w-10 h-10 bg-rose-50 dark:bg-orange-950 rounded-xl flex items-center justify-center text-rose-500 dark:text-orange-400">
                  <span className="text-xl">🎨</span>
                </div>
                <div>
                  <p className="text-xs text-gray-500 dark:text-neutral-400">Design</p>
                  <p className="font-medium text-gray-900 dark:text-white">{booking.designName}</p>
                </div>
              </div>

              {booking.notes && (
                <div className="flex items-start gap-3 py-3 border-b border-rose-100 dark:border-neutral-800">
                  <div className="w-10 h-10 bg-rose-50 dark:bg-orange-950 rounded-xl flex items-center justify-center text-rose-500 dark:text-orange-400 mt-0.5">
                    <MapPin size={20} />
                  </div>
                  <div>
                    <p className="text-xs text-gray-500 dark:text-neutral-400">Notes</p>
                    <p className="font-medium text-gray-900 dark:text-white">{booking.notes}</p>
                  </div>
                </div>
              )}

              <div className="flex items-center gap-3 py-3">
                <div className="w-10 h-10 bg-rose-50 dark:bg-orange-950 rounded-xl flex items-center justify-center text-rose-500 dark:text-orange-400">
                  <span className="text-xl">🆔</span>
                </div>
                <div>
                  <p className="text-xs text-gray-500 dark:text-neutral-400">Booking ID</p>
                  <p className="font-mono text-xs text-gray-500 dark:text-neutral-400">{booking.id}</p>
                </div>
              </div>
            </div>
          </div>

          {/* Client Info */}
          <div className="bg-white dark:bg-neutral-900 rounded-2xl border border-rose-100 dark:border-neutral-800 p-5 shadow-sm">
            <h3 className="font-semibold text-gray-900 dark:text-white mb-4">Client Information</h3>
            <div className="space-y-3">
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 bg-rose-50 dark:bg-orange-950 rounded-xl flex items-center justify-center text-rose-500 dark:text-orange-400">
                  <User size={20} />
                </div>
                <div>
                  <p className="text-xs text-gray-500 dark:text-neutral-400">Name</p>
                  <p className="font-medium text-gray-900 dark:text-white">{booking.clientName}</p>
                </div>
              </div>

              <div className="flex items-center gap-3">
                <div className="w-10 h-10 bg-rose-50 dark:bg-orange-950 rounded-xl flex items-center justify-center text-rose-500 dark:text-orange-400">
                  <Mail size={20} />
                </div>
                <div>
                  <p className="text-xs text-gray-500 dark:text-neutral-400">Email</p>
                  <p className="font-medium text-gray-900 dark:text-white">{booking.clientEmail}</p>
                </div>
              </div>

              {booking.clientPhone && (
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 bg-rose-50 dark:bg-orange-950 rounded-xl flex items-center justify-center text-rose-500 dark:text-orange-400">
                    <Phone size={20} />
                  </div>
                  <div>
                    <p className="text-xs text-gray-500 dark:text-neutral-400">Phone</p>
                    <p className="font-medium text-gray-900 dark:text-white">{booking.clientPhone}</p>
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>

        {/* Actions Sidebar */}
        <div className="lg:col-span-1">
          <div className="bg-white dark:bg-neutral-900 rounded-2xl border border-rose-100 dark:border-neutral-800 p-5 shadow-sm sticky top-24">
            <h3 className="font-semibold text-gray-900 dark:text-white mb-4">Actions</h3>

            <div className="space-y-3">
              {availableActions.map((status) => (
                <Button
                  key={status}
                  variant={status === 'CANCELLED' ? 'outline' : 'default'}
                  className="w-full"
                  onClick={() => handleStatusUpdate(status)}
                  disabled={updating}
                >
                  {updating ? (
                    <>
                      <Loader2 size={18} className="animate-spin mr-2" />
                      Updating...
                    </>
                  ) : (
                    status === 'CONFIRMED' ? 'Confirm Booking' :
                    status === 'COMPLETED' ? 'Mark Complete' :
                    status === 'CANCELLED' ? 'Cancel Booking' : status
                  )}
                </Button>
              ))}

              {availableActions.length === 0 && (
                <div className="text-center py-4 text-sm text-gray-500 dark:text-neutral-400">
                  No further actions available for {booking.status} bookings.
                </div>
              )}

              <div className="pt-4 border-t border-rose-100 dark:border-neutral-800">
                <p className="text-xs text-gray-500 dark:text-neutral-400 text-center">
                  Created: {new Date(booking.createdAt).toLocaleDateString()}
                  <br />
                  Updated: {new Date(booking.updatedAt).toLocaleDateString()}
                </p>
              </div>
            </div>

            {error && (
              <div className="mt-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 px-4 py-3 rounded-xl text-sm flex items-start gap-2">
                <AlertCircle size={18} className="shrink-0 mt-0.5" />
                <span>{error}</span>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
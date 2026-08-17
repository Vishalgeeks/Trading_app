import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { designService } from '../services/designService';
import { availabilityService } from '../services/availabilityService';
import { bookingService } from '../services/bookingService';
import Button from '../components/ui/Button';
import { ArrowLeft, Clock, AlertCircle, CheckCircle } from 'lucide-react';

function formatDuration(minutes) {
  if (!minutes || minutes <= 0) return '';
  const hours = Math.floor(minutes / 60);
  const mins = minutes % 60;
  if (hours > 0 && mins > 0) return `${hours}h ${mins}m`;
  if (hours > 0) return `${hours}h`;
  return `${mins}m`;
}

function formatTimeTo12Hour(time24) {
  if (!time24) return '';
  const [hours, minutes] = time24.split(':').map(Number);
  const period = hours >= 12 ? 'PM' : 'AM';
  const hours12 = hours % 12 || 12;
  return `${hours12}:${minutes.toString().padStart(2, '0')} ${period}`;
}

function getTodayDate() {
  const today = new Date();
  const year = today.getFullYear();
  const month = (today.getMonth() + 1).toString().padStart(2, '0');
  const day = today.getDate().toString().padStart(2, '0');
  return `${year}-${month}-${day}`;
}

export default function Booking() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [design, setDesign] = useState(null);
  const [loading, setLoading] = useState(true);
  const [slotsLoading, setSlotsLoading] = useState(false);
  const [error, setError] = useState('');
  const [selectedDate, setSelectedDate] = useState(getTodayDate);
  const [selectedSlot, setSelectedSlot] = useState(null);
  const [slots, setSlots] = useState([]);
  const [notes, setNotes] = useState('');
  const [bookingLoading, setBookingLoading] = useState(false);
  const [bookingResult, setBookingResult] = useState(null);

  useEffect(() => {
    let cancelled = false;

    async function loadDesign() {
      setLoading(true);
      setError('');

      const result = await designService.getDesign(id);
      if (cancelled) return;

      if (result.error) {
        setError(result.message || 'Failed to load design');
        setLoading(false);
        return;
      }

      setDesign(result.data);
      setLoading(false);
    }

    loadDesign();

    return () => {
      cancelled = true;
    };
  }, [id]);

  useEffect(() => {
    if (!selectedDate || !design?.id) return;

    let cancelled = false;
    setSlotsLoading(true);
    setSelectedSlot(null);

    async function loadSlots() {
      const result = await availabilityService.getAvailableSlots(selectedDate, design.id);
      if (cancelled) return;

      if (!result.error) {
        setSlots(result.data || []);
      } else {
        setSlots([]);
      }
      setSlotsLoading(false);
    }

    loadSlots();

    return () => {
      cancelled = true;
    };
  }, [selectedDate, design?.id]);

  const handleDateChange = (e) => {
    setSelectedDate(e.target.value);
    setSelectedSlot(null);
    setBookingResult(null);
  };

  const handleSlotSelect = (slot) => {
    setSelectedSlot(slot);
    setBookingResult(null);
  };

  const handleConfirmBooking = async () => {
    if (!selectedSlot || !design || !selectedDate) return;

    setBookingLoading(true);
    setBookingResult(null);

    const result = await bookingService.createBooking({
      designId: design.id,
      bookingDate: selectedDate,
      startTime: selectedSlot.start_time,
      notes: notes.trim(),
    });

    setBookingLoading(false);

    if (result.error) {
      if (result.status === 409) {
        setBookingResult({
          type: 'conflict',
          message: 'This time slot was just booked by someone else. Please choose another time.',
        });
        setSelectedSlot(null);
        window.scrollTo({ top: 0, behavior: 'smooth' });
      } else if (result.status === 401) {
        navigate('/login');
      } else {
        setBookingResult({
          type: 'error',
          message: result.message || 'Failed to create booking. Please try again.',
        });
      }
      return;
    }

    setBookingResult({
      type: 'success',
      data: result.data,
    });
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

  if (error || !design) {
    return (
      <div className="flex items-center justify-center h-screen">
        <div className="text-center px-5">
          <p className="text-gray-500 dark:text-neutral-400 mb-4">{error || 'Design not found'}</p>
          <button onClick={() => navigate('/browse')} className="text-rose-500 text-sm font-medium">
            Browse designs
          </button>
        </div>
      </div>
    );
  }

  if (bookingResult?.type === 'success') {
    const booking = bookingResult.data;
    return (
      <div className="pb-24">
        <div className="px-5 pt-6">
          <div className="text-center py-8">
            <div className="w-16 h-16 bg-green-100 dark:bg-green-900/30 rounded-full flex items-center justify-center mx-auto mb-4">
              <CheckCircle size={32} className="text-green-600 dark:text-green-400" />
            </div>
            <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-2">Booking Confirmed</h2>
            <p className="text-sm text-gray-500 dark:text-neutral-400 mb-6">
              Your booking has been created successfully.
            </p>
          </div>

          <div className="bg-white dark:bg-neutral-900 rounded-2xl border border-rose-100 dark:border-neutral-800 p-5 shadow-sm mb-6">
            <div className="space-y-3">
              <div className="flex items-center justify-between">
                <span className="text-sm text-gray-500 dark:text-neutral-400">Booking ID</span>
                <span className="text-sm font-medium text-gray-900 dark:text-white">{booking.id}</span>
              </div>
              <div className="border-t border-rose-100 dark:border-neutral-800" />
              <div className="flex items-center justify-between">
                <span className="text-sm text-gray-500 dark:text-neutral-400">Design</span>
                <span className="text-sm font-medium text-gray-900 dark:text-white">{booking.design_name}</span>
              </div>
              <div className="border-t border-rose-100 dark:border-neutral-800" />
              <div className="flex items-center justify-between">
                <span className="text-sm text-gray-500 dark:text-neutral-400">Date</span>
                <span className="text-sm font-medium text-gray-900 dark:text-white">{booking.booking_date}</span>
              </div>
              <div className="border-t border-rose-100 dark:border-neutral-800" />
              <div className="flex items-center justify-between">
                <span className="text-sm text-gray-500 dark:text-neutral-400">Time</span>
                <span className="text-sm font-medium text-gray-900 dark:text-white">
                  {formatTimeTo12Hour(booking.start_time)} - {formatTimeTo12Hour(booking.end_time)}
                </span>
              </div>
              <div className="border-t border-rose-100 dark:border-neutral-800" />
              <div className="flex items-center justify-between">
                <span className="text-sm text-gray-500 dark:text-neutral-400">Status</span>
                <span className="inline-block px-3 py-1 bg-amber-100 dark:bg-amber-900/30 text-amber-700 dark:text-amber-400 rounded-full text-xs font-semibold">
                  {booking.status}
                </span>
              </div>
            </div>
          </div>

          <div className="flex gap-3">
            <Button variant="outline" className="flex-1" onClick={() => navigate('/bookings')}>
              View Bookings
            </Button>
            <Button className="flex-1" onClick={() => navigate('/browse')}>
              Browse More
            </Button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="pb-24">
      <div className="px-5 pt-6">
        <button
          onClick={() => navigate(-1)}
          className="flex items-center gap-2 text-sm text-gray-600 dark:text-neutral-400 mb-4 hover:text-rose-500"
        >
          <ArrowLeft size={18} />
          Back
        </button>

        <h1 className="text-2xl font-bold text-gray-900 dark:text-white mb-4">Book Appointment</h1>

        <div className="bg-white dark:bg-neutral-900 rounded-2xl border border-rose-100 dark:border-neutral-800 p-4 shadow-sm mb-6">
          <div className="flex items-center gap-3">
            <div className="w-16 h-16 rounded-xl bg-rose-50 dark:bg-neutral-800 flex items-center justify-center text-3xl shrink-0">
              ??
            </div>
            <div className="flex-1 min-w-0">
              <h3 className="font-semibold text-gray-900 dark:text-white text-sm truncate">{design.name}</h3>
              <p className="text-xs text-gray-500 dark:text-neutral-400">{design.category?.name || 'Uncategorized'}</p>
              <div className="flex items-center gap-3 mt-1">
                <span className="text-xs text-gray-500 dark:text-neutral-400 flex items-center gap-1">
                  <Clock size={12} />
                  {formatDuration(design.duration_minutes)}
                </span>
                <span className="text-sm font-bold text-rose-600 dark:text-orange-400">?{design.price}</span>
              </div>
            </div>
          </div>
        </div>

        {bookingResult?.type === 'conflict' && (
          <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 px-4 py-3 rounded-xl mb-4 text-sm flex items-start gap-2">
            <AlertCircle size={18} className="shrink-0 mt-0.5" />
            <span>{bookingResult.message}</span>
          </div>
        )}

        {bookingResult?.type === 'error' && (
          <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 px-4 py-3 rounded-xl mb-4 text-sm flex items-start gap-2">
            <AlertCircle size={18} className="shrink-0 mt-0.5" />
            <span>{bookingResult.message}</span>
          </div>
        )}

        <div className="mb-6">
          <label className="block text-sm font-medium text-gray-700 dark:text-neutral-300 mb-1.5">
            Select Date
          </label>
          <input
            type="date"
            value={selectedDate}
            min={getTodayDate()}
            onChange={handleDateChange}
            className="w-full px-4 py-3 bg-white dark:bg-neutral-800 border border-rose-200 dark:border-neutral-700 rounded-xl text-gray-900 dark:text-white focus:ring-2 focus:ring-rose-500 focus:border-transparent outline-none"
          />
        </div>

        <div className="mb-6">
          <label className="block text-sm font-medium text-gray-700 dark:text-neutral-300 mb-3">
            Available Time Slots
          </label>
          {slotsLoading ? (
            <div className="grid grid-cols-3 gap-2">
              {[1, 2, 3, 4, 5, 6].map((i) => (
                <div key={i} className="h-12 bg-gray-200 dark:bg-neutral-800 rounded-xl animate-pulse" />
              ))}
            </div>
          ) : slots.length > 0 ? (
            <div className="grid grid-cols-3 gap-2">
              {slots.map((slot) => {
                const isSelected = selectedSlot?.start_time === slot.start_time;
                return (
                  <button
                    key={slot.start_time}
                    onClick={() => handleSlotSelect(slot)}
                    disabled={bookingLoading}
                    className={`py-3 px-2 rounded-xl text-sm font-medium transition-all duration-200 border-2 ${
                      isSelected
                        ? 'bg-rose-50 border-rose-500 text-rose-700 dark:bg-orange-950 dark:border-orange-500 dark:text-orange-400'
                        : 'bg-white dark:bg-neutral-800 border-gray-200 dark:border-neutral-700 text-gray-700 dark:text-neutral-300 hover:border-rose-300'
                    } ${bookingLoading ? 'opacity-50 cursor-not-allowed' : ''}`}
                  >
                    {formatTimeTo12Hour(slot.start_time)}
                  </button>
                );
              })}
            </div>
          ) : (
            <div className="text-center py-8 bg-rose-50 dark:bg-neutral-800 rounded-2xl">
              <p className="text-gray-500 dark:text-neutral-400 text-sm mb-2">No available slots for this date.</p>
              <button
                onClick={() => setSelectedDate(getTodayDate())}
                className="text-sm text-rose-500 dark:text-orange-400 font-medium"
              >
                Choose another date
              </button>
            </div>
          )}
        </div>

        {selectedSlot && (
          <div className="bg-rose-50 dark:bg-neutral-800 rounded-2xl p-5 border border-rose-100 dark:border-neutral-700 mb-6">
            <h3 className="text-sm font-bold text-gray-900 dark:text-white mb-4">Booking Summary</h3>
            <div className="space-y-3">
              <div className="flex items-center justify-between">
                <span className="text-sm text-gray-500 dark:text-neutral-400">Design</span>
                <span className="text-sm font-medium text-gray-900 dark:text-white">{design.name}</span>
              </div>
              <div className="border-t border-rose-100 dark:border-neutral-700" />
              <div className="flex items-center justify-between">
                <span className="text-sm text-gray-500 dark:text-neutral-400">Date</span>
                <span className="text-sm font-medium text-gray-900 dark:text-white">{selectedDate}</span>
              </div>
              <div className="border-t border-rose-100 dark:border-neutral-700" />
              <div className="flex items-center justify-between">
                <span className="text-sm text-gray-500 dark:text-neutral-400">Time</span>
                <span className="text-sm font-medium text-gray-900 dark:text-white">
                  {formatTimeTo12Hour(selectedSlot.start_time)} - {formatTimeTo12Hour(selectedSlot.end_time)}
                </span>
              </div>
              <div className="border-t border-rose-100 dark:border-neutral-700" />
              <div className="flex items-center justify-between">
                <span className="text-sm text-gray-500 dark:text-neutral-400">Duration</span>
                <span className="text-sm font-medium text-gray-900 dark:text-white">{formatDuration(design.duration_minutes)}</span>
              </div>
              <div className="border-t border-rose-100 dark:border-neutral-700" />
              <div className="flex items-center justify-between">
                <span className="text-sm text-gray-500 dark:text-neutral-400">Price</span>
                <span className="text-sm font-bold text-rose-600 dark:text-orange-400">?{design.price}</span>
              </div>
            </div>
          </div>
        )}

        <div className="mb-6">
          <label className="block text-sm font-medium text-gray-700 dark:text-neutral-300 mb-1.5">
            Notes (optional)
          </label>
          <textarea
            value={notes}
            onChange={(e) => setNotes(e.target.value)}
            placeholder="Any special requests..."
            rows={3}
            className="w-full px-4 py-3 bg-white dark:bg-neutral-800 border border-rose-200 dark:border-neutral-700 rounded-xl text-gray-900 dark:text-white placeholder-gray-400 focus:ring-2 focus:ring-rose-500 focus:border-transparent outline-none resize-none"
          />
        </div>

        <Button
          className="w-full"
          disabled={!selectedSlot || bookingLoading}
          onClick={handleConfirmBooking}
        >
          {bookingLoading ? 'Confirming...' : 'Confirm Booking'}
        </Button>
      </div>
    </div>
  );
}

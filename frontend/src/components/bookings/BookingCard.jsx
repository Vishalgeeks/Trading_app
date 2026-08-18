import BookingStatus from './BookingStatus';

export default function BookingCard({ booking, onPress }) {
  return (
    <div
      onClick={onPress}
      className="bg-white dark:bg-neutral-900 rounded-2xl border border-rose-100 dark:border-neutral-800 p-4 shadow-sm active:scale-[0.98] transition-transform cursor-pointer"
    >
      <div className="flex items-start justify-between mb-3">
        <div>
          <h3 className="font-semibold text-gray-900 dark:text-white text-sm">
            {booking.designName || booking.design_name}
          </h3>
          <p className="text-xs text-gray-500 dark:text-neutral-400 mt-0.5">
            {booking.designCategory || booking.design_category}
          </p>
        </div>
        <BookingStatus status={booking.status} />
      </div>
      <div className="flex items-center gap-4 text-xs text-gray-600 dark:text-neutral-400">
        <div className="flex items-center gap-1">
          <span>📅</span>
          <span>{booking.booking_date || booking.date}</span>
        </div>
        <div className="flex items-center gap-1">
          <span>🕐</span>
          <span>{booking.start_time || booking.startTime} - {booking.end_time || booking.endTime}</span>
        </div>
      </div>
    </div>
  );
}
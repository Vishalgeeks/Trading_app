import Badge from '../ui/Badge';

export default function AppointmentCard({ appointment }) {
  return (
    <div className="bg-white dark:bg-neutral-900 rounded-2xl border border-rose-100 dark:border-neutral-800 p-4 shadow-sm">
      <div className="flex items-start justify-between mb-2">
        <div>
          <h3 className="font-semibold text-gray-900 dark:text-white text-sm">
            {appointment.clientName || appointment.client_name}
          </h3>
          <p className="text-xs text-gray-500 dark:text-neutral-400">
            {appointment.designName || appointment.design_name}
          </p>
        </div>
        <Badge status={appointment.status}>
          {appointment.status}
        </Badge>
      </div>
      <div className="flex items-center gap-3 text-xs text-gray-600 dark:text-neutral-400 mt-3">
        <div className="flex items-center gap-1">
          <span>📅</span>
          <span>{appointment.date}</span>
        </div>
        <div className="flex items-center gap-1">
          <span>🕐</span>
          <span>{appointment.startTime || appointment.start_time} - {appointment.endTime || appointment.end_time}</span>
        </div>
      </div>
    </div>
  );
}

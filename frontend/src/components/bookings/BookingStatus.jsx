import { Clock, CheckCircle, XCircle, AlertCircle, Calendar } from 'lucide-react';

export default function BookingStatus({ status }) {
  const statusConfig = {
    PENDING: {
      label: 'Pending',
      description: 'Awaiting confirmation',
      color: 'bg-amber-100 dark:bg-amber-900/30 text-amber-700 dark:text-amber-400 border-amber-200 dark:border-amber-800',
      icon: AlertCircle,
      iconColor: 'text-amber-500'
    },
    CONFIRMED: {
      label: 'Confirmed',
      description: 'Booking accepted',
      color: 'bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-400 border-green-200 dark:border-green-800',
      icon: CheckCircle,
      iconColor: 'text-green-500'
    },
    CANCELLED: {
      label: 'Cancelled',
      description: 'Booking cancelled',
      color: 'bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-400 border-red-200 dark:border-red-800',
      icon: XCircle,
      iconColor: 'text-red-500'
    },
    COMPLETED: {
      label: 'Completed',
      description: 'Appointment completed',
      color: 'bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-400 border-blue-200 dark:border-blue-800',
      icon: Calendar,
      iconColor: 'text-blue-500'
    }
  };

  const config = statusConfig[status] || statusConfig.PENDING;
  const Icon = config.icon;

  return (
    <div className={`inline-flex items-center gap-2 px-3 py-1.5 rounded-xl border text-sm font-medium ${config.color}`}>
      <Icon size={14} className={config.iconColor} />
      <span>{config.label}</span>
    </div>
  );
}

export function getStatusConfig(status) {
  return statusConfig[status] || statusConfig.PENDING;
}
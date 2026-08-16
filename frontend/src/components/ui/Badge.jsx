const statusStyles = {
  CONFIRMED: 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400',
  PENDING: 'bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-400',
  CANCELLED: 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400',
  COMPLETED: 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400',
};

export default function Badge({ children, status, className = '' }) {
  return (
    <span
      className={`inline-flex items-center px-3 py-1 rounded-full text-xs font-semibold ${statusStyles[status] || 'bg-gray-100 text-gray-700'} ${className}`}
    >
      {children}
    </span>
  );
}

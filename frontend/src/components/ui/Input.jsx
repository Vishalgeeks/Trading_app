export default function Input({ label, className = '', ...props }) {
  return (
    <div className={className}>
      {label && (
        <label className="block text-sm font-medium text-gray-700 mb-1.5 dark:text-neutral-300">
          {label}
        </label>
      )}
      <input
        className="w-full px-4 py-3 bg-white border border-rose-200 rounded-xl focus:ring-2 focus:ring-rose-500 focus:border-transparent outline-none transition-all duration-200 text-gray-900 dark:bg-neutral-800 dark:border-neutral-700 dark:text-white dark:placeholder-neutral-500"
        {...props}
      />
    </div>
  );
}

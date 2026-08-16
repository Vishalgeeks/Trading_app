export default function Card({ children, className = '', ...props }) {
  return (
    <div
      className={`bg-white rounded-2xl border border-rose-100 shadow-sm p-4 dark:bg-neutral-900 dark:border-neutral-800 ${className}`}
      {...props}
    >
      {children}
    </div>
  );
}

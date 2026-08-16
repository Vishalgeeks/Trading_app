export default function IconButton({ children, variant = 'ghost', className = '', ...props }) {
  const variants = {
    ghost: 'text-gray-500 hover:text-rose-500 hover:bg-rose-50',
    solid: 'bg-rose-500 text-white hover:bg-rose-600',
    outline: 'border border-rose-300 text-rose-500 hover:bg-rose-50',
  };

  const darkVariants = {
    ghost: 'dark:text-neutral-400 dark:hover:text-orange-400 dark:hover:bg-orange-950',
    solid: 'dark:bg-orange-500 dark:hover:bg-orange-600',
    outline: 'dark:border-orange-700 dark:text-orange-400 dark:hover:bg-orange-950',
  };

  return (
    <button
      className={`inline-flex items-center justify-center rounded-full p-2 transition-colors duration-200 ${variants[variant]} ${darkVariants[variant]} ${className}`}
      {...props}
    >
      {children}
    </button>
  );
}

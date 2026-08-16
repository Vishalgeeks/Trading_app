export default function Button({
  children,
  variant = 'primary',
  size = 'md',
  className = '',
  ...props
}) {
  const base = 'inline-flex items-center justify-center gap-2 font-medium rounded-xl transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-offset-2';

  const variants = {
    primary: 'bg-rose-500 hover:bg-rose-600 text-white focus:ring-rose-500',
    secondary: 'bg-amber-400 hover:bg-amber-500 text-gray-900 focus:ring-amber-400',
    outline: 'border-2 border-rose-500 text-rose-500 hover:bg-rose-50 focus:ring-rose-500',
    ghost: 'text-gray-600 hover:text-rose-500 hover:bg-rose-50',
    danger: 'bg-red-500 hover:bg-red-600 text-white focus:ring-red-500',
  };

  const darkVariants = {
    primary: 'dark:bg-orange-500 dark:hover:bg-orange-600',
    secondary: 'dark:bg-amber-400 dark:hover:bg-amber-500',
    outline: 'dark:border-orange-500 dark:text-orange-500 dark:hover:bg-orange-950',
    ghost: 'dark:text-neutral-400 dark:hover:text-orange-400 dark:hover:bg-orange-950',
    danger: 'dark:bg-red-600 dark:hover:bg-red-700',
  };

  const sizes = {
    sm: 'py-2 px-4 text-sm',
    md: 'py-3 px-6 text-base',
    lg: 'py-4 px-8 text-lg',
  };

  return (
    <button
      className={`${base} ${variants[variant]} ${darkVariants[variant]} ${sizes[size]} ${className}`}
      {...props}
    >
      {children}
    </button>
  );
}

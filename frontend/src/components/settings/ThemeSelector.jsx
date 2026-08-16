export default function ThemeSelector({ theme, onSelect }) {
  return (
    <div className="grid grid-cols-2 gap-3">
      <button
        onClick={() => onSelect('light')}
        className={`flex items-center gap-3 p-4 rounded-xl border-2 transition-all duration-200 ${
          theme === 'light'
            ? 'border-rose-500 bg-rose-50 dark:border-orange-500 dark:bg-orange-950'
            : 'border-gray-200 hover:border-rose-300 dark:border-neutral-700 dark:hover:border-orange-700'
        }`}
      >
        <div className="w-8 h-8 rounded-full bg-white border border-gray-200 flex items-center justify-center shadow-sm">
          <span className="text-sm">☀️</span>
        </div>
        <span className={`text-sm font-medium ${theme === 'light' ? 'text-rose-600 dark:text-orange-400' : 'text-gray-700 dark:text-neutral-300'}`}>
          Light
        </span>
      </button>
      <button
        onClick={() => onSelect('dark')}
        className={`flex items-center gap-3 p-4 rounded-xl border-2 transition-all duration-200 ${
          theme === 'dark'
            ? 'border-rose-500 bg-rose-50 dark:border-orange-500 dark:bg-orange-950'
            : 'border-gray-200 hover:border-rose-300 dark:border-neutral-700 dark:hover:border-orange-700'
        }`}
      >
        <div className="w-8 h-8 rounded-full bg-gray-900 border border-gray-700 flex items-center justify-center shadow-sm">
          <span className="text-sm">🌙</span>
        </div>
        <span className={`text-sm font-medium ${theme === 'dark' ? 'text-rose-600 dark:text-orange-400' : 'text-gray-700 dark:text-neutral-300'}`}>
          Dark
        </span>
      </button>
    </div>
  );
}

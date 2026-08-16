import { ChevronRight } from 'lucide-react';

export default function SettingsRow({ icon: Icon, label, description, onClick, trailing }) {
  return (
    <button
      onClick={onClick}
      className="w-full flex items-center gap-4 p-4 rounded-xl hover:bg-rose-50 dark:hover:bg-orange-950/50 transition-colors text-left"
    >
      {Icon && (
        <div className="w-10 h-10 rounded-xl bg-rose-50 dark:bg-orange-950 flex items-center justify-center text-rose-500 dark:text-orange-400 shrink-0">
          <Icon size={20} />
        </div>
      )}
      <div className="flex-1 min-w-0">
        <p className="text-sm font-medium text-gray-900 dark:text-white">{label}</p>
        {description && (
          <p className="text-xs text-gray-500 dark:text-neutral-400 mt-0.5">{description}</p>
        )}
      </div>
      {trailing || <ChevronRight size={18} className="text-gray-400 dark:text-neutral-500 shrink-0" />}
    </button>
  );
}

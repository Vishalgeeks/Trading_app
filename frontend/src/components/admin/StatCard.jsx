import { BarChart3 } from 'lucide-react';

const colorMap = {
  blue: 'bg-blue-500',
  green: 'bg-green-500',
  yellow: 'bg-amber-500',
  purple: 'bg-purple-500',
  red: 'bg-red-500',
  orange: 'bg-orange-500',
};

export default function StatCard({ label, value, color = 'orange', icon: Icon }) {
  return (
    <div className="bg-white dark:bg-neutral-900 rounded-2xl border border-rose-100 dark:border-neutral-800 p-4 shadow-sm">
      <div className={`${colorMap[color] || colorMap.orange} rounded-xl p-3 mb-3 w-fit`}>
        {Icon ? <Icon size={20} className="text-white" /> : <BarChart3 size={20} className="text-white" />}
      </div>
      <p className="text-2xl font-bold text-gray-900 dark:text-white">{value}</p>
      <p className="text-xs text-gray-500 dark:text-neutral-400 mt-1">{label}</p>
    </div>
  );
}

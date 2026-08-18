import { Heart } from 'lucide-react';
import IconButton from '../ui/IconButton';

export default function DesignCard({ design, onFavoriteToggle, onPress }) {
  const isFavorite = Boolean(design.favorite);

  return (
    <div
      onClick={onPress}
      className="bg-white dark:bg-neutral-900 rounded-2xl border border-rose-100 dark:border-neutral-800 overflow-hidden shadow-sm active:scale-[0.98] transition-transform cursor-pointer"
    >
      <div className="relative aspect-[4/3] bg-rose-50 dark:bg-neutral-800 flex items-center justify-center">
        <span className="text-5xl">🎨</span>
        {onFavoriteToggle && (
          <button
            onClick={(e) => {
              e.stopPropagation();
              onFavoriteToggle(design.id);
            }}
            className="absolute top-3 right-3"
          >
            <IconButton variant="ghost">
              <Heart
                size={18}
                fill={isFavorite ? '#f43f5e' : 'none'}
                className={isFavorite ? 'text-rose-500' : 'text-gray-400'}
              />
            </IconButton>
          </button>
        )}
        <div className="absolute bottom-3 left-3">
          <span className="inline-block px-2.5 py-1 bg-white/90 dark:bg-neutral-800/90 backdrop-blur-sm rounded-full text-xs font-medium text-gray-700 dark:text-neutral-300">
            {design.category_name || design.category?.name || 'Uncategorized'}
          </span>
        </div>
      </div>
      <div className="p-3">
        <h3 className="font-semibold text-gray-900 dark:text-white text-sm mb-1 truncate">
          {design.name}
        </h3>
        <div className="flex items-center justify-between">
          <span className="text-xs text-gray-500 dark:text-neutral-400">{design.duration || ''}</span>
          <span className="text-sm font-bold text-rose-600 dark:text-orange-400">₹{design.price || ''}</span>
        </div>
      </div>
    </div>
  );
}
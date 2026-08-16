import { useState } from 'react';
import { mockDesigns } from '../mockData';
import DesignCard from '../components/designs/DesignCard';
import { Heart } from 'lucide-react';

export default function Favorites() {
  const [favorites, setFavorites] = useState(
    mockDesigns.filter((d) => d.favorite)
  );

  const toggleFavorite = (id) => {
    setFavorites((prev) => {
      const exists = prev.find((d) => d.id === id);
      if (exists) {
        return prev.filter((d) => d.id !== id);
      }
      const design = mockDesigns.find((d) => d.id === id);
      if (design) return [...prev, { ...design, favorite: true }];
      return prev;
    });
  };

  return (
    <div className="pb-24">
      <div className="px-5 pt-6">
        <h1 className="text-2xl font-bold text-gray-900 dark:text-white mb-1">Favorites</h1>
        <p className="text-sm text-gray-500 dark:text-neutral-400 mb-6">
          {favorites.length} saved design{favorites.length !== 1 ? 's' : ''}
        </p>

        {favorites.length > 0 ? (
          <div className="grid grid-cols-2 gap-3">
            {favorites.map((design) => (
              <DesignCard
                key={design.id}
                design={design}
                onFavoriteToggle={toggleFavorite}
              />
            ))}
          </div>
        ) : (
          <div className="text-center py-16">
            <div className="w-16 h-16 bg-rose-50 dark:bg-orange-950 rounded-full flex items-center justify-center mx-auto mb-4">
              <Heart size={28} className="text-rose-400 dark:text-orange-400" />
            </div>
            <p className="text-gray-500 dark:text-neutral-400 text-sm mb-4">No favorites yet</p>
            <p className="text-xs text-gray-400 dark:text-neutral-500">
              Save designs you love by tapping the heart icon
            </p>
          </div>
        )}
      </div>
    </div>
  );
}

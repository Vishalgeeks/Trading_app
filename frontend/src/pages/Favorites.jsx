import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { favoriteService } from '../services/favoriteService';
import { designService } from '../services/designService';
import DesignCard from '../components/designs/DesignCard';
import { Heart } from 'lucide-react';

export default function Favorites() {
  const [designs, setDesigns] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const navigate = useNavigate();

  useEffect(() => {
    let cancelled = false;

    async function load() {
      setLoading(true);
      setError('');

      const result = await favoriteService.listFavorites();
      if (cancelled) return;

      if (result.error) {
        if (result.status === 401) {
          navigate('/login');
          return;
        }
        setError(result.message || 'Failed to load favorites');
        setLoading(false);
        return;
      }

      const favorites = result.data || [];
      if (favorites.length === 0) {
        setDesigns([]);
        setLoading(false);
        return;
      }

      const designPromises = favorites.map((fav) => designService.getDesign(fav.design_id));
      const designResults = await Promise.all(designPromises);
      const loadedDesigns = designResults
        .filter((r) => !r.error && r.data)
        .map((r) => ({ ...r.data, favorite: true }));

      if (!cancelled) {
        setDesigns(loadedDesigns);
      }
      setLoading(false);
    }

    load();

    return () => {
      cancelled = true;
    };
  }, [navigate]);

  const handleToggleFavorite = async (designId) => {
    setDesigns((prev) => {
      const exists = prev.find((d) => d.id === designId);
      if (exists) {
        favoriteService.removeFavorite(designId);
        return prev.filter((d) => d.id !== designId);
      }
      return prev;
    });
  };

  if (loading) {
    return (
      <div className="pb-24">
        <div className="px-5 pt-6">
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white mb-6">Favorites</h1>
          <div className="grid grid-cols-2 gap-3">
            {[1, 2, 3, 4].map((i) => (
              <div key={i} className="h-48 bg-gray-200 dark:bg-neutral-800 rounded-2xl animate-pulse" />
            ))}
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="pb-24">
      <div className="px-5 pt-6">
        <h1 className="text-2xl font-bold text-gray-900 dark:text-white mb-1">Favorites</h1>
        <p className="text-sm text-gray-500 dark:text-neutral-400 mb-6">
          {designs.length} saved design{designs.length !== 1 ? 's' : ''}
        </p>

        {error && (
          <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 px-4 py-3 rounded-xl mb-4 text-sm">
            {error}
          </div>
        )}

        {designs.length > 0 ? (
          <div className="grid grid-cols-2 gap-3">
            {designs.map((design) => (
              <DesignCard
                key={design.id}
                design={design}
                onFavoriteToggle={handleToggleFavorite}
              />
            ))}
          </div>
        ) : (
          <div className="text-center py-16">
            <div className="w-16 h-16 bg-rose-50 dark:bg-orange-950 rounded-full flex items-center justify-center mx-auto mb-4">
              <Heart size={28} className="text-rose-400 dark:text-orange-400" />
            </div>
            <p className="text-gray-500 dark:text-neutral-400 text-sm mb-4">No favorites yet</p>
            <button
              onClick={() => navigate('/browse')}
              className="text-sm text-rose-500 dark:text-orange-400 font-medium"
            >
              Browse designs to save favorites
            </button>
          </div>
        )}
      </div>
    </div>
  );
}

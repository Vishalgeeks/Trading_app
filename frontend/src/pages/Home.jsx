import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { categoryService } from '../services/categoryService';
import { designService } from '../services/designService';
import Button from '../components/ui/Button';
import DesignCard from '../components/designs/DesignCard';

const categoryIcons = {
  bridal: '💍',
  party: '🎉',
  casual: '☕',
  festive: '🪔',
};

function getCategoryIcon(slug, name) {
  if (categoryIcons[slug]) return categoryIcons[slug];
  const lower = (name || '').toLowerCase();
  if (lower.includes('bridal')) return '💍';
  if (lower.includes('party')) return '🎉';
  if (lower.includes('casual')) return '☕';
  if (lower.includes('festive')) return '🪔';
  return '✨';
}

export default function Home() {
  const [categories, setCategories] = useState([]);
  const [designs, setDesigns] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    let cancelled = false;

    async function load() {
      setLoading(true);
      setError('');

      const [catResult, designResult] = await Promise.all([
        categoryService.listCategories(),
        designService.listDesigns({ limit: 20 }),
      ]);

      if (cancelled) return;

      if (catResult.error) {
        setError(catResult.message);
      } else {
        setCategories(catResult.data || []);
      }

      if (designResult.error) {
        setError(designResult.message);
      } else {
        setDesigns(designResult.data || []);
      }

      setLoading(false);
    }

    load();
    return () => {
      cancelled = true;
    };
  }, []);

  const featuredDesigns = designs.slice(0, 4);
  const bridalDesigns = designs.filter((d) => d.category?.slug === 'bridal');

  if (loading) {
    return (
      <div className="pb-24">
        <div className="relative bg-gradient-to-b from-rose-500 to-rose-600 dark:from-orange-600 dark:to-orange-700 text-white">
          <div className="px-5 pt-12 pb-8">
            <div className="flex items-center gap-3 mb-6">
              <div className="w-12 h-12 rounded-full bg-white/20 animate-pulse" />
              <div className="flex-1">
                <div className="h-4 bg-white/20 rounded w-24 mb-2 animate-pulse" />
                <div className="h-5 bg-white/20 rounded w-40 animate-pulse" />
              </div>
            </div>
            <div className="h-8 bg-white/20 rounded w-72 mb-3 animate-pulse" />
            <div className="h-4 bg-white/20 rounded w-full mb-6 animate-pulse" />
            <div className="h-12 bg-white rounded-xl animate-pulse" />
          </div>
        </div>
        <div className="px-5 pt-8">
          <div className="h-6 bg-gray-200 dark:bg-neutral-800 rounded w-32 mb-4 animate-pulse" />
          <div className="flex gap-3 overflow-x-auto pb-2">
            {[1, 2, 3, 4].map((i) => (
              <div key={i} className="flex flex-col items-center gap-2 min-w-[72px]">
                <div className="w-14 h-14 rounded-2xl bg-gray-200 dark:bg-neutral-800 animate-pulse" />
                <div className="h-3 bg-gray-200 dark:bg-neutral-800 rounded w-12 animate-pulse" />
              </div>
            ))}
          </div>
          <div className="mt-8">
            <div className="h-6 bg-gray-200 dark:bg-neutral-800 rounded w-36 mb-4 animate-pulse" />
            <div className="grid grid-cols-2 gap-3">
              {[1, 2, 3, 4].map((i) => (
                <div key={i} className="h-48 bg-gray-200 dark:bg-neutral-800 rounded-2xl animate-pulse" />
              ))}
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (error && categories.length === 0 && designs.length === 0) {
    return (
      <div className="pb-24">
        <div className="px-5 pt-12 text-center">
          <p className="text-red-500 text-sm mb-4">{error}</p>
          <button onClick={() => window.location.reload()} className="text-rose-500 text-sm font-medium">
            Retry
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="pb-24">
      <div className="relative bg-gradient-to-b from-rose-500 to-rose-600 dark:from-orange-600 dark:to-orange-700 text-white">
        <div className="px-5 pt-12 pb-8">
          <div className="flex items-center gap-3 mb-6">
            <div className="w-12 h-12 rounded-full bg-white/20 backdrop-blur-sm flex items-center justify-center text-2xl">
              👋
            </div>
            <div>
              <p className="text-sm text-white/80">Welcome back</p>
              <h1 className="text-xl font-bold">Discover Mehndi Art</h1>
            </div>
          </div>

          <h2 className="text-2xl font-bold mb-3 leading-tight">
            Beautiful Henna for Your Special Moments
          </h2>
          <p className="text-white/80 text-sm mb-6 leading-relaxed">
            Book stunning mehndi designs online with easy scheduling and premium service.
          </p>
          <Link to="/browse">
            <Button size="lg" className="w-full bg-white text-rose-600 hover:bg-gray-100 dark:bg-orange-400 dark:text-gray-900 dark:hover:bg-orange-300">
              Browse Designs
            </Button>
          </Link>
        </div>
        <div className="absolute -bottom-6 left-0 right-0">
          <div className="w-16 h-4 bg-gray-50 dark:bg-black mx-auto rounded-b-2xl"></div>
        </div>
      </div>

      <div className="px-5 pt-8">
        <div className="mb-8">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-lg font-bold text-gray-900 dark:text-white">Categories</h2>
            <Link to="/browse" className="text-xs text-rose-500 dark:text-orange-400 font-medium">
              View All
            </Link>
          </div>
          <div className="flex gap-3 overflow-x-auto pb-2 -mx-5 px-5 scrollbar-hide">
            {categories.map((cat) => (
              <Link
                key={cat.id}
                to={`/browse?category=${cat.id}`}
                className="flex flex-col items-center gap-2 min-w-[72px]"
              >
                <div className="w-14 h-14 rounded-2xl bg-rose-50 dark:bg-orange-950 flex items-center justify-center text-2xl border border-rose-100 dark:border-orange-900">
                  {getCategoryIcon(cat.slug, cat.name)}
                </div>
                <span className="text-xs font-medium text-gray-700 dark:text-neutral-300 whitespace-nowrap">
                  {cat.name}
                </span>
              </Link>
            ))}
          </div>
        </div>

        {featuredDesigns.length > 0 && (
          <div className="mb-8">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-bold text-gray-900 dark:text-white">Featured Designs</h2>
              <Link to="/browse" className="text-xs text-rose-500 dark:text-orange-400 font-medium">
                See All
              </Link>
            </div>
            <div className="grid grid-cols-2 gap-3">
              {featuredDesigns.map((design) => (
                <DesignCard key={design.id} design={design} />
              ))}
            </div>
          </div>
        )}

        {bridalDesigns.length > 0 && (
          <div className="mb-8">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-bold text-gray-900 dark:text-white">Bridal Collection</h2>
              <Link to={`/browse?category=${bridalDesigns[0]?.category_id}`} className="text-xs text-rose-500 dark:text-orange-400 font-medium">
                View All
              </Link>
            </div>
            <div className="grid grid-cols-2 gap-3">
              {bridalDesigns.map((design) => (
                <DesignCard key={design.id} design={design} />
              ))}
            </div>
          </div>
        )}

        {designs.length === 0 && !error && (
          <div className="text-center py-16">
            <p className="text-gray-500 dark:text-neutral-400 text-sm">No designs available yet.</p>
          </div>
        )}
      </div>
    </div>
  );
}

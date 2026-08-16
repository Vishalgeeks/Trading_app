import { useEffect, useState, useCallback, useRef } from 'react';
import { useSearchParams } from 'react-router-dom';
import { categoryService } from '../services/categoryService';
import { designService } from '../services/designService';
import DesignCard from '../components/designs/DesignCard';
import { Search, SlidersHorizontal } from 'lucide-react';

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

export default function Browse() {
  const [searchParams, setSearchParams] = useSearchParams();
  const [categories, setCategories] = useState([]);
  const [designs, setDesigns] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [searchQuery, setSearchQuery] = useState('');
  const [count, setCount] = useState(0);

  const categoryId = searchParams.get('category') || 'all';
  const searchTimerRef = useRef(null);

  const loadDesigns = useCallback(async (catId, query) => {
    setLoading(true);
    setError('');

    const params = {};
    if (catId && catId !== 'all') {
      params.category_id = catId;
    }
    if (query && query.trim()) {
      params.q = query.trim();
    }

    const result = await designService.listDesigns(params);
    if (result.error) {
      setError(result.message);
      setDesigns([]);
      setCount(0);
    } else {
      setDesigns(result.data || []);
      setCount(result.count || 0);
    }
    setLoading(false);
  }, []);

  useEffect(() => {
    let cancelled = false;

    async function loadCategories() {
      const result = await categoryService.listCategories();
      if (!cancelled && !result.error) {
        setCategories(result.data || []);
      }
    }

    loadCategories();
    return () => {
      cancelled = true;
    };
  }, []);

  useEffect(() => {
    if (searchTimerRef.current) {
      clearTimeout(searchTimerRef.current);
    }

    searchTimerRef.current = setTimeout(() => {
      loadDesigns(categoryId, searchQuery);
    }, 300);

    return () => {
      if (searchTimerRef.current) {
        clearTimeout(searchTimerRef.current);
      }
    };
  }, [categoryId, searchQuery, loadDesigns]);

  const handleCategoryChange = (catId) => {
    if (catId === 'all') {
      searchParams.delete('category');
    } else {
      searchParams.set('category', catId);
    }
    setSearchParams(searchParams);
  };

  return (
    <div className="pb-24">
      <div className="px-5 pt-6">
        <h1 className="text-2xl font-bold text-gray-900 dark:text-white mb-4">Browse Designs</h1>

        <div className="relative mb-4">
          <Search className="absolute left-4 top-1/2 -translate-y-1/2 text-gray-400" size={18} />
          <input
            type="text"
            placeholder="Search designs..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="w-full pl-11 pr-4 py-3 bg-white dark:bg-neutral-800 border border-rose-200 dark:border-neutral-700 rounded-xl text-sm text-gray-900 dark:text-white placeholder-gray-400 focus:ring-2 focus:ring-rose-500 focus:border-transparent outline-none"
          />
        </div>

        <div className="flex gap-2 overflow-x-auto pb-3 -mx-5 px-5 scrollbar-hide mb-2">
          <button
            onClick={() => handleCategoryChange('all')}
            className={`px-4 py-2 rounded-full text-sm font-medium whitespace-nowrap transition-all duration-200 ${
              categoryId === 'all'
                ? 'bg-rose-500 text-white dark:bg-orange-500'
                : 'bg-white dark:bg-neutral-800 text-gray-600 dark:text-neutral-400 border border-rose-100 dark:border-neutral-700 hover:border-rose-300'
            }`}
          >
            ✨ All
          </button>
          {categories.map((cat) => (
            <button
              key={cat.id}
              onClick={() => handleCategoryChange(cat.id)}
              className={`px-4 py-2 rounded-full text-sm font-medium whitespace-nowrap transition-all duration-200 ${
                categoryId === cat.id
                  ? 'bg-rose-500 text-white dark:bg-orange-500'
                  : 'bg-white dark:bg-neutral-800 text-gray-600 dark:text-neutral-400 border border-rose-100 dark:border-neutral-700 hover:border-rose-300'
              }`}
            >
              {getCategoryIcon(cat.slug, cat.name)} {cat.name}
            </button>
          ))}
        </div>
      </div>

      <div className="px-5">
        <div className="flex items-center justify-between mb-4">
          <p className="text-sm text-gray-500 dark:text-neutral-400">
            {loading ? 'Loading...' : `${count} design${count !== 1 ? 's' : ''} found`}
          </p>
          <button className="flex items-center gap-1 text-xs text-gray-500 dark:text-neutral-400">
            <SlidersHorizontal size={14} />
            Filter
          </button>
        </div>

        {error && !loading && (
          <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 px-4 py-3 rounded-xl mb-4 text-sm">
            {error}
          </div>
        )}

        {!loading && designs.length > 0 ? (
          <div className="grid grid-cols-2 gap-3">
            {designs.map((design) => (
              <DesignCard key={design.id} design={design} onPress={() => {}} />
            ))}
          </div>
        ) : (
          !loading && (
            <div className="text-center py-12">
              <p className="text-gray-500 dark:text-neutral-400">
                {searchQuery || categoryId !== 'all'
                  ? 'No designs found matching your criteria.'
                  : 'No designs available yet.'}
              </p>
            </div>
          )
        )}
      </div>
    </div>
  );
}

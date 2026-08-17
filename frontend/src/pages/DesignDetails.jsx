import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { designService } from '../services/designService';
import Button from '../components/ui/Button';
import DesignCard from '../components/designs/DesignCard';
import { Heart, Clock, IndianRupee, ArrowLeft } from 'lucide-react';

function formatDuration(minutes) {
  if (!minutes || minutes <= 0) return '';
  const hours = Math.floor(minutes / 60);
  const mins = minutes % 60;
  if (hours > 0 && mins > 0) {
    return `${hours}h ${mins}m`;
  }
  if (hours > 0) {
    return `${hours}h`;
  }
  return `${mins}m`;
}

export default function DesignDetails() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [design, setDesign] = useState(null);
  const [relatedDesigns, setRelatedDesigns] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    let cancelled = false;

    async function load() {
      setLoading(true);
      setError('');
      setDesign(null);
      setRelatedDesigns([]);

      const result = await designService.getDesign(id);
      if (cancelled) return;

      if (result.error) {
        setError(result.message || 'Failed to load design');
        setLoading(false);
        return;
      }

      setDesign(result.data);

      if (result.data?.category_id) {
        const relatedResult = await designService.listDesignsByCategory(result.data.category_id);
        if (!cancelled && !relatedResult.error) {
          setRelatedDesigns((relatedResult.data || []).filter((d) => d.id !== result.data.id).slice(0, 4));
        }
      }

      setLoading(false);
    }

    load();

    return () => {
      cancelled = true;
    };
  }, [id]);

  if (loading) {
    return (
      <div className="pb-24">
        <div className="relative aspect-[4/3] bg-rose-50 dark:bg-neutral-800 animate-pulse" />
        <div className="px-5 pt-6">
          <div className="h-8 bg-gray-200 dark:bg-neutral-800 rounded w-48 mb-4 animate-pulse" />
          <div className="h-4 bg-gray-200 dark:bg-neutral-800 rounded w-full mb-6 animate-pulse" />
          <div className="flex gap-4 mb-8">
            <div className="flex-1 h-24 bg-gray-200 dark:bg-neutral-800 rounded-xl animate-pulse" />
            <div className="flex-1 h-24 bg-gray-200 dark:bg-neutral-800 rounded-xl animate-pulse" />
          </div>
        </div>
      </div>
    );
  }

  if (error || !design) {
    return (
      <div className="flex items-center justify-center h-screen">
        <div className="text-center px-5">
          <p className="text-gray-500 dark:text-neutral-400 mb-4">{error || 'Design not found'}</p>
          <button onClick={() => navigate('/browse')} className="text-rose-500 text-sm font-medium">
            Browse designs
          </button>
        </div>
      </div>
    );
  }

  const imageUrl = design.image_url;
  const categoryName = design.category?.name || 'Uncategorized';

  return (
    <div className="pb-24">
      <div className="relative aspect-[4/3] bg-rose-50 dark:bg-neutral-800 flex items-center justify-center">
        {imageUrl ? (
          <img
            src={imageUrl}
            alt={design.name}
            className="w-full h-full object-cover"
            onError={(e) => {
              e.target.style.display = 'none';
            }}
          />
        ) : (
          <span className="text-8xl">🎨</span>
        )}
        <button
          onClick={() => navigate(-1)}
          className="absolute top-4 left-4 w-10 h-10 bg-white/90 dark:bg-neutral-800/90 backdrop-blur-sm rounded-full flex items-center justify-center shadow-sm"
        >
          <ArrowLeft size={20} className="text-gray-700 dark:text-white" />
        </button>
        <button className="absolute top-4 right-4 w-10 h-10 bg-white/90 dark:bg-neutral-800/90 backdrop-blur-sm rounded-full flex items-center justify-center shadow-sm">
          <Heart size={20} className="text-gray-700 dark:text-white" />
        </button>
      </div>

      <div className="px-5 pt-6">
        <div className="flex items-start justify-between mb-2">
          <div>
            <span className="inline-block px-2.5 py-1 bg-rose-100 dark:bg-orange-950 text-rose-700 dark:text-orange-400 rounded-full text-xs font-medium mb-2">
              {categoryName}
            </span>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-white">{design.name}</h1>
          </div>
          <div className="text-right">
            <p className="text-2xl font-bold text-rose-600 dark:text-orange-400">₹{design.price}</p>
          </div>
        </div>

        <p className="text-gray-600 dark:text-neutral-400 text-sm leading-relaxed mt-3 mb-6">
          {design.description || 'No description available.'}
        </p>

        <div className="flex gap-4 mb-8">
          <div className="flex-1 bg-rose-50 dark:bg-neutral-800 rounded-xl p-4 flex items-center gap-3">
            <div className="w-10 h-10 rounded-full bg-rose-100 dark:bg-orange-950 flex items-center justify-center text-rose-500 dark:text-orange-400">
              <Clock size={18} />
            </div>
            <div>
              <p className="text-xs text-gray-500 dark:text-neutral-400">Duration</p>
              <p className="text-sm font-semibold text-gray-900 dark:text-white">{formatDuration(design.duration_minutes)}</p>
            </div>
          </div>
          <div className="flex-1 bg-rose-50 dark:bg-neutral-800 rounded-xl p-4 flex items-center gap-3">
            <div className="w-10 h-10 rounded-full bg-rose-100 dark:bg-orange-950 flex items-center justify-center text-rose-500 dark:text-orange-400">
              <IndianRupee size={18} />
            </div>
            <div>
              <p className="text-xs text-gray-500 dark:text-neutral-400">Starting at</p>
              <p className="text-sm font-semibold text-gray-900 dark:text-white">₹{design.price}</p>
            </div>
          </div>
        </div>

        <div className="mb-8">
          <Button className="w-full" onClick={() => navigate(`/booking/${design.id}`)}>
            Book Now
          </Button>
        </div>

        {relatedDesigns.length > 0 && (
          <div>
            <h3 className="text-lg font-bold text-gray-900 dark:text-white mb-4">You May Also Like</h3>
            <div className="grid grid-cols-2 gap-3">
              {relatedDesigns.map((d) => (
                <DesignCard key={d.id} design={d} onPress={() => navigate(`/designs/${d.id}`)} />
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

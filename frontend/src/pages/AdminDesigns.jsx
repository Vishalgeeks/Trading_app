import { mockDesigns } from '../mockData';
import DesignCard from '../components/designs/DesignCard';
import Button from '../components/ui/Button';
import { Plus, Edit, Trash2 } from 'lucide-react';

export default function AdminDesigns() {
  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white mb-1">Designs</h1>
          <p className="text-sm text-gray-500 dark:text-neutral-400">
            Manage your design catalog
          </p>
        </div>
        <Button size="sm" className="gap-2">
          <Plus size={16} />
          Add New
        </Button>
      </div>

      <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
        {mockDesigns.map((design) => (
          <div key={design.id} className="relative group">
            <DesignCard design={design} />
            <div className="absolute inset-0 bg-black/20 dark:bg-black/40 rounded-2xl opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center gap-2">
              <button className="w-8 h-8 bg-white rounded-full flex items-center justify-center shadow-lg">
                <Edit size={14} className="text-gray-700" />
              </button>
              <button className="w-8 h-8 bg-red-500 rounded-full flex items-center justify-center shadow-lg">
                <Trash2 size={14} className="text-white" />
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

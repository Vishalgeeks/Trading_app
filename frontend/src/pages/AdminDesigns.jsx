import { useEffect, useState, useMemo } from 'react';
import { designService } from '../services/designService';
import { categoryService } from '../services/categoryService';
import Button from '../components/ui/Button';
import DesignCard from '../components/designs/DesignCard';
import { Plus, Edit, Trash2, X, Search } from 'lucide-react';

const EMPTY_FORM = {
  name: '',
  slug: '',
  description: '',
  image_url: '',
  price: '',
  duration_minutes: '',
  category_id: '',
};

export default function AdminDesigns() {
  const [designs, setDesigns] = useState([]);
  const [categories, setCategories] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [search, setSearch] = useState('');
  const [showModal, setShowModal] = useState(false);
  const [editing, setEditing] = useState(null);
  const [form, setForm] = useState(EMPTY_FORM);
  const [saving, setSaving] = useState(false);
  const [formError, setFormError] = useState('');
  const [deleting, setDeleting] = useState(null);
  const [imageFile, setImageFile] = useState(null);
  const [imagePreview, setImagePreview] = useState('');

  useEffect(() => {
    let cancelled = false;

    async function load() {
      setLoading(true);
      setError('');

      const [designsResult, categoriesResult] = await Promise.all([
        designService.listAdminDesigns(),
        categoryService.listAdminCategories(),
      ]);

      if (cancelled) return;

      if (designsResult.error) {
        setError(designsResult.message || 'Failed to load designs');
        setLoading(false);
        return;
      }

      setDesigns(designsResult.data || []);
      setCategories(categoriesResult.error ? [] : (categoriesResult.data || []));
      setLoading(false);
    }

    load();

    return () => {
      cancelled = false;
    };
  }, []);

  const filteredDesigns = useMemo(() => {
    if (!search.trim()) return designs;
    const q = search.trim().toLowerCase();
    return designs.filter((d) =>
      d.name.toLowerCase().includes(q) ||
      (d.description || '').toLowerCase().includes(q)
    );
  }, [designs, search]);

  const openCreate = () => {
    setEditing(null);
    setForm(EMPTY_FORM);
    setImageFile(null);
    setImagePreview('');
    setFormError('');
    setShowModal(true);
  };

  const openEdit = (design) => {
    setEditing(design);
    setForm({
      name: design.name || '',
      slug: design.slug || '',
      description: design.description || '',
      image_url: design.image_url || '',
      price: design.price || '',
      duration_minutes: design.duration_minutes || '',
      category_id: design.category_id || '',
    });
    setImageFile(null);
    setImagePreview(design.image_url || '');
    setFormError('');
    setShowModal(true);
  };

  const closeModal = () => {
    setShowModal(false);
    setEditing(null);
    setForm(EMPTY_FORM);
    setImageFile(null);
    setImagePreview('');
    setFormError('');
    setSaving(false);
  };

  const handleImageChange = (e) => {
    const file = e.target.files?.[0];
    if (!file) return;
    setImageFile(file);
    setImagePreview(URL.createObjectURL(file));
    setForm((prev) => ({ ...prev, image_url: '' }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setSaving(true);
    setFormError('');

    const payload = {
      name: form.name.trim(),
      slug: form.slug.trim(),
      description: form.description.trim() || null,
      price: form.price.trim(),
      duration_minutes: parseInt(form.duration_minutes, 10) || 0,
      category_id: form.category_id,
    };

    if (!payload.name || !payload.slug || !payload.price || !payload.category_id) {
      setFormError('Please fill in all required fields.');
      setSaving(false);
      return;
    }

    if (isNaN(payload.duration_minutes) || payload.duration_minutes <= 0) {
      setFormError('Duration must be a positive number.');
      setSaving(false);
      return;
    }

    let result;
    if (editing) {
      result = await designService.updateDesign(editing.id, payload, imageFile);
    } else {
      result = await designService.createDesign(payload, imageFile);
    }

    if (result.error) {
      setFormError(result.message || 'Failed to save design.');
      setSaving(false);
      return;
    }

    if (editing) {
      setDesigns((prev) => prev.map((d) => (d.id === editing.id ? result.data : d)));
    } else {
      setDesigns((prev) => [result.data, ...prev]);
    }

    closeModal();
  };

  const handleDelete = async (design) => {
    if (!window.confirm(`Deactivate "${design.name}"?`)) return;
    setDeleting(design.id);
    const result = await designService.deleteDesign(design.id);
    setDeleting(null);
    if (result.error) {
      setError(result.message || 'Failed to deactivate design.');
      return;
    }
    setDesigns((prev) => prev.filter((d) => d.id !== design.id));
  };

  const handleToggleActive = async (design) => {
    const result = await designService.updateDesign(design.id, {
      ...design,
      is_active: !design.is_active,
    });
    if (result.error) {
      setError(result.message || 'Failed to update design status.');
      return;
    }
    setDesigns((prev) => prev.map((d) => (d.id === design.id ? result.data : d)));
  };

  if (loading) {
    return (
      <div>
        <div className="flex items-center justify-between mb-6">
          <div>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-white mb-1">Designs</h1>
            <p className="text-sm text-gray-500 dark:text-neutral-400">Loading...</p>
          </div>
        </div>
        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
          {[1, 2, 3, 4].map((i) => (
            <div key={i} className="h-64 bg-gray-200 dark:bg-neutral-800 rounded-2xl animate-pulse" />
          ))}
        </div>
      </div>
    );
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white mb-1">Designs</h1>
          <p className="text-sm text-gray-500 dark:text-neutral-400">
            Manage your design catalog ({designs.length})
          </p>
        </div>
        <Button size="sm" className="gap-2" onClick={openCreate}>
          <Plus size={16} />
          Add New
        </Button>
      </div>

      {error && (
        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 px-4 py-3 rounded-xl mb-4 text-sm">
          {error}
        </div>
      )}

      <div className="mb-6">
        <div className="relative">
          <Search size={18} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
          <input
            type="text"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Search designs..."
            className="w-full pl-10 pr-4 py-2 bg-white dark:bg-neutral-800 border border-rose-200 dark:border-neutral-700 rounded-xl text-sm text-gray-900 dark:text-white focus:ring-2 focus:ring-rose-500 focus:border-transparent outline-none"
          />
        </div>
      </div>

      {filteredDesigns.length > 0 ? (
        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
          {filteredDesigns.map((design) => (
            <div key={design.id} className="relative group">
              <DesignCard design={design} />
              <div className="absolute inset-0 bg-black/20 dark:bg-black/40 rounded-2xl opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center gap-2">
                <button
                  onClick={() => openEdit(design)}
                  className="w-8 h-8 bg-white rounded-full flex items-center justify-center shadow-lg"
                >
                  <Edit size={14} className="text-gray-700" />
                </button>
                <button
                  onClick={() => handleDelete(design)}
                  disabled={deleting === design.id}
                  className="w-8 h-8 bg-red-500 rounded-full flex items-center justify-center shadow-lg disabled:opacity-50"
                >
                  <Trash2 size={14} className="text-white" />
                </button>
              </div>
              <div className="absolute top-2 right-2">
                <button
                  onClick={() => handleToggleActive(design)}
                  className={`px-2 py-1 rounded-full text-[10px] font-semibold ${
                    design.is_active
                      ? 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400'
                      : 'bg-gray-200 text-gray-600 dark:bg-neutral-700 dark:text-neutral-400'
                  }`}
                >
                  {design.is_active ? 'Active' : 'Inactive'}
                </button>
              </div>
            </div>
          ))}
        </div>
      ) : (
        <div className="text-center py-12">
          <p className="text-gray-500 dark:text-neutral-400 text-sm">
            {search ? 'No designs match your search.' : 'No designs yet.'}
          </p>
        </div>
      )}

      {showModal && (
        <div className="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4">
          <div className="bg-white dark:bg-neutral-900 rounded-2xl border border-rose-100 dark:border-neutral-800 shadow-xl w-full max-w-lg max-h-[90vh] overflow-y-auto">
            <div className="flex items-center justify-between p-5 border-b border-rose-100 dark:border-neutral-800">
              <h2 className="text-lg font-bold text-gray-900 dark:text-white">
                {editing ? 'Edit Design' : 'New Design'}
              </h2>
              <button onClick={closeModal} className="text-gray-400 hover:text-gray-600">
                <X size={20} />
              </button>
            </div>

            <form onSubmit={handleSubmit} className="p-5 space-y-4">
              {formError && (
                <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 px-4 py-3 rounded-xl text-sm">
                  {formError}
                </div>
              )}

              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-neutral-300 mb-1.5">Name *</label>
                <input
                  type="text"
                  value={form.name}
                  onChange={(e) => setForm({ ...form, name: e.target.value })}
                  className="w-full px-4 py-3 bg-white dark:bg-neutral-800 border border-rose-200 dark:border-neutral-700 rounded-xl text-gray-900 dark:text-white focus:ring-2 focus:ring-rose-500 focus:border-transparent outline-none"
                  placeholder="Design name"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-neutral-300 mb-1.5">Slug *</label>
                <input
                  type="text"
                  value={form.slug}
                  onChange={(e) => setForm({ ...form, slug: e.target.value })}
                  className="w-full px-4 py-3 bg-white dark:bg-neutral-800 border border-rose-200 dark:border-neutral-700 rounded-xl text-gray-900 dark:text-white focus:ring-2 focus:ring-rose-500 focus:border-transparent outline-none"
                  placeholder="e.g. bridal-royal"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-neutral-300 mb-1.5">Category *</label>
                <select
                  value={form.category_id}
                  onChange={(e) => setForm({ ...form, category_id: e.target.value })}
                  className="w-full px-4 py-3 bg-white dark:bg-neutral-800 border border-rose-200 dark:border-neutral-700 rounded-xl text-gray-900 dark:text-white focus:ring-2 focus:ring-rose-500 focus:border-transparent outline-none"
                >
                  <option value="">Select category</option>
                  {categories.map((cat) => (
                    <option key={cat.id} value={cat.id}>{cat.name}</option>
                  ))}
                </select>
              </div>

              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-neutral-300 mb-1.5">Price (?) *</label>
                  <input
                    type="text"
                    value={form.price}
                    onChange={(e) => setForm({ ...form, price: e.target.value })}
                    className="w-full px-4 py-3 bg-white dark:bg-neutral-800 border border-rose-200 dark:border-neutral-700 rounded-xl text-gray-900 dark:text-white focus:ring-2 focus:ring-rose-500 focus:border-transparent outline-none"
                    placeholder="2500"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-neutral-300 mb-1.5">Duration (mins) *</label>
                  <input
                    type="number"
                    value={form.duration_minutes}
                    onChange={(e) => setForm({ ...form, duration_minutes: e.target.value })}
                    className="w-full px-4 py-3 bg-white dark:bg-neutral-800 border border-rose-200 dark:border-neutral-700 rounded-xl text-gray-900 dark:text-white focus:ring-2 focus:ring-rose-500 focus:border-transparent outline-none"
                    placeholder="120"
                  />
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-neutral-300 mb-1.5">Image</label>
                <input
                  type="file"
                  accept="image/*"
                  onChange={handleImageChange}
                  className="w-full px-4 py-3 bg-white dark:bg-neutral-800 border border-rose-200 dark:border-neutral-700 rounded-xl text-gray-900 dark:text-white focus:ring-2 focus:ring-rose-500 focus:border-transparent outline-none file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-rose-50 file:text-rose-700 hover:file:bg-rose-100"
                />
                {imagePreview && (
                  <div className="mt-3">
                    <img
                      src={imagePreview}
                      alt="Preview"
                      className="w-full h-40 object-cover rounded-xl border border-rose-100 dark:border-neutral-700"
                    />
                  </div>
                )}
                {!imagePreview && form.image_url && (
                  <div className="mt-3">
                    <img
                      src={form.image_url}
                      alt="Current"
                      className="w-full h-40 object-cover rounded-xl border border-rose-100 dark:border-neutral-700"
                      onError={(e) => { e.target.style.display = 'none'; }}
                    />
                  </div>
                )}
                <p className="text-xs text-gray-400 dark:text-neutral-500 mt-1">
                  Upload an image file or paste a URL below.
                </p>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-neutral-300 mb-1.5">Image URL</label>
                <input
                  type="text"
                  value={form.image_url}
                  onChange={(e) => {
                    setForm({ ...form, image_url: e.target.value });
                    if (e.target.value) {
                      setImagePreview(e.target.value);
                      setImageFile(null);
                    }
                  }}
                  className="w-full px-4 py-3 bg-white dark:bg-neutral-800 border border-rose-200 dark:border-neutral-700 rounded-xl text-gray-900 dark:text-white focus:ring-2 focus:ring-rose-500 focus:border-transparent outline-none"
                  placeholder="https://example.com/image.jpg"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-neutral-300 mb-1.5">Description</label>
                <textarea
                  value={form.description}
                  onChange={(e) => setForm({ ...form, description: e.target.value })}
                  rows={3}
                  className="w-full px-4 py-3 bg-white dark:bg-neutral-800 border border-rose-200 dark:border-neutral-700 rounded-xl text-gray-900 dark:text-white focus:ring-2 focus:ring-rose-500 focus:border-transparent outline-none resize-none"
                  placeholder="Describe the design..."
                />
              </div>

              <div className="flex gap-3 pt-2">
                <Button type="button" variant="outline" className="flex-1" onClick={closeModal}>
                  Cancel
                </Button>
                <Button type="submit" className="flex-1" disabled={saving}>
                  {saving ? 'Saving...' : editing ? 'Update' : 'Create'}
                </Button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}

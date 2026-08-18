import { useEffect, useState } from 'react';
import { categoryService } from '../services/categoryService';
import Button from '../components/ui/Button';
import { Plus, Edit, Trash2, X } from 'lucide-react';

const EMPTY_FORM = {
  name: '',
  slug: '',
  description: '',
  is_active: true,
};

export default function AdminCategories() {
  const [categories, setCategories] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showModal, setShowModal] = useState(false);
  const [editing, setEditing] = useState(null);
  const [form, setForm] = useState(EMPTY_FORM);
  const [saving, setSaving] = useState(false);
  const [formError, setFormError] = useState('');
  const [deleting, setDeleting] = useState(null);

  useEffect(() => {
    let cancelled = false;

    async function load() {
      setLoading(true);
      setError('');
      const result = await categoryService.listAdminCategories();
      if (cancelled) return;
      if (result.error) {
        setError(result.message || 'Failed to load categories');
      } else {
        setCategories(result.data || []);
      }
      setLoading(false);
    }

    load();

    return () => {
      cancelled = true;
    };
  }, []);

  const openCreate = () => {
    setEditing(null);
    setForm(EMPTY_FORM);
    setFormError('');
    setShowModal(true);
  };

  const openEdit = (category) => {
    setEditing(category);
    setForm({
      name: category.name || '',
      slug: category.slug || '',
      description: category.description || '',
      is_active: category.is_active,
    });
    setFormError('');
    setShowModal(true);
  };

  const closeModal = () => {
    setShowModal(false);
    setEditing(null);
    setForm(EMPTY_FORM);
    setFormError('');
    setSaving(false);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setSaving(true);
    setFormError('');

    const payload = {
      name: form.name.trim(),
      slug: form.slug.trim(),
      description: form.description.trim() || null,
      is_active: form.is_active,
    };

    if (!payload.name || !payload.slug) {
      setFormError('Name and slug are required.');
      setSaving(false);
      return;
    }

    let result;
    if (editing) {
      result = await categoryService.updateCategory(editing.id, payload);
    } else {
      result = await categoryService.createCategory(payload);
    }

    if (result.error) {
      setFormError(result.message || 'Failed to save category.');
      setSaving(false);
      return;
    }

    if (editing) {
      setCategories((prev) => prev.map((c) => (c.id === editing.id ? result.data : c)));
    } else {
      setCategories((prev) => [result.data, ...prev]);
    }

    closeModal();
  };

  const handleDelete = async (category) => {
    if (!window.confirm(`Deactivate "${category.name}"?`)) return;
    setDeleting(category.id);
    const result = await categoryService.deleteCategory(category.id);
    setDeleting(null);
    if (result.error) {
      setError(result.message || 'Failed to deactivate category.');
      return;
    }
    setCategories((prev) => prev.filter((c) => c.id !== category.id));
  };

  if (loading) {
    return (
      <div>
        <div className="flex items-center justify-between mb-6">
          <div>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-white mb-1">Categories</h1>
            <p className="text-sm text-gray-500 dark:text-neutral-400">Loading...</p>
          </div>
        </div>
        <div className="space-y-3">
          {[1, 2, 3].map((i) => (
            <div key={i} className="h-16 bg-gray-200 dark:bg-neutral-800 rounded-2xl animate-pulse" />
          ))}
        </div>
      </div>
    );
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white mb-1">Categories</h1>
          <p className="text-sm text-gray-500 dark:text-neutral-400">
            Manage design categories ({categories.length})
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

      {categories.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-gray-500 dark:text-neutral-400 text-sm mb-4">
            No categories yet.
          </p>
          <Button size="sm" onClick={openCreate}>Add your first category</Button>
        </div>
      ) : (
        <div className="space-y-3">
          {categories.map((category) => (
            <div
              key={category.id}
              className="bg-white dark:bg-neutral-900 rounded-2xl border border-rose-100 dark:border-neutral-800 p-5 shadow-sm flex items-center justify-between"
            >
              <div>
                <h3 className="font-semibold text-gray-900 dark:text-white text-sm">{category.name}</h3>
                <p className="text-xs text-gray-500 dark:text-neutral-400 mt-0.5">/{category.slug}</p>
                {category.description && (
                  <p className="text-xs text-gray-500 dark:text-neutral-400 mt-1">{category.description}</p>
                )}
              </div>
              <div className="flex items-center gap-2">
                <span
                  className={`px-2 py-1 rounded-full text-[10px] font-semibold ${
                    category.is_active
                      ? 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400'
                      : 'bg-gray-200 text-gray-600 dark:bg-neutral-700 dark:text-neutral-400'
                  }`}
                >
                  {category.is_active ? 'Active' : 'Inactive'}
                </span>
                <button
                  onClick={() => openEdit(category)}
                  className="w-8 h-8 bg-white dark:bg-neutral-700 rounded-full flex items-center justify-center shadow-sm"
                >
                  <Edit size={14} className="text-gray-700 dark:text-white" />
                </button>
                <button
                  onClick={() => handleDelete(category)}
                  disabled={deleting === category.id}
                  className="w-8 h-8 bg-red-500 rounded-full flex items-center justify-center shadow-sm disabled:opacity-50"
                >
                  <Trash2 size={14} className="text-white" />
                </button>
              </div>
            </div>
          ))}
        </div>
      )}

      {showModal && (
        <div className="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4">
          <div className="bg-white dark:bg-neutral-900 rounded-2xl border border-rose-100 dark:border-neutral-800 shadow-xl w-full max-w-sm">
            <div className="flex items-center justify-between p-5 border-b border-rose-100 dark:border-neutral-800">
              <h2 className="text-lg font-bold text-gray-900 dark:text-white">
                {editing ? 'Edit Category' : 'New Category'}
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
                  placeholder="Category name"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-neutral-300 mb-1.5">Slug *</label>
                <input
                  type="text"
                  value={form.slug}
                  onChange={(e) => setForm({ ...form, slug: e.target.value })}
                  className="w-full px-4 py-3 bg-white dark:bg-neutral-800 border border-rose-200 dark:border-neutral-700 rounded-xl text-gray-900 dark:text-white focus:ring-2 focus:ring-rose-500 focus:border-transparent outline-none"
                  placeholder="e.g. bridal"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-neutral-300 mb-1.5">Description</label>
                <textarea
                  value={form.description}
                  onChange={(e) => setForm({ ...form, description: e.target.value })}
                  rows={3}
                  className="w-full px-4 py-3 bg-white dark:bg-neutral-800 border border-rose-200 dark:border-neutral-700 rounded-xl text-gray-900 dark:text-white focus:ring-2 focus:ring-rose-500 focus:border-transparent outline-none resize-none"
                  placeholder="Optional description"
                />
              </div>

              <div className="flex items-center gap-2">
                <input
                  id="is_active"
                  type="checkbox"
                  checked={form.is_active}
                  onChange={(e) => setForm({ ...form, is_active: e.target.checked })}
                  className="w-4 h-4 rounded border-rose-300 text-rose-600 focus:ring-rose-500"
                />
                <label htmlFor="is_active" className="text-sm text-gray-700 dark:text-neutral-300">
                  Active
                </label>
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

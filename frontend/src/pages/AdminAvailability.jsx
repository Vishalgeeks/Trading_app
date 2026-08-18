import { useEffect, useState } from 'react';
import { availabilityService } from '../services/availabilityService';
import Button from '../components/ui/Button';
import { Plus, Edit, Trash2, X, Clock } from 'lucide-react';

const DAY_NAMES = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];
const SHORT_DAYS = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];

const EMPTY_FORM = {
  day_of_week: 1,
  start_time: '09:00',
  end_time: '17:00',
};

export default function AdminAvailability() {
  const [availabilities, setAvailabilities] = useState([]);
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
      const result = await availabilityService.getAdminAvailability(true);
      if (cancelled) return;
      if (result.error) {
        setError(result.message || 'Failed to load availability');
      } else {
        setAvailabilities(result.data || []);
      }
      setLoading(false);
    }

    load();

    return () => {
      cancelled = true;
    };
  }, []);

  const grouped = availabilities.reduce((acc, av) => {
    const day = av.day_of_week;
    if (!acc[day]) acc[day] = [];
    acc[day].push(av);
    return acc;
  }, {});

  const openCreate = () => {
    setEditing(null);
    setForm(EMPTY_FORM);
    setFormError('');
    setShowModal(true);
  };

  const openEdit = (av) => {
    setEditing(av);
    setForm({
      day_of_week: av.day_of_week,
      start_time: av.start_time,
      end_time: av.end_time,
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
      day_of_week: parseInt(form.day_of_week, 10),
      start_time: form.start_time,
      end_time: form.end_time,
    };

    if (payload.day_of_week < 0 || payload.day_of_week > 6) {
      setFormError('Day must be between 0 and 6.');
      setSaving(false);
      return;
    }

    if (payload.start_time >= payload.end_time) {
      setFormError('End time must be after start time.');
      setSaving(false);
      return;
    }

    let result;
    if (editing) {
      result = await availabilityService.updateAvailability(editing.id, payload);
    } else {
      result = await availabilityService.createAvailability(payload);
    }

    if (result.error) {
      setFormError(result.message || 'Failed to save availability.');
      setSaving(false);
      return;
    }

    if (editing) {
      setAvailabilities((prev) => prev.map((a) => (a.id === editing.id ? result.data : a)));
    } else {
      setAvailabilities((prev) => [...prev, result.data]);
    }

    closeModal();
  };

  const handleDelete = async (id) => {
    if (!window.confirm('Remove this availability slot?')) return;
    setDeleting(id);
    const result = await availabilityService.deleteAvailability(id);
    setDeleting(null);
    if (result.error) {
      setError(result.message || 'Failed to delete availability.');
      return;
    }
    setAvailabilities((prev) => prev.filter((a) => a.id !== id));
  };

  if (loading) {
    return (
      <div>
        <h1 className="text-2xl font-bold text-gray-900 dark:text-white mb-1">Schedule</h1>
        <p className="text-sm text-gray-500 dark:text-neutral-400 mb-6">Loading...</p>
        <div className="space-y-4">
          {[1, 2, 3].map((i) => (
            <div key={i} className="h-32 bg-gray-200 dark:bg-neutral-800 rounded-2xl animate-pulse" />
          ))}
        </div>
      </div>
    );
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white mb-1">Schedule</h1>
          <p className="text-sm text-gray-500 dark:text-neutral-400">
            Manage your availability and time slots
          </p>
        </div>
        <Button size="sm" className="gap-2" onClick={openCreate}>
          <Plus size={16} />
          Add Slot
        </Button>
      </div>

      {error && (
        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 px-4 py-3 rounded-xl mb-4 text-sm">
          {error}
        </div>
      )}

      {availabilities.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-gray-500 dark:text-neutral-400 text-sm mb-4">
            No availability configured.
          </p>
          <Button size="sm" onClick={openCreate}>Add your first slot</Button>
        </div>
      ) : (
        <div className="space-y-4">
          {SHORT_DAYS.map((day, index) => {
            const slots = grouped[index] || [];
            if (slots.length === 0) return null;
            return (
              <div
                key={index}
                className="bg-white dark:bg-neutral-900 rounded-2xl border border-rose-100 dark:border-neutral-800 p-5 shadow-sm"
              >
                <h3 className="text-sm font-bold text-gray-900 dark:text-white mb-3">
                  {day} — {DAY_NAMES[index]}
                </h3>
                <div className="space-y-2">
                  {slots.map((slot) => (
                    <div
                      key={slot.id}
                      className="flex items-center justify-between py-3 px-4 bg-rose-50 dark:bg-neutral-800 rounded-xl"
                    >
                      <div className="flex items-center gap-2 text-sm text-gray-700 dark:text-neutral-300">
                        <Clock size={16} className="text-rose-500 dark:text-orange-400" />
                        <span className="font-medium">
                          {slot.start_time} – {slot.end_time}
                        </span>
                      </div>
                      <div className="flex items-center gap-2">
                        <button
                          onClick={() => openEdit(slot)}
                          className="w-8 h-8 bg-white dark:bg-neutral-700 rounded-full flex items-center justify-center shadow-sm"
                        >
                          <Edit size={14} className="text-gray-700 dark:text-white" />
                        </button>
                        <button
                          onClick={() => handleDelete(slot.id)}
                          disabled={deleting === slot.id}
                          className="w-8 h-8 bg-red-500 rounded-full flex items-center justify-center shadow-sm disabled:opacity-50"
                        >
                          <Trash2 size={14} className="text-white" />
                        </button>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            );
          })}
        </div>
      )}

      {showModal && (
        <div className="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4">
          <div className="bg-white dark:bg-neutral-900 rounded-2xl border border-rose-100 dark:border-neutral-800 shadow-xl w-full max-w-sm">
            <div className="flex items-center justify-between p-5 border-b border-rose-100 dark:border-neutral-800">
              <h2 className="text-lg font-bold text-gray-900 dark:text-white">
                {editing ? 'Edit Slot' : 'New Slot'}
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
                <label className="block text-sm font-medium text-gray-700 dark:text-neutral-300 mb-1.5">Day</label>
                <select
                  value={form.day_of_week}
                  onChange={(e) => setForm({ ...form, day_of_week: e.target.value })}
                  className="w-full px-4 py-3 bg-white dark:bg-neutral-800 border border-rose-200 dark:border-neutral-700 rounded-xl text-gray-900 dark:text-white focus:ring-2 focus:ring-rose-500 focus:border-transparent outline-none"
                >
                  {SHORT_DAYS.map((day, i) => (
                    <option key={i} value={i}>{day} — {DAY_NAMES[i]}</option>
                  ))}
                </select>
              </div>

              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-neutral-300 mb-1.5">Start Time</label>
                  <input
                    type="time"
                    value={form.start_time}
                    onChange={(e) => setForm({ ...form, start_time: e.target.value })}
                    className="w-full px-4 py-3 bg-white dark:bg-neutral-800 border border-rose-200 dark:border-neutral-700 rounded-xl text-gray-900 dark:text-white focus:ring-2 focus:ring-rose-500 focus:border-transparent outline-none"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-neutral-300 mb-1.5">End Time</label>
                  <input
                    type="time"
                    value={form.end_time}
                    onChange={(e) => setForm({ ...form, end_time: e.target.value })}
                    className="w-full px-4 py-3 bg-white dark:bg-neutral-800 border border-rose-200 dark:border-neutral-700 rounded-xl text-gray-900 dark:text-white focus:ring-2 focus:ring-rose-500 focus:border-transparent outline-none"
                  />
                </div>
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

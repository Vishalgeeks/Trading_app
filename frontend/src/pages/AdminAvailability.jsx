import { useState } from 'react';
import Button from '../components/ui/Button';
import { Save } from 'lucide-react';

const timeSlots = ['9:00 AM', '10:00 AM', '11:00 AM', '12:00 PM', '1:00 PM', '2:00 PM', '3:00 PM', '4:00 PM'];

export default function AdminAvailability() {
  const [selectedDate, setSelectedDate] = useState('2025-09-15');
  const [slots, setSlots] = useState(() => {
    const initial = {};
    timeSlots.forEach((slot) => {
      initial[slot] = Math.random() > 0.5 ? 'available' : 'blocked';
    });
    return initial;
  });

  const toggleSlot = (slot) => {
    setSlots((prev) => ({
      ...prev,
      [slot]: prev[slot] === 'available' ? 'blocked' : 'available',
    }));
  };

  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-900 dark:text-white mb-1">Schedule</h1>
      <p className="text-sm text-gray-500 dark:text-neutral-400 mb-6">
        Manage your availability and time slots
      </p>

      <div className="bg-white dark:bg-neutral-900 rounded-2xl border border-rose-100 dark:border-neutral-800 p-5 shadow-sm mb-6">
        <label className="block text-sm font-medium text-gray-700 dark:text-neutral-300 mb-2">
          Select Date
        </label>
        <input
          type="date"
          value={selectedDate}
          onChange={(e) => setSelectedDate(e.target.value)}
          className="w-full px-4 py-3 bg-white dark:bg-neutral-800 border border-rose-200 dark:border-neutral-700 rounded-xl text-gray-900 dark:text-white focus:ring-2 focus:ring-rose-500 focus:border-transparent outline-none"
        />
      </div>

      <div className="bg-white dark:bg-neutral-900 rounded-2xl border border-rose-100 dark:border-neutral-800 p-5 shadow-sm mb-6">
        <h3 className="text-sm font-bold text-gray-900 dark:text-white mb-4">Time Slots</h3>
        <div className="grid grid-cols-2 sm:grid-cols-4 gap-3">
          {timeSlots.map((slot) => (
            <button
              key={slot}
              onClick={() => toggleSlot(slot)}
              className={`py-3 px-4 rounded-xl text-sm font-medium transition-all duration-200 border-2 ${
                slots[slot] === 'available'
                  ? 'bg-rose-50 border-rose-500 text-rose-700 dark:bg-orange-950 dark:border-orange-500 dark:text-orange-400'
                  : 'bg-gray-50 border-gray-200 text-gray-400 dark:bg-neutral-800 dark:border-neutral-700 dark:text-neutral-500'
              }`}
            >
              {slots[slot] === 'available' ? '● ' : '○ '}{slot}
            </button>
          ))}
        </div>
        <p className="text-xs text-gray-400 dark:text-neutral-500 mt-3">
          Tap a slot to toggle availability. Orange = available, Gray = blocked.
        </p>
      </div>

      <Button className="w-full gap-2">
        <Save size={18} />
        Save Schedule
      </Button>
    </div>
  );
}

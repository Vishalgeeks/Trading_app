import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import Button from '../components/ui/Button';
import { userService } from '../services/userService';
import {
  ChevronRight,
  Heart,
  Calendar,
  Shield,
  LogOut,
} from 'lucide-react';

export default function Profile() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const [profileUser, setProfileUser] = useState(null);
  const [, setLoading] = useState(true);
  const [, setError] = useState('');

  useEffect(() => {
    let cancelled = false;

    async function loadProfile() {
      setLoading(true);
      setError('');

      const result = await userService.getCurrentUser();
      if (cancelled) return;

      if (!result.error && result.data) {
        setProfileUser(result.data);
      }
      setLoading(false);
    }

    loadProfile();

    return () => {
      cancelled = true;
    };
  }, []);

  const handleLogout = async () => {
    await logout();
    navigate('/');
  };

  const menuItems = [
    { icon: Heart, label: 'Favorites', description: 'View saved designs', to: '/favorites' },
    { icon: Calendar, label: 'My Bookings', description: 'View and manage bookings', to: '/bookings' },
    { icon: Settings, label: 'Settings', description: 'Appearance, notifications, privacy', to: '/settings' },
    { icon: Shield, label: 'Privacy & Security', description: 'Password and data settings', to: '/settings' },
  ];

  return (
    <div className="pb-24">
      <div className="relative bg-gradient-to-b from-rose-500 to-rose-600 dark:from-orange-600 dark:to-orange-700 text-white">
        <div className="px-5 pt-12 pb-8">
          <div className="flex items-center gap-4">
            <div className="w-16 h-16 rounded-full bg-white/20 backdrop-blur-sm flex items-center justify-center text-3xl border-2 border-white/30">
              👤
            </div>
            <div>
              <h1 className="text-xl font-bold">{profileUser?.name || user?.name || 'Guest User'}</h1>
              <p className="text-sm text-white/80">{profileUser?.email || user?.email || 'guest@example.com'}</p>
            </div>
          </div>
        </div>
        <div className="absolute -bottom-6 left-0 right-0">
          <div className="w-16 h-4 bg-gray-50 dark:bg-black mx-auto rounded-b-2xl"></div>
        </div>
      </div>

      <div className="px-5 pt-8">
        <div className="space-y-2 mb-8">
          {menuItems.map((item) => (
            <button
              key={item.label}
              onClick={() => navigate(item.to)}
              className="w-full flex items-center gap-4 p-4 bg-white dark:bg-neutral-900 rounded-2xl border border-rose-100 dark:border-neutral-800 shadow-sm hover:shadow-md transition-shadow"
            >
              <div className="w-10 h-10 rounded-xl bg-rose-50 dark:bg-orange-950 flex items-center justify-center text-rose-500 dark:text-orange-400 shrink-0">
                <item.icon size={20} />
              </div>
              <div className="flex-1 text-left">
                <p className="text-sm font-semibold text-gray-900 dark:text-white">{item.label}</p>
                <p className="text-xs text-gray-500 dark:text-neutral-400">{item.description}</p>
              </div>
              <ChevronRight size={18} className="text-gray-400 dark:text-neutral-500" />
            </button>
          ))}
        </div>

        <Button variant="outline" className="w-full" onClick={handleLogout}>
          <LogOut size={18} />
          Log Out
        </Button>
      </div>
    </div>
  );
}

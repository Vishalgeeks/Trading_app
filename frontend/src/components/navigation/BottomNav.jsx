import { NavLink } from 'react-router-dom';
import { Home, Compass, Calendar, User } from 'lucide-react';

const navItems = [
  { to: '/', icon: Home, label: 'Home' },
  { to: '/browse', icon: Compass, label: 'Browse' },
  { to: '/bookings', icon: Calendar, label: 'Bookings' },
  { to: '/profile', icon: User, label: 'Profile' },
];

export default function BottomNav() {
  return (
    <nav className="fixed bottom-0 left-0 right-0 bg-white border-t border-rose-100 dark:bg-neutral-900 dark:border-neutral-800 z-50">
      <div className="max-w-lg mx-auto flex items-center justify-around py-2 pb-safe">
        {navItems.map((item) => (
          <NavLink
            key={item.to}
            to={item.to}
            className={({ isActive }) =>
              `flex flex-col items-center gap-0.5 py-1 px-3 rounded-lg transition-colors duration-200 ${
                isActive
                  ? 'text-rose-500 dark:text-orange-500'
                  : 'text-gray-400 hover:text-rose-500 dark:text-neutral-500 dark:hover:text-orange-400'
              }`
            }
          >
            {({ isActive }) => (
              <>
                <item.icon size={22} strokeWidth={isActive ? 2.5 : 2} />
                <span className="text-[10px] font-medium">{item.label}</span>
              </>
            )}
          </NavLink>
        ))}
      </div>
    </nav>
  );
}

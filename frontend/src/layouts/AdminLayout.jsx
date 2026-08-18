import { NavLink, Outlet } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import {
  LayoutDashboard,
  CalendarDays,
  Image,
  Settings2,
  FolderTree,
  LogOut,
} from 'lucide-react';

const adminNavItems = [
  { to: '/admin', icon: LayoutDashboard, label: 'Dashboard', end: true },
  { to: '/admin/bookings', icon: CalendarDays, label: 'Bookings' },
  { to: '/admin/designs', icon: Image, label: 'Designs' },
  { to: '/admin/categories', icon: FolderTree, label: 'Categories' },
  { to: '/admin/availability', icon: Settings2, label: 'Schedule' },
];

export default function AdminLayout() {
  const { user } = useAuth();
  

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-black">
      <div className="flex">
        <aside className="hidden md:flex flex-col w-64 bg-white dark:bg-neutral-900 border-r border-rose-100 dark:border-neutral-800 min-h-screen sticky top-0">
          <div className="p-6">
            <h2 className="text-xl font-bold text-rose-500 dark:text-orange-500">Henna Booking</h2>
            <p className="text-xs text-gray-500 dark:text-neutral-400 mt-1">Admin Panel</p>
          </div>
          <nav className="flex-1 px-4 space-y-1">
            {adminNavItems.map((item) => (
              <NavLink
                key={item.to}
                to={item.to}
                end={item.end}
                className={({ isActive: navActive }) =>
                  `flex items-center gap-3 px-4 py-3 rounded-xl text-sm font-medium transition-colors duration-200 ${
                    navActive
                      ? 'bg-rose-50 text-rose-600 dark:bg-orange-950 dark:text-orange-400'
                      : 'text-gray-600 hover:bg-rose-50 hover:text-rose-600 dark:text-neutral-400 dark:hover:bg-orange-950 dark:hover:text-orange-400'
                  }`
                }
              >
                <item.icon size={18} />
                {item.label}
              </NavLink>
            ))}
          </nav>
          <div className="p-4 border-t border-rose-100 dark:border-neutral-800">
            <div className="flex items-center gap-3 px-4 py-2 mb-2">
              <div className="w-8 h-8 rounded-full bg-rose-100 dark:bg-orange-900 flex items-center justify-center text-rose-600 dark:text-orange-400 font-semibold text-sm">
                {user?.name?.charAt(0) || 'A'}
              </div>
              <div>
                <p className="text-sm font-medium text-gray-900 dark:text-white">{user?.name}</p>
                <p className="text-xs text-gray-500 dark:text-neutral-400">{user?.email}</p>
              </div>
            </div>
            <NavLink
              to="/"
              className={({ isActive: navActive }) =>
                `flex items-center gap-3 px-4 py-3 text-sm font-medium transition-colors ${
                  navActive
                    ? 'text-rose-600 dark:text-orange-400'
                    : 'text-gray-600 hover:text-rose-600 dark:text-neutral-400 dark:hover:text-orange-400'
                }`
              }
            >
              <LogOut size={18} />
              Logout
            </NavLink>
          </div>
        </aside>

        <main className="flex-1 md:ml-0 pb-20 md:pb-0">
          <div className="md:hidden bg-white dark:bg-neutral-900 border-b border-rose-100 dark:border-neutral-800 px-4 py-3 flex items-center justify-between sticky top-0 z-40">
            <h2 className="text-lg font-bold text-rose-500 dark:text-orange-500">Henna Booking</h2>
          </div>
          <div className="p-4 md:p-8 max-w-5xl mx-auto">
            <Outlet />
          </div>
        </main>
      </div>

      <nav className="md:hidden fixed bottom-0 left-0 right-0 bg-white dark:bg-neutral-900 border-t border-rose-100 dark:border-neutral-800 z-50">
        <div className="flex items-center justify-around py-2">
          {adminNavItems.map((item) => (
            <NavLink
              key={item.to}
              to={item.to}
              end={item.end}
              className={({ isActive: navActive }) =>
                `flex flex-col items-center gap-0.5 py-1 px-2 rounded-lg transition-colors duration-200 ${
                  navActive
                    ? 'text-rose-500 dark:text-orange-500'
                    : 'text-gray-400 dark:text-neutral-500'
                }`
              }
            >
              {({ isActive: navActive }) => (
                <>
                  <item.icon size={20} strokeWidth={navActive ? 2.5 : 2} />
                  <span className="text-[10px] font-medium">{item.label}</span>
                </>
              )}
            </NavLink>
          ))}
        </div>
      </nav>
    </div>
  );
}

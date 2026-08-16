import { Outlet, useLocation } from 'react-router-dom';
import BottomNav from '../components/navigation/BottomNav';

export default function ClientLayout() {
  const location = useLocation();
  const hideBottomNav = ['/login', '/register'].includes(location.pathname);

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-black">
      <main className={hideBottomNav ? '' : 'pb-20'}>
        <Outlet />
      </main>
      {!hideBottomNav && <BottomNav />}
    </div>
  );
}

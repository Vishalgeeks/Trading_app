import { useNavigate } from 'react-router-dom';
import { useTheme } from '../context/ThemeContext';
import Button from '../components/ui/Button';
import ThemeSelector from '../components/settings/ThemeSelector';
import SettingsRow from '../components/settings/SettingsRow';
import {
  Bell,
  Globe,
  CreditCard,
  Shield,
  HelpCircle,
  LogOut,
  Palette,
} from 'lucide-react';

export default function Settings() {
  const { theme, toggleTheme } = useTheme();
  const navigate = useNavigate();

  return (
    <div className="pb-24">
      <div className="px-5 pt-6">
        <h1 className="text-2xl font-bold text-gray-900 dark:text-white mb-1">Settings</h1>
        <p className="text-sm text-gray-500 dark:text-neutral-400 mb-6">
          Manage your app preferences
        </p>

        <div className="space-y-6">
          <div className="bg-white dark:bg-neutral-900 rounded-2xl border border-rose-100 dark:border-neutral-800 p-5 shadow-sm">
            <div className="flex items-center gap-3 mb-4">
              <div className="w-8 h-8 rounded-lg bg-rose-50 dark:bg-orange-950 flex items-center justify-center text-rose-500 dark:text-orange-400">
                <Palette size={18} />
              </div>
              <h3 className="text-sm font-bold text-gray-900 dark:text-white uppercase tracking-wide">
                Appearance
              </h3>
            </div>
            <ThemeSelector theme={theme} onSelect={(t) => t !== theme && toggleTheme()} />
            <p className="text-xs text-gray-500 dark:text-neutral-400 mt-3">
              Choose your preferred theme
            </p>
          </div>

          <div className="bg-white dark:bg-neutral-900 rounded-2xl border border-rose-100 dark:border-neutral-800 shadow-sm overflow-hidden">
            <div className="px-5 py-3 border-b border-rose-100 dark:border-neutral-800">
              <h3 className="text-sm font-bold text-gray-900 dark:text-white uppercase tracking-wide">
                Preferences
              </h3>
            </div>
            <SettingsRow
              icon={Bell}
              label="Notifications"
              description="Manage push and email notifications"
              onClick={() => {}}
            />
            <div className="mx-4 border-t border-rose-100 dark:border-neutral-800" />
            <SettingsRow
              icon={Globe}
              label="Language"
              description="English (US)"
              trailing={<span className="text-xs text-gray-400 dark:text-neutral-500">EN ▾</span>}
              onClick={() => {}}
            />
            <div className="mx-4 border-t border-rose-100 dark:border-neutral-800" />
            <SettingsRow
              icon={CreditCard}
              label="Payment Methods"
              description="Manage payment options"
              onClick={() => {}}
            />
          </div>

          <div className="bg-white dark:bg-neutral-900 rounded-2xl border border-rose-100 dark:border-neutral-800 shadow-sm overflow-hidden">
            <div className="px-5 py-3 border-b border-rose-100 dark:border-neutral-800">
              <h3 className="text-sm font-bold text-gray-900 dark:text-white uppercase tracking-wide">
                Support
              </h3>
            </div>
            <SettingsRow
              icon={Shield}
              label="Privacy & Security"
              description="Password and data settings"
              onClick={() => {}}
            />
            <div className="mx-4 border-t border-rose-100 dark:border-neutral-800" />
            <SettingsRow
              icon={HelpCircle}
              label="Help & Support"
              description="FAQ, contact us"
              onClick={() => {}}
            />
          </div>

          <Button variant="outline" className="w-full" onClick={() => navigate('/profile')}>
            <LogOut size={18} />
            Log Out
          </Button>
        </div>
      </div>
    </div>
  );
}

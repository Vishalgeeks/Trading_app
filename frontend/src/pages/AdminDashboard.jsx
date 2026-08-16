import { useState, useEffect } from 'react';
import { api } from '../services/api';
import PageContainer from '../components/PageContainer';

export default function AdminDashboard() {
  const [stats, setStats] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    fetchStats();
  }, []);

  const fetchStats = async () => {
    const token = localStorage.getItem('token');
    if (!token) {
      setError('Please login to view admin dashboard');
      setLoading(false);
      return;
    }

    const result = await api.get('/admin/bookings/stats');
    if (result.error) {
      setError(result.message);
    } else {
      setStats(result.data);
    }
    setLoading(false);
  };

  if (loading) {
    return (
      <PageContainer title="Admin Dashboard">
        <div className="text-center py-12">
          <p className="text-gray-500">Loading dashboard...</p>
        </div>
      </PageContainer>
    );
  }

  if (error) {
    return (
      <PageContainer title="Admin Dashboard">
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
          {error}
        </div>
      </PageContainer>
    );
  }

  const statCards = [
    { label: 'Total Bookings', value: stats?.total || 0, color: 'bg-blue-500' },
    { label: 'Pending', value: stats?.pending || 0, color: 'bg-yellow-500' },
    { label: 'Confirmed', value: stats?.confirmed || 0, color: 'bg-green-500' },
    { label: 'Completed', value: stats?.completed || 0, color: 'bg-purple-500' },
    { label: 'Cancelled', value: stats?.cancelled || 0, color: 'bg-red-500' },
    { label: 'Upcoming', value: stats?.upcoming || 0, color: 'bg-rose-500' },
  ];

  return (
    <PageContainer title="Admin Dashboard">
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
        {statCards.map((card) => (
          <div key={card.label} className="card">
            <div className={`${card.color} text-white rounded-lg p-4 mb-3`}>
              <p className="text-3xl font-bold">{card.value}</p>
            </div>
            <p className="text-gray-600 font-medium">{card.label}</p>
          </div>
        ))}
      </div>
    </PageContainer>
  );
}

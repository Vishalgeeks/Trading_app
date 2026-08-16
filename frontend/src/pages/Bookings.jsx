import { useState, useEffect } from 'react';
import { api } from '../services/api';
import PageContainer from '../components/PageContainer';

export default function Bookings() {
  const [bookings, setBookings] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    fetchBookings();
  }, []);

  const fetchBookings = async () => {
    const token = localStorage.getItem('token');
    if (!token) {
      setError('Please login to view bookings');
      setLoading(false);
      return;
    }

    const result = await api.get('/bookings');
    if (result.error) {
      setError(result.message);
    } else {
      setBookings(result.data?.bookings || []);
    }
    setLoading(false);
  };

  if (loading) {
    return (
      <PageContainer title="My Bookings">
        <div className="text-center py-12">
          <p className="text-gray-500">Loading bookings...</p>
        </div>
      </PageContainer>
    );
  }

  if (error) {
    return (
      <PageContainer title="My Bookings">
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
          {error}
        </div>
      </PageContainer>
    );
  }

  return (
    <PageContainer title="My Bookings">
      <div className="space-y-4">
        {bookings.map((booking) => (
          <div key={booking.id} className="card">
            <div className="flex justify-between items-start">
              <div>
                <h3 className="text-lg font-semibold">{booking.design_name}</h3>
                <p className="text-gray-600">
                  {booking.booking_date} • {booking.start_time} - {booking.end_time}
                </p>
                <span className={`inline-block mt-2 px-3 py-1 rounded-full text-sm font-medium ${
                  booking.status === 'CONFIRMED' ? 'bg-green-100 text-green-800' :
                  booking.status === 'PENDING' ? 'bg-yellow-100 text-yellow-800' :
                  booking.status === 'CANCELLED' ? 'bg-red-100 text-red-800' :
                  'bg-gray-100 text-gray-800'
                }`}>
                  {booking.status}
                </span>
              </div>
            </div>
          </div>
        ))}
        {bookings.length === 0 && (
          <p className="text-center text-gray-500 mt-8">No bookings yet.</p>
        )}
      </div>
    </PageContainer>
  );
}

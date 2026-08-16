import { useState, useEffect } from 'react';
import { api } from '../services/api';
import PageContainer from '../components/PageContainer';

export default function Designs() {
  const [designs, setDesigns] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    fetchDesigns();
  }, []);

  const fetchDesigns = async () => {
    const result = await api.get('/designs');
    if (result.error) {
      setError(result.message);
    } else {
      setDesigns(result.data?.designs || []);
    }
    setLoading(false);
  };

  if (loading) {
    return (
      <PageContainer title="Designs">
        <div className="text-center py-12">
          <p className="text-gray-500">Loading designs...</p>
        </div>
      </PageContainer>
    );
  }

  if (error) {
    return (
      <PageContainer title="Designs">
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
          {error}
        </div>
      </PageContainer>
    );
  }

  return (
    <PageContainer title="Designs">
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
        {designs.map((design) => (
          <div key={design.id} className="card hover:shadow-md transition-shadow">
            <h3 className="text-lg font-semibold mb-2">{design.name}</h3>
            <p className="text-gray-600 text-sm mb-4">{design.description || 'Beautiful henna design'}</p>
            <div className="flex justify-between items-center">
              <span className="text-sm text-gray-500">Category: {design.category_name || 'General'}</span>
            </div>
          </div>
        ))}
      </div>
      {designs.length === 0 && (
        <p className="text-center text-gray-500 mt-8">No designs available yet.</p>
      )}
    </PageContainer>
  );
}

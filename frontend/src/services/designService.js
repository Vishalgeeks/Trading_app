import { api } from './api';

function formatDuration(minutes) {
  if (!minutes || minutes <= 0) return '';
  const hours = Math.floor(minutes / 60);
  const mins = minutes % 60;
  if (hours > 0 && mins > 0) {
    return `${hours}h ${mins}m`;
  }
  if (hours > 0) {
    return `${hours}h`;
  }
  return `${mins}m`;
}

function mapDesign(raw) {
  if (!raw) return null;
  return {
    id: raw.id,
    name: raw.name,
    slug: raw.slug,
    description: raw.description || '',
    image_url: raw.image_url || null,
    price: raw.price || '',
    duration_minutes: raw.duration_minutes || 0,
    duration: formatDuration(raw.duration_minutes),
    is_active: raw.is_active,
    category: raw.category || null,
    category_id: raw.category?.id || null,
    category_name: raw.category?.name || 'Uncategorized',
    created_at: raw.created_at,
    updated_at: raw.updated_at,
  };
}

export const designService = {
  async listDesigns({ category_id, q, page = 1, limit = 20 } = {}) {
    const params = new URLSearchParams();
    if (category_id) params.set('category_id', category_id);
    if (q) params.set('q', q);
    if (page) params.set('page', String(page));
    if (limit) params.set('limit', String(limit));

    const query = params.toString();
    const result = await api.get(`/designs${query ? `?${query}` : ''}`);
    if (result.error) {
      return result;
    }
    const designs = (result.data?.designs || []).map(mapDesign);
    return { error: false, data: designs, count: result.data?.count || designs.length };
  },

  async getDesign(id) {
    const result = await api.get(`/designs/${id}`);
    if (result.error) {
      return result;
    }
    return { error: false, data: mapDesign(result.data) };
  },

  async searchDesigns(query) {
    if (!query || query.trim() === '') {
      return { error: false, data: [] };
    }
    const result = await api.get(`/designs/search?q=${encodeURIComponent(query.trim())}`);
    if (result.error) {
      return result;
    }
    const designs = (result.data?.designs || []).map(mapDesign);
    return { error: false, data: designs, count: result.data?.count || designs.length };
  },

  async listDesignsByCategory(categoryId) {
    const result = await api.get(`/designs/category/${categoryId}`);
    if (result.error) {
      return result;
    }
    const designs = (result.data?.designs || []).map(mapDesign);
    return { error: false, data: designs, count: result.data?.count || designs.length };
  },
};

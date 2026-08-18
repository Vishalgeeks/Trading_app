import { api } from './api';

const UPLOAD_BASE_URL = import.meta.env.VITE_UPLOAD_BASE_URL || 'http://localhost:8080/uploads';

function resolveImageUrl(image_url) {
  if (!image_url) return null;
  if (image_url.startsWith('http://') || image_url.startsWith('https://')) {
    return image_url;
  }
  if (image_url.startsWith('/uploads/')) {
    return `${UPLOAD_BASE_URL}${image_url.replace('/uploads', '')}`;
  }
  return image_url;
}

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
    image_url: resolveImageUrl(raw.image_url),
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

function toFormData(data, file) {
  const fd = new FormData();
  if (data.name) fd.append('name', data.name);
  if (data.slug) fd.append('slug', data.slug);
  if (data.description !== undefined) fd.append('description', data.description || '');
  if (data.price) fd.append('price', data.price);
  if (data.duration_minutes) fd.append('duration_minutes', String(data.duration_minutes));
  if (data.category_id) fd.append('category_id', data.category_id);
  if (file) {
    fd.append('image', file);
  } else if (data.image_url) {
    fd.append('image_url', data.image_url);
  }
  return fd;
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

  // Admin design management
  async listAdminDesigns({ category_id, q, page = 1, limit = 50 } = {}) {
    const params = new URLSearchParams();
    if (category_id) params.set('category_id', category_id);
    if (q) params.set('q', q);
    if (page) params.set('page', String(page));
    if (limit) params.set('limit', String(limit));

    const query = params.toString();
    const result = await api.get(`/admin/designs${query ? `?${query}` : ''}`);
    if (result.error) {
      return result;
    }
    const designs = (result.data?.designs || []).map(mapDesign);
    return { error: false, data: designs, count: result.data?.count || designs.length };
  },

  async getAdminDesign(id) {
    const result = await api.get(`/admin/designs/${encodeURIComponent(id)}`);
    if (result.error) {
      return result;
    }
    return { error: false, data: mapDesign(result.data) };
  },

  async createDesign(data, file) {
    const fd = toFormData(data, file);
    const result = await api.post('/admin/designs', fd);
    if (result.error) {
      return result;
    }
    return { error: false, data: mapDesign(result.data) };
  },

  async updateDesign(id, data, file) {
    const fd = toFormData(data, file);
    const result = await api.patch(`/admin/designs/${encodeURIComponent(id)}`, fd);
    if (result.error) {
      return result;
    }
    return { error: false, data: mapDesign(result.data) };
  },

  async deleteDesign(id) {
    const result = await api.delete(`/admin/designs/${encodeURIComponent(id)}`);
    if (result.error) {
      return result;
    }
    return { error: false };
  },
};

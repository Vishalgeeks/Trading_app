import { api } from './api';

export const categoryService = {
  async listCategories() {
    const result = await api.get('/categories');
    if (result.error) {
      return result;
    }
    return { error: false, data: result.data?.categories || [] };
  },

  async getCategory(id) {
    const result = await api.get(`/categories/${id}`);
    if (result.error) {
      return result;
    }
    return { error: false, data: result.data };
  },

  // Admin category management
  async listAdminCategories() {
    const result = await api.get('/admin/categories');
    if (result.error) {
      return result;
    }
    return { error: false, data: result.data?.categories || [], count: result.data?.count || 0 };
  },

  async createCategory(data) {
    const result = await api.post('/admin/categories', data);
    if (result.error) {
      return result;
    }
    return { error: false, data: result.data };
  },

  async updateCategory(id, data) {
    const result = await api.patch(`/admin/categories/${encodeURIComponent(id)}`, data);
    if (result.error) {
      return result;
    }
    return { error: false, data: result.data };
  },

  async deleteCategory(id) {
    const result = await api.delete(`/admin/categories/${encodeURIComponent(id)}`);
    if (result.error) {
      return result;
    }
    return { error: false };
  },
};

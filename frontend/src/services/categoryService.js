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
};

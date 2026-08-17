import { api } from './api';

export const favoriteService = {
  async listFavorites() {
    const result = await api.get('/favorites');
    if (result.error) {
      return result;
    }
    const favorites = result.data?.favorites || [];
    return { error: false, data: favorites, count: result.data?.count || favorites.length };
  },

  async addFavorite(designId) {
    const result = await api.post('/favorites', { design_id: designId });
    if (result.error) {
      return result;
    }
    return { error: false, data: result.data };
  },

  async removeFavorite(designId) {
    const result = await api.delete('/favorites/delete', { design_id: designId });
    if (result.error) {
      return result;
    }
    return { error: false, data: result.data };
  },
};

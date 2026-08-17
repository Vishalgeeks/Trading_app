import { api } from './api';

export const userService = {
  async getCurrentUser() {
    const result = await api.get('/auth/me');
    if (result.error) {
      return result;
    }
    return { error: false, data: result.data };
  },
};

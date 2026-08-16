import { api } from './api';

export const authService = {
  async login(email, password) {
    const result = await api.post('/auth/login', { email, password });
    if (result.error) {
      return result;
    }
    return result;
  },

  async register(name, email, phone, password) {
    const result = await api.post('/auth/register', { name, email, phone, password });
    if (result.error) {
      return result;
    }
    return result;
  },

  async logout() {
    const result = await api.post('/auth/logout');
    if (result.error) {
      return result;
    }
    return result;
  },

  async getCurrentUser() {
    const result = await api.get('/auth/me');
    if (result.error) {
      return result;
    }
    return result;
  },
};

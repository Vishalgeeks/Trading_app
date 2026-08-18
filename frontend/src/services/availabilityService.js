import { api } from './api';

export const availabilityService = {
  async getAvailableSlots(date, designId) {
    const result = await api.get(`/availability/slots?design_id=${encodeURIComponent(designId)}&date=${encodeURIComponent(date)}`);
    if (result.error) {
      return result;
    }
    const slots = result.data?.slots || [];
    return { error: false, data: slots, count: result.data?.count || slots.length };
  },

  // Admin availability endpoints
  async getAdminAvailability(activeOnly = true) {
    const result = await api.get(`/admin/availability?active=${activeOnly}`);
    if (result.error) {
      return result;
    }
    return { error: false, data: result.data?.availabilities || [], count: result.data?.count || 0 };
  },

  async createAvailability(data) {
    const result = await api.post('/admin/availability', data);
    if (result.error) {
      return result;
    }
    return { error: false, data: result.data };
  },

  async getAvailability(id) {
    const result = await api.get(`/admin/availability?id=${encodeURIComponent(id)}`);
    if (result.error) {
      return result;
    }
    return { error: false, data: result.data };
  },

  async updateAvailability(id, data) {
    const result = await api.patch(`/admin/availability?id=${encodeURIComponent(id)}`, data);
    if (result.error) {
      return result;
    }
    return { error: false, data: result.data };
  },

  async deleteAvailability(id) {
    const result = await api.delete(`/admin/availability?id=${encodeURIComponent(id)}`);
    if (result.error) {
      return result;
    }
    return { error: false };
  },
};

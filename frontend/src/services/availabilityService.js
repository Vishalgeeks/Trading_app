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
};

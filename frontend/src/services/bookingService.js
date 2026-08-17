import { api } from './api';

export const bookingService = {
  async createBooking({ designId, bookingDate, startTime, notes = '' }) {
    const result = await api.post('/bookings', {
      design_id: designId,
      booking_date: bookingDate,
      start_time: startTime,
      notes,
    });
    if (result.error) {
      return result;
    }
    return { error: false, data: result.data };
  },
};

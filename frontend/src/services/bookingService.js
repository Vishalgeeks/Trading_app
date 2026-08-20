import { api } from './api';

function mapBooking(raw) {
  if (!raw) return null;
  return {
    id: raw.id,
    designId: raw.design_id,
    designName: raw.design_name,
    userName: raw.user_name,
    userEmail: raw.user_email,
    userPhone: raw.user_phone,
    bookingDate: raw.booking_date,
    booking_date: raw.booking_date,
    startTime: raw.start_time,
    start_time: raw.start_time,
    endTime: raw.end_time,
    end_time: raw.end_time,
    status: raw.status,
    notes: raw.notes,
    createdAt: raw.created_at,
    created_at: raw.created_at,
    updatedAt: raw.updated_at,
    updated_at: raw.updated_at,
  };
}

function mapAdminBooking(raw) {
  if (!raw) return null;
  return {
    id: raw.id,
    userId: raw.user_id,
    userName: raw.user_name,
    clientName: raw.user_name,
    userEmail: raw.user_email,
    clientEmail: raw.user_email,
    userPhone: raw.user_phone,
    clientPhone: raw.user_phone,
    designId: raw.design_id,
    designName: raw.design_name,
    bookingDate: raw.booking_date,
    booking_date: raw.booking_date,
    startTime: raw.start_time,
    start_time: raw.start_time,
    endTime: raw.end_time,
    end_time: raw.end_time,
    status: raw.status,
    notes: raw.notes,
    createdAt: raw.created_at,
    created_at: raw.created_at,
    updatedAt: raw.updated_at,
    updated_at: raw.updated_at,
  };
}

export const bookingService = {
  // Client booking endpoints
  async createBooking(data) {
    const payload = {
      design_id: data.designId || data.design_id,
      booking_date: data.bookingDate || data.booking_date,
      start_time: data.startTime || data.start_time,
      notes: data.notes || '',
    };
    const result = await api.post('/bookings', payload);
    if (result.error) {
      return result;
    }
    return { error: false, data: mapBooking(result.data) };
  },

  async getMyBookings(params = {}) {
    const query = new URLSearchParams();
    if (params.status) query.append('status', params.status);
    if (params.date) query.append('date', params.date);
    if (params.from) query.append('from', params.from);
    if (params.to) query.append('to', params.to);
    if (params.limit) query.append('limit', params.limit);
    if (params.offset) query.append('offset', params.offset);
    
    const result = await api.get(`/bookings?${query.toString()}`);
    if (result.error) return result;
    const bookings = (result.data?.bookings || []).map(mapBooking);
    return { 
      error: false, 
      data: bookings, 
      total: result.data?.total,
      limit: result.data?.limit,
      offset: result.data?.offset 
    };
  },

  async getMyBooking(id) {
    const result = await api.get(`/bookings?id=${encodeURIComponent(id)}`);
    if (result.error) return result;
    return { error: false, data: mapBooking(result.data) };
  },

  async cancelBooking(id) {
    const result = await api.patch(`/bookings/${encodeURIComponent(id)}/cancel`, {});
    if (result.error) return result;
    return { error: false, data: mapBooking(result.data) };
  },

  async getUpcomingBookings() {
    const result = await api.get('/bookings/upcoming');
    if (result.error) return result;
    const bookings = (result.data?.bookings || []).map(mapBooking);
    return { error: false, data: bookings };
  },

  async getBookingHistory(params = {}) {
    const query = new URLSearchParams();
    if (params.limit) query.append('limit', params.limit);
    if (params.offset) query.append('offset', params.offset);
    
    const result = await api.get(`/bookings/history?${query.toString()}`);
    if (result.error) return result;
    const bookings = (result.data?.bookings || []).map(mapBooking);
    return { 
      error: false, 
      data: bookings, 
      total: result.data?.total,
      limit: result.data?.limit,
      offset: result.data?.offset 
    };
  },

  // Admin booking endpoints
  async getAdminBookings(params = {}) {
    const query = new URLSearchParams();
    if (params.status) query.append('status', params.status);
    if (params.date) query.append('date', params.date);
    if (params.from) query.append('from', params.from);
    if (params.to) query.append('to', params.to);
    if (params.search) query.append('search', params.search);
    if (params.limit) query.append('limit', params.limit);
    if (params.offset) query.append('offset', params.offset);
    
    const result = await api.get(`/admin/bookings?${query.toString()}`);
    if (result.error) return result;
    const bookings = (result.data?.bookings || []).map(mapAdminBooking);
    return { 
      error: false, 
      data: bookings, 
      total: result.data?.total,
      limit: result.data?.limit,
      offset: result.data?.offset 
    };
  },

  async getAdminBooking(id) {
    const result = await api.get(`/admin/bookings?id=${encodeURIComponent(id)}`);
    if (result.error) return result;
    return { error: false, data: mapAdminBooking(result.data) };
  },

  async updateBookingStatus(id, status) {
    const result = await api.patch(`/admin/bookings/${encodeURIComponent(id)}/status`, { status });
    if (result.error) return result;
    return { error: false, data: mapAdminBooking(result.data) };
  },

  async getAdminUpcomingBookings() {
    const result = await api.get('/admin/bookings/upcoming');
    if (result.error) return result;
    const bookings = (result.data?.bookings || []).map(mapAdminBooking);
    return { error: false, data: bookings };
  },

  async getAdminBookingHistory(params = {}) {
    const query = new URLSearchParams();
    if (params.limit) query.append('limit', params.limit);
    if (params.offset) query.append('offset', params.offset);
    
    const result = await api.get(`/admin/bookings/history?${query.toString()}`);
    if (result.error) return result;
    const bookings = (result.data?.bookings || []).map(mapAdminBooking);
    return { 
      error: false, 
      data: bookings, 
      total: result.data?.total,
      limit: result.data?.limit,
      offset: result.data?.offset 
    };
  },

  async getAdminBookingStats() {
    const result = await api.get('/admin/bookings/stats');
    if (result.error) return result;
    return { error: false, data: result.data };
  }
};

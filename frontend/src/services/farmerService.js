import api from './api';

export const farmerService = {
  fetchFarmers: async () => {
    try {
      const response = await api.get('/farmers');
      return response.data;
    } catch (error) {
      console.error('Error fetching farmers:', error);
      throw error;
    }
  },

  getFarmerById: async (id) => {
    try {
      const response = await api.get(`/farmers/${id}`);
      return response.data;
    } catch (error) {
      console.error('Error fetching farmer:', error);
      throw error;
    }
  },

  createFarmer: async (farmerData) => {
    try {
      const response = await api.post('/farmers', farmerData);
      return response.data;
    } catch (error) {
      console.error('Error creating farmer:', error);
      throw error;
    }
  },

  updateFarmer: async (id, farmerData) => {
    try {
      const response = await api.put(`/farmers/${id}`, farmerData);
      return response.data;
    } catch (error) {
      console.error('Error updating farmer:', error);
      throw error;
    }
  },

  deleteFarmer: async (id) => {
    try {
      const response = await api.delete(`/farmers/${id}`);
      return response.data;
    } catch (error) {
      console.error('Error deleting farmer:', error);
      throw error;
    }
  }
};

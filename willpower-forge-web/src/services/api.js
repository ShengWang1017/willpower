import axios from 'axios';
import { TOKEN_KEY } from '../constants/index.js';

const api = axios.create({
  baseURL: '/api/v1'
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem(TOKEN_KEY);
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Goal APIs
export const deleteGoal = (goalId) => api.delete(`/goals/${goalId}`);
export const getDeletedGoals = () => api.get('/goals/recycle-bin');
export const restoreGoal = (goalId) => api.post(`/goals/${goalId}/restore`);
export const permanentDeleteGoal = (goalId) => api.delete(`/goals/${goalId}/permanent`);

export default api;

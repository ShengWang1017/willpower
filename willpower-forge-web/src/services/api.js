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

export default api;

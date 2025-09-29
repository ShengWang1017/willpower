import { defineStore } from 'pinia';
import api from '../services/api';
import { TOKEN_KEY } from '../constants/index.js';

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem(TOKEN_KEY) || '',
    user: null
  }),
  getters: {
    isAuthenticated: (state) => Boolean(state.token)
  },
  actions: {
    setToken(token) {
      this.token = token;
      if (token) {
        localStorage.setItem(TOKEN_KEY, token);
      } else {
        localStorage.removeItem(TOKEN_KEY);
      }
    },
    async register(payload) {
      await api.post('/auth/register', payload);
    },
    async login(credentials) {
      const response = await api.post('/auth/login', credentials);
      const { token, user_id: userId } = response.data.data;
      this.setToken(token);
      this.user = { id: userId, username: credentials.username };
    },
    logout() {
      this.setToken('');
      this.user = null;
    }
  }
});

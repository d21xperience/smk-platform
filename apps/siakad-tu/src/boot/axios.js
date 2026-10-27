import axios from 'axios';
import { useContextStore } from '../stores/contextStore.js';

const api = axios.create({
  baseURL: process.env.API_URL || 'http://localhost:8080',
});

api.interceptors.request.use((config) => {
  const contextStore = useContextStore();
  const ctx = contextStore.currentContext;

  if (ctx.schoolId) config.headers['X-School-Id'] = ctx.schoolId;
  if (ctx.userId) config.headers['X-User-Id'] = ctx.userId;
  
  if (ctx.semesterId) {
    config.headers['X-Academic-Period'] = ctx.semesterId; 
  }

  return config;
});

export default api;

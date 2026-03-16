import axios from 'axios';

function getBaseURL(): string {
  // In dev with Vite proxy
  if (window.location.port === '5173') {
    return '/api/v1';
  }
  // In Wails desktop or production: call backend directly
  return 'http://localhost:3001/api/v1';
}

const api = axios.create({
  baseURL: getBaseURL(),
  timeout: 30000,
  headers: { 'Content-Type': 'application/json' },
});

export default api;

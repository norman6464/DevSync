import axios from 'axios';

const client = axios.create({
  baseURL: '/api/v1',
  withCredentials: true, // httpOnly Cookieを自動送信
});

client.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // 認証チェック（/auth/me）の401はloadUser内で処理するためスキップ
      const url = error.config?.url || '';
      if (!url.includes('/auth/me') && window.location.pathname !== '/login') {
        window.location.href = '/login';
      }
    }
    return Promise.reject(error);
  }
);

export default client;

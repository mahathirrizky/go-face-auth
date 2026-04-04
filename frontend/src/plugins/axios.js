import axios from 'axios';
import { useAuthStore } from '../stores/auth';
import { useWebSocketStore } from '../stores/websocket';

export function setupAxios(router) {
    // Determine API base URL
    const envApiBase = process.env.VITE_API_BASE_URL || 'http://localhost:8080/';
    const isDev = Boolean(typeof import.meta !== 'undefined' && import.meta.env && import.meta.env.DEV);
    const axiosBase = isDev ? '/api' : envApiBase;

    axios.defaults.baseURL = axiosBase;
    axios.defaults.withCredentials = true;

    // Request interceptor
    axios.interceptors.request.use(config => {
        const auth = useAuthStore();
        if (auth && auth.token) {
            config.headers = config.headers || {};
            config.headers.Authorization = `Bearer ${auth.token}`;
        }
        return config;
    }, error => Promise.reject(error));

    // Response interceptor
    axios.interceptors.response.use(response => response, async error => {
        const auth = useAuthStore();
        const webSocketStore = useWebSocketStore();

        if (error && error.response && error.response.status === 401) {
            console.log('401 Unauthorized — clearing auth and redirecting');
            try { auth.clearAuth(); } catch (e) { /* ignore */ }
            if (webSocketStore && webSocketStore.isConnected) {
                try { webSocketStore.closeWebSocket(); } catch (e) { /* ignore */ }
            }

            if (router) {
                try { await router.push('/'); } catch (e) { window.location.href = '/'; }
            } else {
                window.location.href = '/';
            }
        }
        return Promise.reject(error);
    });
}

import { useAuthStore } from '../stores/auth';
import { useWebSocketStore } from '../stores/websocket';

export function setupWebSocket(router, subdomain) {
    const authStore = useAuthStore();
    const webSocketStore = useWebSocketStore();

    // Determine API base URL
    const envApiBase = process.env.VITE_API_BASE_URL || 'http://localhost:8080/';

    // Determine backend host for WebSocket
    let backendUrlForWs;
    try {
        backendUrlForWs = new URL(envApiBase);
    } catch (e) {
        backendUrlForWs = new URL(window.location.origin);
    }

    const wsProtocol = backendUrlForWs.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsHostPort = backendUrlForWs.host; // host:port

    let wsPath = '';
    if (subdomain === 'admin') {
        wsPath = '/ws/dashboard';
    } else if (subdomain === 'superadmin') {
        wsPath = '/ws/superadmin-dashboard';
    }

    const wsUrl = wsPath ? `${wsProtocol}//${wsHostPort}${wsPath}` : '';

    // Initialize WebSocket connection if token exists and a relevant subdomain
    if (authStore.token && wsUrl) {
        console.log(`Initializing WebSocket for ${subdomain} at ${wsUrl}`);
        webSocketStore.initWebSocket(wsUrl);
    }

    // Add a navigation guard to close WebSocket on logout or token invalidation
    router.beforeEach((to, from, next) => {
        if (!authStore.token && webSocketStore.isConnected) {
            console.log('Auth token missing, closing WebSocket connection.');
            try { webSocketStore.closeWebSocket(); } catch (e) { /* ignore */ }
        }
        next();
    });
}

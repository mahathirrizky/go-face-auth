import { createApp } from 'vue';
import { createPinia } from 'pinia';
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate';

import App from './App.vue';
import router from './router'; // Main router
import routeradmin from './router/admin'; // Admin router
import routersuperadmin from './router/superadmin'; // SuperAdmin router
import { getSubdomain } from './utils/subdomain';
import Toast,{POSITION} from "vue-toastification";


/* import the fontawesome core */
import { library } from '@fortawesome/fontawesome-svg-core'

/* import font awesome icon component */
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

/* import specific icons */
import { faTachometerAlt, faBuilding, faReceipt, faChartLine, faBars, faUsers, faCalendarCheck, faCog, faBoxOpen, faBullhorn, faCalendarAlt, faEdit, faTrashAlt, faMapMarkerAlt } from '@fortawesome/free-solid-svg-icons'
import { faEye, faEyeSlash } from '@fortawesome/free-regular-svg-icons'

/* add icons to the library */
library.add(faTachometerAlt, faBuilding, faReceipt, faChartLine, faBars, faUsers, faCalendarCheck, faCog, faBoxOpen, faEye, faEyeSlash, faBullhorn, faCalendarAlt, faEdit, faTrashAlt, faMapMarkerAlt)

import axios from 'axios';
import { useAuthStore } from './stores/auth';
import { useAdminBroadcastStore } from './stores/adminBroadcast.js'; // New import
import { useWebSocketStore } from './stores/websocket'; // New import

const app = createApp(App);
const pinia = createPinia();
pinia.use(piniaPluginPersistedstate);

// Get the current subdomain
const subdomain = getSubdomain();
let selectedRouter;
const apiBaseUrl = process.env.VITE_API_BASE_URL || 'http://localhost:8080/api'; // Fallback for development
axios.defaults.baseURL = apiBaseUrl; // Set Axios base URL

app.use(pinia);
const authStore = useAuthStore();
console.log("main.js: authStore.companyId after pinia init:", authStore.companyId);

// Add a request interceptor to set the Authorization header
axios.interceptors.request.use(config => {
  const authStore = useAuthStore();
  if (authStore.token) {
    config.headers.Authorization = `Bearer ${authStore.token}`;
  }
  return config;
}, error => {
  return Promise.reject(error);
});

// Add a response interceptor to handle 401 Unauthorized errors
axios.interceptors.response.use(response => {
  return response;
}, async error => {
  const authStore = useAuthStore();
  const webSocketStore = useWebSocketStore(); // Get WebSocket store instance
  // Check if the error is a 401 Unauthorized response
  if (error.response && error.response.status === 401) {
    console.log('401 Unauthorized response received. Token might be expired or invalid.');
    // Clear authentication state
    authStore.clearAuth();
    // Close WebSocket connection on 401
    if (webSocketStore.isConnected) {
      console.log('Closing WebSocket due to 401 Unauthorized.');
      webSocketStore.closeWebSocket();
    }
    // Redirect to the login page (or root path which handles redirection to login)
    // Use the selectedRouter to push to the appropriate login path
    if (selectedRouter) { // Ensure router is initialized
      console.log('Redirecting to login page...');
      await selectedRouter.push('/'); // Assuming '/' is the unauthenticated landing page
    } else {
      console.error('Router not initialized for redirection.');
      // Fallback if router is not yet available (should not happen in normal flow)
      window.location.href = '/';
    }
  }
  return Promise.reject(error);
});

console.log('Current subdomain:', subdomain);
// Determine which router to use based on the subdomain
if (subdomain === 'admin') {
    selectedRouter = routeradmin; // Use the admin router
} else if (subdomain === 'superadmin') {
    selectedRouter = routersuperadmin; // Use the superadmin router
} else {
    selectedRouter = router; // Use the main router
}

const options = {
    position: POSITION.TOP_CENTER,
    timeout: 5000,
  };

// Use the selected router
app.use(selectedRouter);
app.use(Toast, options);
app.component('font-awesome-icon', FontAwesomeIcon);

// --- Start WebSocket Integration ---
const webSocketStore = useWebSocketStore();

// Parse apiBaseUrl to get the host and port for WebSocket
const url = new URL(apiBaseUrl);
const wsProtocol = url.protocol === 'https:' ? 'wss:' : 'ws:';
const wsHostPort = url.host; // This will be 'localhost:8080'

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
selectedRouter.beforeEach((to, from, next) => {
    // Check if the user is navigating away from authenticated routes
    // and if the token is no longer present (e.g., after logout or 401 interceptor)
    if (!authStore.token && webSocketStore.isConnected) {
        console.log('Auth token missing, closing WebSocket connection.');
        webSocketStore.closeWebSocket();
    }
    next();
});
// --- End WebSocket Integration ---

app.mount('#app');
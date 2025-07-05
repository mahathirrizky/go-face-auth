import { createApp } from 'vue';
import { createPinia } from 'pinia';
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate';
import './style.css';
import App from './App.vue';
import router from './router'; // Main router
import routeradmin from './router/admin'; // Admin router
import { getSubdomain } from './utils/subdomain';
import Toast,{POSITION} from "vue-toastification";
import "vue-toastification/dist/index.css";

import axios from 'axios';
import { useAuthStore } from './stores/auth';

const app = createApp(App);
const pinia = createPinia();
pinia.use(piniaPluginPersistedstate);

// Get the current subdomain
const subdomain = getSubdomain();
let selectedRouter;
window.base_url = 'http://localhost:8080';
axios.defaults.baseURL = window.base_url; // Set Axios base URL

app.use(pinia);
const authStore = useAuthStore();

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

console.log('Current subdomain:', subdomain);
// Determine which router to use based on the subdomain
if (subdomain === 'admin') {
    selectedRouter = routeradmin; // Use the admin router
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
app.mount('#app');
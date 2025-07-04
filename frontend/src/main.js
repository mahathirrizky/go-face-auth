import { createApp } from 'vue';
import './style.css';
import App from './App.vue';
import router from './router'; // Main router
import routeradmin from './router/admin'; // Admin router
import { getSubdomain } from './utils/subdomain';
import Toast,{POSITION} from "vue-toastification";
import "vue-toastification/dist/index.css";

import axios from 'axios';

const app = createApp(App);

// Get the current subdomain
const subdomain = getSubdomain();
let selectedRouter;
window.base_url = 'https://19db-36-84-48-46.ngrok-free.app';
axios.defaults.baseURL = window.base_url; // Set Axios base URL
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
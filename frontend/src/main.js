import { createApp } from 'vue';
import './style.css';
import App from './App.vue';
import router from './router'; // Main router
import routeradmin from './router/admin'; // Admin router
import { getSubdomain } from './utils/subdomain';
import Toast,{POSITION} from "vue-toastification";
import "vue-toastification/dist/index.css";

const app = createApp(App);

// Get the current subdomain
const subdomain = getSubdomain();
let selectedRouter;
window.base_url = 'https://817c-2001-448a-c140-6fe-5fbc-26aa-8977-bed3.ngrok-free.app';
// Determine which router to use based on the subdomain
if (subdomain === 'owner') {
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

import { createApp } from 'vue';
import { createPinia } from 'pinia';
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate';
import PrimeVue from 'primevue/config';
import Aura from '@primeuix/themes/aura';
import App from './App.vue';
import ToastService from 'primevue/toastservice';
import ConfirmationService from 'primevue/confirmationservice';
import Tooltip from 'primevue/tooltip';

import { getSubdomain } from './utils/subdomain';
import { getSelectedRouter } from './utils/router_selector';
import { setupAxios } from './plugins/axios';
import { setupWebSocket } from './plugins/websocket';

const app = createApp(App);
const pinia = createPinia();
pinia.use(piniaPluginPersistedstate);

app.use(pinia);

app.use(PrimeVue, {
    theme: {
        preset: Aura,
        options: {
            cssLayer: {
                name: 'primevue',
                order: 'theme, base, primevue'
            }
        }
    }
});

// 1. Get Subdomain and Router
const subdomain = getSubdomain();
console.log('Current subdomain:', subdomain);
const selectedRouter = getSelectedRouter(subdomain);

// 2. Setup Plugins (Axios & WebSocket)
// Note: These must run after Pinia is installed because they use stores
setupAxios(selectedRouter);
setupWebSocket(selectedRouter, subdomain);

// 3. Mount App
app.use(selectedRouter);
app.use(ToastService);
app.use(ConfirmationService);
app.directive('tooltip', Tooltip);

app.mount('#app');
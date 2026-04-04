import router from '../router';
import routeradmin from '../router/admin';
import routersuperadmin from '../router/superadmin';

export function getSelectedRouter(subdomain) {
    if (subdomain === 'admin') {
        return routeradmin;
    } else if (subdomain === 'superadmin') {
        return routersuperadmin;
    } else {
        return router;
    }
}

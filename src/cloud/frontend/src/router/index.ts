import { createRouter, createWebHistory } from 'vue-router';
import Home from "@/components/cloud/Home.vue";
import Login from "@/components/cloud/Login.vue";
import axios from "axios";
import {isSecurityEnabled} from "@/components/cloud/Config";
import HubComponent from "@/components/hub/HubHome.vue";
import HubLogin from "@/components/hub/HubLogin.vue";
import HubRegistration from "@/components/hub/HubRegistration.vue";
import HubChangePassword from "@/components/hub/HubChangePassword.vue";
import HubTagManagement from "@/components/hub/HubTagManagement.vue";
import {session} from "@/components/hub/shared";

const routes = [
    {
        path: '/',
        name: 'Home',
        component: Home,
        meta: { requiresAuth: true },
    },
    {
        path: '/login',
        name: 'Login',
        component: Login,
    },
    {
        path: '/hub',
        name: 'HubComponent',
        component: HubComponent,
    },
    {
        path: '/hub/login',
        name: 'HubLogin',
        component: HubLogin,
    },
    {
        path: '/hub/registration',
        name: 'HubRegistration',
        component: HubRegistration,
    },
    {
        path: '/hub/change-password',
        name: 'HubChangePassword',
        component: HubChangePassword,
    },
    {
        path: '/hub/tag-management',
        name: 'HubTagManagement',
        component: HubTagManagement,
    },
];

const router = createRouter({
    history: createWebHistory(process.env.BASE_URL),
    routes,
});

router.beforeEach(async (to, from, next) => {
    if (to.path.startsWith('/hub')) {
        if (to.path == '/hub/login' || to.path == '/hub/registration' || session.isAuthenticated) {
            next();
        } else {
            const isHubSessionValid = await isThereValidHubSessionCookie();
            if (isHubSessionValid) {
                next();
            } else {
                next({ name: 'HubLogin' });
            }
        }
        return;
    }

    // TODO Apply the upper approach to this router.
    // TODO Get rid of "isSecurityEnabled"
    if (isSecurityEnabled && to.matched.some(record => record.meta.requiresAuth) && !(await isThereValidCloudSessionCookie())) {
        next({ name: 'Login' });
    } else {
        next();
    }
});

async function isThereValidHubSessionCookie(): Promise<boolean> {
    try {
        const response = await axios.get('http://localhost:8082/auth-check');
        if (response.status === 200) {
            session.user = response.data.value;
            session.isAuthenticated = true
            return true;
        }
        return false
    } catch (error) {
        return false;
    }
}

async function isThereValidCloudSessionCookie(): Promise<boolean> {
    try {
        // TODO I think the first part of the URL is missing, right?
        await axios.get('/api/check-session', { withCredentials: true });
        return true;
    } catch (error) {
        return false;
    }
}

export default router;
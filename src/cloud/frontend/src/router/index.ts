import { createRouter, createWebHistory } from 'vue-router';
import Home from "@/components/Home.vue";
import Login from "@/components/Login.vue";
import axios from "axios";
import {isSecurityEnabled} from "@/components/Config";
import Hub from "@/components/Hub.vue";
import HubLogin from "@/components/HubLogin.vue";
import HubRegistration from "@/components/HubRegistration.vue";

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
        name: 'Hub',
        component: Hub,
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
];

const router = createRouter({
    history: createWebHistory(process.env.BASE_URL),
    routes,
});

router.beforeEach(async (to, from, next) => {
    if (isSecurityEnabled && to.matched.some(record => record.meta.requiresAuth) && !(await isSessionCookieValid())) {
        next({ name: 'Login' });
    } else {
        next();
    }
});

async function isSessionCookieValid(): Promise<boolean> {
    try {
        await axios.get('/api/check-session', { withCredentials: true });
        return true;
    } catch (error) {
        return false;
    }
}

export default router;

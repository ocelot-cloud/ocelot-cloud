import { createRouter, createWebHistory } from 'vue-router';
import Home from "@/components/cloud/Home.vue";
import Login from "@/components/cloud/Login.vue";
import axios from "axios";
import {isSecurityEnabled} from "@/components/cloud/Config";
import HubComponent from "@/components/hub/HubHome.vue";
import HubLogin from "@/components/hub/HubLogin.vue";
import HubRegistration from "@/components/hub/HubRegistration.vue";
import HubChangePassword from "@/components/hub/HubChangePassword.vue";

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

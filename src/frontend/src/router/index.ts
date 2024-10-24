import { createRouter, createWebHistory } from 'vue-router';
import Home from "@/components/Home.vue";
import Login from "@/components/Login.vue";
import axios from "axios";
import {globalConfig} from "@/components/global_config";

export interface Session {
    user: string;
    isAuthenticated: boolean;
}

export const cloudSession: Session = {
    user: "",
    isAuthenticated: false,
};

const routes = [
    {
        path: '/',
        name: 'Home',
        component: Home,
    },
    {
        path: '/login',
        name: 'Login',
        component: Login,
    },
];

const router = createRouter({
    history: createWebHistory(import.meta.env.VITE_BASE_URL),
    routes,
});

async function isThereValidSession(session: Session, apiUrl: string): Promise<boolean> {
    try {
        const response = await axios.get(apiUrl);
        if (response.status === 200) {
            session.user = response.data.value;
            session.isAuthenticated = true;
            return true;
        }
        return false;
    } catch (error) {
        return false;
    }
}

router.beforeEach(async (to, from, next) => {
    if (to.path === '/login' || cloudSession.isAuthenticated) {
        next();
    } else {
        if (await isThereValidSession(cloudSession, globalConfig.cloudBaseUrl + '/api/check-auth')) {
            next();
        } else {
            next({ name: 'Login' });
        }
    }
});

export default router;

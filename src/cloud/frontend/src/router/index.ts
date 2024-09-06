import { createRouter, createWebHistory } from 'vue-router';
import Home from "@/components/cloud/Home.vue";
import Login from "@/components/cloud/Login.vue";
import axios from "axios";
import HubComponent from "@/components/hub/HubHome.vue";
import HubLogin from "@/components/hub/HubLogin.vue";
import HubRegistration from "@/components/hub/HubRegistration.vue";
import HubChangePassword from "@/components/hub/HubChangePassword.vue";
import HubTagManagement from "@/components/hub/HubTagManagement.vue";
import { cloudBaseUrl } from "@/components/cloud/Config";

export interface Session {
    user: string;
    isAuthenticated: boolean;
}

export const cloudSession: Session = {
    user: "",
    isAuthenticated: false,
};

export const hubSession: Session = {
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
    if (to.path.startsWith('/hub')) {
        if (to.path === '/hub/login' || to.path === '/hub/registration' || hubSession.isAuthenticated) {
            next();
        } else {
            if (await isThereValidSession(hubSession, 'http://localhost:8082/auth-check')) {
                next();
            } else {
                next({ name: 'HubLogin' });
            }
        }
        return;
    }

    if (to.path === '/login' || cloudSession.isAuthenticated) {
        next();
    } else {
        if (await isThereValidSession(cloudSession, cloudBaseUrl + '/api/check-auth')) {
            next();
        } else {
            next({ name: 'Login' });
        }
    }
});

export default router;

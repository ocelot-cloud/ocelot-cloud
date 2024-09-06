import { createRouter, createWebHistory } from 'vue-router';
import Home from "@/components/cloud/Home.vue";
import Login from "@/components/cloud/Login.vue";
import axios from "axios";
import HubComponent from "@/components/hub/HubHome.vue";
import HubLogin from "@/components/hub/HubLogin.vue";
import HubRegistration from "@/components/hub/HubRegistration.vue";
import HubChangePassword from "@/components/hub/HubChangePassword.vue";
import HubTagManagement from "@/components/hub/HubTagManagement.vue";
import {hubSession} from "@/components/hub/shared";
import {backendBaseUrl} from "@/components/cloud/Config";

// TODO Duplication with hubSession
export const cloudSession = {
    isAuthenticated: false,
}

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
})

router.beforeEach(async (to, from, next) => {
    if (to.path.startsWith('/hub')) {
        if (to.path == '/hub/login' || to.path == '/hub/registration' || hubSession.isAuthenticated) {
            next();
        } else {
            if (await isThereValidHubSessionCookie()) {
                next();
            } else {
                next({ name: 'HubLogin' });
            }
        }
        return;
    }

    // TODO Apply the upper approach to this router.
    if (to.path == '/login' || cloudSession.isAuthenticated) {
        next();
    } else {
        if (await isThereValidCloudSessionCookie()) {
            next();
        } else {
            next({ name: 'Login' });
        }
    }
});

async function isThereValidHubSessionCookie(): Promise<boolean> {
    try {
        const response = await axios.get('http://localhost:8082/auth-check');
        if (response.status === 200) {
            hubSession.user = response.data.value;
            hubSession.isAuthenticated = true
            return true;
        }
        return false
    } catch (error) {
        return false;
    }
}

// TODO Duplication wit hisThereValidHubSessionCookie
// TODO Here should a request to the backend happen to check if a cookie is valid or not. Like I already did in the hub.
async function isThereValidCloudSessionCookie(): Promise<boolean> {
    try {
        const response = await axios.get(backendBaseUrl + '/api/check-auth');
        if (response.status === 200) {
            // TODO cloudSession.user = response.data.value;
            cloudSession.isAuthenticated = true
            return true;
        }
        return false
    } catch (error) {
        return false;
    }
}


export default router;
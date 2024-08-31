import {BackendClient} from "@/components/cloud/Shared";
import {BackendClientImpl} from "@/components/cloud/BackendClientImpl";

export const scheme = "http"
export const baseDomain = "localhost"

let PROFILE: string
export let backendBaseUrl: string
export let stackUrl: string
export let backendClient: BackendClient

enum PROFILE_VALUES {
    TEST = "TEST",
    PROD = "PROD",
}

export function initializeGlobalConfig() {
    PROFILE = import.meta.env.VITE_APP_PROFILE || PROFILE_VALUES.PROD
    if (PROFILE === PROFILE_VALUES.TEST) {
        backendBaseUrl = 'http://localhost:8080'
        backendClient = new BackendClientImpl()
    } else if (PROFILE === PROFILE_VALUES.PROD) {
        backendBaseUrl = scheme + '://ocelot-cloud.' + baseDomain
        backendClient = new BackendClientImpl()
    } else {
        throw new Error("error, provided VITE_APP_PROFILE is not known: " + PROFILE);
    }
    stackUrl = backendBaseUrl + '/api/stacks/';
}

initializeGlobalConfig()
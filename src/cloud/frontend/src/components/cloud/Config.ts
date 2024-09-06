import {BackendClient} from "@/components/cloud/BackendClient";

export const scheme = "http"
export const baseDomain = "localhost"

let PROFILE: string
export let cloudBaseUrl: string
export let stackUrl: string
export let backendClient: BackendClient

enum PROFILE_VALUES {
    TEST = "TEST",
    PROD = "PROD",
}

export function initializeGlobalConfig() {
    PROFILE = import.meta.env.VITE_APP_PROFILE || PROFILE_VALUES.PROD
    if (PROFILE === PROFILE_VALUES.TEST) {
        cloudBaseUrl = 'http://localhost:8080'
        backendClient = new BackendClient()
    } else if (PROFILE === PROFILE_VALUES.PROD) {
        cloudBaseUrl = scheme + '://ocelot-cloud.' + baseDomain
        backendClient = new BackendClient()
    } else {
        throw new Error("error, provided VITE_APP_PROFILE is not known: " + PROFILE);
    }
    stackUrl = cloudBaseUrl + '/api/stacks/';
}

initializeGlobalConfig()

export class Stack {
    name: string;
    state: string;
    urlPath: string;

    constructor(name: string, state: string, urlPath: string) {
        this.name = name;
        this.state = state;
        this.urlPath = urlPath;
    }
}
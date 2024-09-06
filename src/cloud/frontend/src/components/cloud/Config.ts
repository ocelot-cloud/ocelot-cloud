import {BackendClient} from "@/components/cloud/BackendClient";
import axios, {AxiosResponse} from "axios";

export const scheme = "http"
export const baseDomain = "localhost"

interface GlobalConfig {
    cloudBaseUrl: string;
    stackUrl: string;
    backendClient: BackendClient;
}

function getGlobalConfig(): GlobalConfig {
    const PROFILE = import.meta.env.VITE_APP_PROFILE || PROFILE_VALUES.PROD;
    let cloudBaseUrl: string;
    let backendClient: BackendClient;

    if (PROFILE === PROFILE_VALUES.TEST) {
        cloudBaseUrl = 'http://localhost:8080';
        backendClient = new BackendClient();
    } else if (PROFILE === PROFILE_VALUES.PROD) {
        cloudBaseUrl = 'http://ocelot-cloud.' + 'localhost';
        backendClient = new BackendClient();
    } else {
        throw new Error("error, provided VITE_APP_PROFILE is not known: " + PROFILE);
    }

    const stackUrl = cloudBaseUrl + '/api/stacks/';

    return {
        cloudBaseUrl,
        stackUrl,
        backendClient
    };
}

export const config = getGlobalConfig()

enum PROFILE_VALUES {
    TEST = "TEST",
    PROD = "PROD",
}

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
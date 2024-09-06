export const scheme = "http"
export const baseDomain = "localhost"

interface GlobalConfig {
    cloudBaseUrl: string;
    stackUrl: string;
}

function getGlobalConfig(): GlobalConfig {
    const PROFILE = import.meta.env.VITE_APP_PROFILE || PROFILE_VALUES.PROD;
    let cloudBaseUrl: string;

    if (PROFILE === PROFILE_VALUES.TEST) {
        cloudBaseUrl = 'http://localhost:8080';
    } else if (PROFILE === PROFILE_VALUES.PROD) {
        cloudBaseUrl = 'http://ocelot-cloud.' + 'localhost';
    } else {
        throw new Error("error, provided VITE_APP_PROFILE is not known: " + PROFILE);
    }

    const stackUrl = cloudBaseUrl + '/api/stacks/';

    return {
        cloudBaseUrl,
        stackUrl,
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
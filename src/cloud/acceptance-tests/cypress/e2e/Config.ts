
export const scheme = "http"

export enum PROFILE_VALUES {
    TEST = "TEST",
    PROD = "PROD",
}

export let rootDomain: string
export const PROFILE = Cypress.env('PROFILE') || PROFILE_VALUES.PROD;
export let ocelotUrl: string

function initializeGlobalConfig() {
    rootDomain = "localhost"
    if (PROFILE === PROFILE_VALUES.PROD) {
        ocelotUrl = "http://ocelot-cloud." + rootDomain
    } else if (PROFILE === PROFILE_VALUES.TEST){
        ocelotUrl = "http://localhost:8081"
    } else {

    }
}

initializeGlobalConfig()

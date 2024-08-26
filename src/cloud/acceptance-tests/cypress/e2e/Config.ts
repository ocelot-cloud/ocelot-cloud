
export const scheme = "http"

export enum PROFILE_VALUES {
    TEST = "TEST",
    PROD = "PROD",
}

// Takes the value from CYPRESS_PROFILE
export let rootDomain: string
export const CYPRESS_PROFILE = Cypress.env('PROFILE') || PROFILE_VALUES.PROD;
export let ocelotUrl: string
export let isSecurityEnabled: boolean

function initializeGlobalConfig() {
    rootDomain = "localhost"
    if (CYPRESS_PROFILE === PROFILE_VALUES.PROD) {
        ocelotUrl = "http://ocelot-cloud." + rootDomain
        isSecurityEnabled = true
    } else if (CYPRESS_PROFILE === PROFILE_VALUES.TEST){
        ocelotUrl = "http://localhost:8081"
        isSecurityEnabled = false // TODO should be enabled in both cases
    } else {

    }
}

initializeGlobalConfig()


export const scheme = "http"

export enum CYPRESS_PROFILE_VALUES {
    separateGuiProfile = "development-setup",
    separateGuiWithBackendMockProfile = "backend-mock",
    productionProfile = "production",
}

// Takes the value from CYPRESS_PROFILE
export let rootDomain: string
export const CYPRESS_PROFILE = Cypress.env('PROFILE') || CYPRESS_PROFILE_VALUES.productionProfile;
export let ocelotUrl: string
export let areRealDockerContainersUsed: boolean
export let isSecurityEnabled: boolean

function initializeGlobalConfig() {
    rootDomain = "localhost"
    if (CYPRESS_PROFILE === CYPRESS_PROFILE_VALUES.productionProfile) {
        ocelotUrl = "http://ocelot-cloud." + rootDomain
        areRealDockerContainersUsed = true
        isSecurityEnabled = true
    } else if (CYPRESS_PROFILE === CYPRESS_PROFILE_VALUES.separateGuiProfile){
        ocelotUrl = "http://localhost:8081"
        areRealDockerContainersUsed = false
        isSecurityEnabled = false
    } else if (CYPRESS_PROFILE === CYPRESS_PROFILE_VALUES.separateGuiWithBackendMockProfile){
        ocelotUrl = "http://localhost:8081"
        areRealDockerContainersUsed = false
        isSecurityEnabled = false
    }
}

initializeGlobalConfig()

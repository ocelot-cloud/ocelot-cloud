import {BackendClient} from "@/components/cloud/Shared";
import {BackendClientImpl} from "@/components/cloud/BackendClientImpl";
import {BackendClientMock} from "@/components/cloud/BackendClientMock";

export const scheme = "http"
export const baseDomain = "localhost"

let VUE_APP_PROFILE: string
export let backendBaseUrl: string
export let stackUrl: string
export let backendClient: BackendClient
export let waitTimeInMillis: number
export let isSecurityEnabled: boolean

enum VUE_APP_PROFILE_VALUES {
    separateGuiProfile = "development-setup",
    separateGuiWithBackendMockProfile = "backend-mock",
    productionProfile = "production",
}

export function initializeGlobalConfig() {
    VUE_APP_PROFILE = process.env.VUE_APP_PROFILE || VUE_APP_PROFILE_VALUES.productionProfile
    if (VUE_APP_PROFILE === VUE_APP_PROFILE_VALUES.separateGuiProfile) {
        backendBaseUrl = 'http://localhost:8080'
        backendClient = new BackendClientImpl()
        isSecurityEnabled = false
    } else if (VUE_APP_PROFILE === VUE_APP_PROFILE_VALUES.separateGuiWithBackendMockProfile) {
        backendBaseUrl = 'http://fake.domain'
        backendClient = new BackendClientMock()
        isSecurityEnabled = false
    } else if (VUE_APP_PROFILE === VUE_APP_PROFILE_VALUES.productionProfile) {
        backendBaseUrl = scheme + '://ocelot-cloud.' + baseDomain
        backendClient = new BackendClientImpl()
        isSecurityEnabled = true
    } else {
        throw new Error("error, provided VUE_APP_PROFILE is not known: " + VUE_APP_PROFILE);
    }
    stackUrl = backendBaseUrl + '/api/stacks/';
    waitTimeInMillis = VUE_APP_PROFILE === VUE_APP_PROFILE_VALUES.separateGuiWithBackendMockProfile ? 100 : 1000;
}

initializeGlobalConfig()
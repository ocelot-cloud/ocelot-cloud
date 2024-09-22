import router from "@/router";
import {AxiosResponse} from "axios";
import {doRequest} from "@/components/shared/requests";
import {globalConfig} from "@/components/shared/global_config";

export function goToHubPage(path: string) {
    router.push('/hub' + path)
}

export function goToCloudPage(path: string) {
    router.push(path)
}

/* TODO
export async function doCloudRequest(path: string, data: any): Promise<(AxiosResponse | null)> {
    return doRequest(cloudBaseUrl + path, data)
}
 */

export async function doHubRequest(path: string, data: any): Promise<(AxiosResponse | null)> {
    return doRequest(globalConfig.hubBaseUrl + path, data)
}

export const defaultAllowedSymbols = '[0-9a-zA-Z]';
export const tagAllowedSymbols = '[0-9a-zA-Z.]';
export const defaultMinLength = 3;
export const defaultMaxLength = 20;
export const minLengthPassword = 8;
export const maxLengthPassword = 30;

export function generateInvalidInputMessage(fieldName: string, allowedSymbols: string, minLength: number, maxLength: number): string {
    return `Invalid ${fieldName}, allowed symbols are ${allowedSymbols} and the length must be between ${minLength} and ${maxLength}.`;
}

export function getDefaultValidationRegex(): RegExp {
    return new RegExp(`^${defaultAllowedSymbols}{${defaultMinLength},${defaultMaxLength}}$`)
}

// TODO Somehow the backend needs to tell the frontend the domain of the hub.
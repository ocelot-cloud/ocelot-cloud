import router from "@/router";
import axios, {AxiosResponse} from "axios";

export const baseUrl = 'http://localhost:8082';

export function goToHubPage(path: string) {
    router.push('/hub' + path)
}

export const session = {
    user: "",
    isAuthenticated: false,
};

export function alertError(error: any) {
    if (axios.isAxiosError(error) && error.response) {
        const errorMessage = error.response.data || 'An unknown error occurred';
        alert(`An error occurred: ${errorMessage}`);
    } else {
        alert('An unknown error occurred');
    }
}

export async function doRequest(path: string, data: any): Promise<(AxiosResponse | null)> {
    try {
        const response = await axios.post(baseUrl + path, data);
        if (response.status !== 200) {
            throw new Error(response.data);
        }
        return response
    } catch (error) {
        alertError(error);
        return null
    }
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
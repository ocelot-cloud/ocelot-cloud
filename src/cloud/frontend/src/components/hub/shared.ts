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

export const allowedSymbols = '[0-9a-zA-Z]';
export const minLength = 3;
export const maxLength = 20;

export const usernamePattern = new RegExp(`^${allowedSymbols}{${minLength},${maxLength}}$`);
export const passwordPattern = new RegExp(`^${allowedSymbols}{${minLength},${maxLength}}$`);

export const globalUsernameErrorMessage = `Invalid input, allowed symbols are ${allowedSymbols} and the length must be between ${minLength} and ${maxLength}.`;
export const globalPasswordErrorMessage = `Invalid input, allowed symbols are ${allowedSymbols} and the length must be between ${minLength} and ${maxLength}.`;

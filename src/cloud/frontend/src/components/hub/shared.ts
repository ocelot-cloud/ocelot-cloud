import router from "@/router";
import axios, {AxiosResponse} from "axios";

export function goToHubPage(path: string) {
    router.push('/hub' + path)
}

export const session = {
    user: "",
    isAuthenticated: false,
};

export async function doRequest(path: string, data: any): Promise<(AxiosResponse | null)> {
    const baseUrl = 'http://localhost:8082';

    try {
        const response = await axios.post(baseUrl + path, data);
        if (response.status !== 200) {
            throw new Error(response.data);
        }
        return response
    } catch (error) {
        if (axios.isAxiosError(error) && error.response) {
            const errorMessage = error.response.data || 'An unknown error occurred';
            alert(`An error occurred: ${errorMessage}`);
        } else {
            alert('An unknown error occurred');
        }
        return null
    }
}
import router from "@/router";
import axios from "axios";

export function goToHubPage(path: string) {
    router.push('/hub' + path)
}

export const session = {
    user: "",
    isAuthenticated: false,
};

export async function doRequest(method: string, path: string, data: any) {
    const baseUrl = 'http://localhost:8082';

    try {
        const response = await axios.request({
            method,
            url: baseUrl + path,
            data: data,
        });
        if (response.status !== 200) {
            throw new Error(response.data);
        }
    } catch (error) {
        if (axios.isAxiosError(error) && error.response) {
            const errorMessage = error.response.data || 'An unknown error occurred';
            alert(`An error occurred: ${errorMessage}`);
        } else {
            alert('An unknown error occurred');
        }
    }
}
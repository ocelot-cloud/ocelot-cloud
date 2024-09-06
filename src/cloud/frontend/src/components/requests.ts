import axios, {AxiosResponse} from "axios";
import {globalConfig} from "@/components/GlobalConfig";

export function alertError(error: any) {
    if (axios.isAxiosError(error) && error.response) {
        const errorMessage = error.response.data || 'An unknown error occurred';
        alert(`An error occurred: ${errorMessage}`);
    } else {
        alert('An unknown error occurred');
    }
}

export async function doRequest(url: string, data: any): Promise<(AxiosResponse | null)> {
    try {
        const response = await axios.post(url, data, {
            withCredentials: true,
        });
        if (response.status !== 200) {
            throw new Error(response.data);
        }
        return response
    } catch (error) {
        alertError(error);
        return null
    }
}

export async function doCloudRequest(path: string, data: any): Promise<(AxiosResponse | null)> {
    return doRequest(globalConfig.cloudBaseUrl + path, data)
}
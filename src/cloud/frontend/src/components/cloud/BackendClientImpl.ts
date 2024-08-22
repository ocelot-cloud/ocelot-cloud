import {BackendClient} from "@/components/cloud/Shared";

export class BackendClientImpl implements BackendClient {
    Logout(baseUrl: string): void {
        fetch(baseUrl + '/api/logout', {
            method: 'POST',
        }).then(response => {
            if (!response.ok) {
                console.error("Logout failed with status:", response.status);
            }
        })
            .catch(error => {
                console.error("Error during logout:", error);
            });
    }

    async getResponsePromise(stackUrl: string): Promise<Response> {
        return await fetch(stackUrl + 'read');
    }

    async postRequest(name: string, stackUrl: string, endpoint: string): Promise<void> {
        const data = {
            name: name
        };

        try {
            const response = await fetch(stackUrl + endpoint, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(data)
            });
            await response;
            console.log(response.status);
        } catch (error) {
            console.error("There was an error with the fetch operation:", error);
        }
    }

}

export class BackendClient {
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
                credentials: 'include',
                body: JSON.stringify(data)
            });
            await response;
            console.log(response.status);
        } catch (error) {
            console.error("There was an error with the fetch operation:", error);
        }
    }

}

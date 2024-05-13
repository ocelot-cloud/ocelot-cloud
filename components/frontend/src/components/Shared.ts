export class Stack {
    name: string;
    state: string;
    urlPath: string;

    constructor(name: string, state: string, urlPath: string) {
        this.name = name;
        this.state = state;
        this.urlPath = urlPath;
    }
}

export interface BackendClient {
    getResponsePromise(stackUrl: string): Promise<Response>
    postRequest(name: string, stackUrl: string, endpoint: string): Promise<void>
    Logout(baseUrl: string): void
}
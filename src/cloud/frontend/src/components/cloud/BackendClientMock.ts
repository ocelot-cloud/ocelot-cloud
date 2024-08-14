import {BackendClient, Stack} from "@/components/cloud/Shared";


const stack1 = new Stack('nginx-default', 'Uninitialized', "/");
const stack2 = new Stack('nginx-custom-path', 'Uninitialized', "/custom-path");
const stack3 = new Stack('nginx-slow-start', 'Uninitialized', "/");

export class BackendClientMock implements BackendClient {

    stacks = [stack1, stack2, stack3]

    Logout(baseUrl: string): void {
        console.log("Mock logout called for baseUrl:", baseUrl);
    }

    async getResponsePromise(stackUrl: string): Promise<Response> {
        console.log("Mock getResponsePromise called with stackUrl:", stackUrl);
        const jsonResponse = JSON.stringify(this.stacks);
        return Promise.resolve(new Response(jsonResponse, {
            status: 200,
            statusText: 'OK',
            headers: {
                'Content-Type': 'application/json'
            }
        }));
    }

    async postRequest(name: string, stackUrl: string, endpoint: string): Promise<void> {
        console.log("Mock postRequest called with:", { name, stackUrl, endpoint });
        if (endpoint === "deploy") {
            this.changeStackStateTo(name, "Available")
            return
        } else if (endpoint === "stop") {
            this.changeStackStateTo(name, "Uninitialized")
            return
        } else {
            throw new Error("error, unknown endpoint addresses: " + endpoint);
        }
    }

    changeStackStateTo(name: string, state: string) {
        this.stacks.forEach(stack => {
            if (stack.name === name) {
                stack.state = state
            }
        })
    }

}
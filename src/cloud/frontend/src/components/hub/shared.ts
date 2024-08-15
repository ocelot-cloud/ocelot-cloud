import router from "@/router";

export function goToHubPage(path: string) {
    router.push('/hub' + path)
}

export const session = {
    user: "",
    isAuthenticated: false
};
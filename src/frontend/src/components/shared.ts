import router from "@/router";

export function goToPage(path: string) {
    router.push(path)
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
import { writable, get} from 'svelte/store';

export const token = writable(null);
export function bearer() {
    let t = get(token)
    if (t === null || t === undefined) {
        return ""
    }
    if (t.token === undefined) {
        return ""
    }
    return "Bearer " + t.token
}
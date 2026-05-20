const TOKEN_KEY = 'admin_token';

export function getToken(): string | null {
	if (typeof localStorage === 'undefined') return null;
	return localStorage.getItem(TOKEN_KEY);
}

export function setToken(token: string): void {
	localStorage.setItem(TOKEN_KEY, token);
}

export function clearToken(): void {
	localStorage.removeItem(TOKEN_KEY);
}

function authHeaders(): Record<string, string> {
	const token = getToken();
	return token ? { Authorization: `Bearer ${token}` } : {};
}

export async function adminFetch(path: string, init: RequestInit = {}): Promise<Response> {
	const res = await fetch(path, {
		...init,
		headers: { ...authHeaders(), ...(init.headers as Record<string, string>) }
	});
	if (res.status === 401) {
		clearToken();
		location.href = '/admin/login';
	}
	return res;
}

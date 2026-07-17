import { tokenJmeno } from '@/stores';

type JsonData = ReturnType<JSON['parse']>;

export type RequestConfig = {
    headers?: HeadersInit;
    params?: Record<string, string | number | boolean | null | undefined>;
    signal?: AbortSignal;
};

export type ApiResponse<T> = {
    data: T;
    status: number;
    headers: Headers;
};

export class ApiError<T = unknown> extends Error {
    response: ApiResponse<T>;

    constructor(response: ApiResponse<T>) {
        super(`Request failed with status code ${response.status}`);
        this.name = 'ApiError';
        this.response = response;
    }
}

export function getApiErrorMessage(error: unknown): string | undefined {
    if (!(error instanceof ApiError)) return undefined;

    const data = error.response.data;
    if (typeof data === 'string') return data;
    if (data && typeof data === 'object' && 'error' in data) {
        const message = (data as { error?: unknown }).error;
        if (typeof message === 'string') return message;
    }

    return undefined;
}

async function parseResponse(response: Response): Promise<unknown> {
    if (response.status === 204) return undefined;

    const text = await response.text();
    if (!text) return undefined;

    const contentType = response.headers.get('content-type') ?? '';
    if (contentType.includes('application/json') || contentType.includes('+json')) {
        try {
            return JSON.parse(text);
        } catch {
            if (response.ok) throw new SyntaxError(`Invalid JSON response from ${response.url}`);
            return text;
        }
    }

    return text;
}

async function request<T>(method: string, url: string, data?: unknown, config: RequestConfig = {}): Promise<ApiResponse<T>> {
    const requestUrl = new URL(`/api${url}`, window.location.origin);
    for (const [key, value] of Object.entries(config.params ?? {})) {
        if (value !== null && value !== undefined) requestUrl.searchParams.set(key, String(value));
    }

    const headers = new Headers(config.headers);
    const legacyToken = localStorage.getItem(tokenJmeno);
    if (legacyToken && !headers.has('Authorization')) headers.set('Authorization', `Bearer ${legacyToken}`);
    const init: RequestInit = { method, headers, credentials: 'same-origin', signal: config.signal };

    if (data !== undefined) {
        if (data instanceof FormData || data instanceof URLSearchParams || typeof data === 'string' || data instanceof Blob) {
            init.body = data;
        } else {
            headers.set('Content-Type', 'application/json');
            init.body = JSON.stringify(data);
        }
    }

    const response = await fetch(requestUrl, init);
    const apiResponse: ApiResponse<T> = {
        data: (await parseResponse(response)) as T,
        status: response.status,
        headers: response.headers,
    };

    if (!response.ok) throw new ApiError(apiResponse);
    return apiResponse;
}

const api = {
    get<T = JsonData>(url: string, config?: RequestConfig) {
        return request<T>('GET', url, undefined, config);
    },
    post<T = JsonData>(url: string, data?: unknown, config?: RequestConfig) {
        return request<T>('POST', url, data, config);
    },
    put<T = JsonData>(url: string, data?: unknown, config?: RequestConfig) {
        return request<T>('PUT', url, data, config);
    },
    patch<T = JsonData>(url: string, data?: unknown, config?: RequestConfig) {
        return request<T>('PATCH', url, data, config);
    },
    delete<T = JsonData>(url: string, config?: RequestConfig) {
        return request<T>('DELETE', url, undefined, config);
    },
};

export default api;

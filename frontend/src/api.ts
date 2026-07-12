type JsonData = ReturnType<JSON['parse']>;

export type RequestConfig = {
    headers?: HeadersInit;
    params?: Record<string, string | number | boolean | null | undefined>;
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

async function parseResponse(response: Response): Promise<unknown> {
    if (response.status === 204) return undefined;

    const contentType = response.headers.get('content-type');
    if (contentType?.includes('application/json')) return response.json();

    const text = await response.text();
    return text || undefined;
}

async function request<T>(method: string, url: string, data?: unknown, config: RequestConfig = {}): Promise<ApiResponse<T>> {
    const requestUrl = new URL(`/api${url}`, window.location.origin);
    for (const [key, value] of Object.entries(config.params ?? {})) {
        if (value !== null && value !== undefined) requestUrl.searchParams.set(key, String(value));
    }

    const headers = new Headers(config.headers);
    const init: RequestInit = { method, headers };

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

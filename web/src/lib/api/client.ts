const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export class ApiError extends Error {
  constructor(
    public status: number,
    public statusText: string,
    public body?: unknown
  ) {
    super(`API Error: ${status} ${statusText}`);
    this.name = 'ApiError';
  }
}

// Helper to get session from localStorage (works outside React components)
function getSessionHeaders(): Record<string, string> {
  if (typeof window === 'undefined') return {};

  try {
    const stored = localStorage.getItem('git-resume-session');
    if (!stored) return {};

    const { state } = JSON.parse(stored);
    const headers: Record<string, string> = {};

    if (state?.session?.id) {
      headers['X-Session-ID'] = state.session.id;
    }
    if (state?.apiKey) {
      headers['X-API-Key'] = state.apiKey;
    }

    return headers;
  } catch {
    return {};
  }
}

class ApiClient {
  private baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  private async request<T>(
    endpoint: string,
    options?: RequestInit,
    skipAuth = false
  ): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`;

    const sessionHeaders = skipAuth ? {} : getSessionHeaders();

    const response = await fetch(url, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...sessionHeaders,
        ...options?.headers,
      },
    });

    if (!response.ok) {
      let body;
      try {
        body = await response.json();
      } catch {
        body = await response.text();
      }
      throw new ApiError(response.status, response.statusText, body);
    }

    // Handle empty responses
    const text = await response.text();
    if (!text) return {} as T;

    return JSON.parse(text);
  }

  get<T>(endpoint: string, skipAuth = false): Promise<T> {
    return this.request<T>(endpoint, { method: 'GET' }, skipAuth);
  }

  post<T>(endpoint: string, body?: unknown, skipAuth = false): Promise<T> {
    return this.request<T>(
      endpoint,
      {
        method: 'POST',
        body: body ? JSON.stringify(body) : undefined,
      },
      skipAuth
    );
  }

  delete<T>(endpoint: string, skipAuth = false): Promise<T> {
    return this.request<T>(endpoint, { method: 'DELETE' }, skipAuth);
  }
}

export const api = new ApiClient(API_BASE);

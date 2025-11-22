// This service handles all API calls related to authentication.

const BASE_URL = "http://localhost:8080";

// A helper function to handle fetch responses
const handleResponse = async (response) => {
    if (!response.ok) {
        const error = await response.json().catch(() => ({ error: 'Server error with no JSON response' }));
        throw new Error(error.error || 'Request failed');
    }
    // Handle 204 No Content response
    if (response.status === 204) {
        return null;
    }
    return response.json();
};

export const authService = {
  async login({ login, password }) {
    const response = await fetch(`${BASE_URL}/auth`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ login, password }),
      credentials: 'include', // This is crucial for sending cookies
    });
    return handleResponse(response);
  },

  async register({ login, email, password }) {
    const response = await fetch(`${BASE_URL}/register`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ login, email, password }),
      credentials: 'include',
    });
    return handleResponse(response);
  },

  async logout() {
    const response = await fetch(`${BASE_URL}/api/logout`, {
        method: "POST",
        credentials: 'include',
    });
    return handleResponse(response);
  },

  async getMe() {
    const response = await fetch(`${BASE_URL}/api/me`, {
        method: "GET",
        credentials: 'include',
    });
    return handleResponse(response);
  }
};
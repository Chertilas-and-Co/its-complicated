import React, { createContext, useState, useContext, useEffect } from 'react';
import { authService } from '../services/authService';

// 1. Create the context
const AuthContext = createContext(null);

// 2. Create the provider component
export function AuthProvider({ children }) {
    const [user, setUser] = useState(null);
    const [isLoading, setIsLoading] = useState(true); // To handle initial page load check

    // Check if user is already logged in on initial load
    useEffect(() => {
        authService.getMe()
            .then(userData => {
                setUser(userData);
            })
            .catch(() => {
                setUser(null);
            })
            .finally(() => {
                setIsLoading(false);
            });
    }, []);

    const login = async (credentials) => {
        // credentials are { login, password }
        await authService.login(credentials);
        const userData = await authService.getMe();
        setUser(userData);
    };

    const logout = async () => {
        await authService.logout();
        setUser(null);
    };

    const authValue = {
        user,
        isLoading,
        login,
        logout,
    };

    return (
        <AuthContext.Provider value={authValue}>
            {children}
        </AuthContext.Provider>
    );
}

// 3. Create a custom hook for easy access to the context
export function useAuth() {
    return useContext(AuthContext);
}

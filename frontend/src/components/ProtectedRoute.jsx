import React from 'react';
import { Navigate, Outlet } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

const ProtectedRoute = () => {
    const { user, isLoading } = useAuth();

    // While the initial check is running, show a loading indicator
    if (isLoading) {
        return <div>Loading authentication status...</div>;
    }

    // If loading is finished and there is no user, redirect to login
    if (!user) {
        return <Navigate to="/login" replace />;
    }

    // If user is authenticated, render the child route
    return <Outlet />;
};

export default ProtectedRoute;

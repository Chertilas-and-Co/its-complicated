import React from 'react';
import Button from '@mui/material/Button';
import { useNavigate } from 'react-router-dom';
import { useAuth } from './context/AuthContext';

export default function Home() {
  const { logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <div>
      <h1>Welcome to the Home Page!</h1>
      <Button variant="contained" onClick={handleLogout}>
        Logout
      </Button>
    </div>
  );
}

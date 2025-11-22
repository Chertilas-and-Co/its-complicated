import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from './context/AuthContext'; // Import useAuth
import './Navbar.css';

const Navbar = () => {
  const { user, logout } = useAuth(); // Get user and logout function from context
  const navigate = useNavigate();

  const handleLogout = async () => {
    await logout();
    navigate('/login'); // Redirect to login after logout
  };

  return (
    <nav className="navbar" style={{fontFamily: 'Arial, sans-serif'}}>
      <div className="navbar-left">
        <Link to="/" className="logo">Страдать</Link>
      </div>
      <div className="navbar-center">
        <ul className="nav-links">
          {/* Show profile and friends only if user is logged in */}
          {user && (
            <>
              <li><Link to="/user">Профиль</Link></li>
              <li><Link to="/friends">Друзья</Link></li>
            </>
          )}
          <li><Link to="/">Сообщества</Link></li>
        </ul>
      </div>
      <div className="navbar-right">
        {user ? (
          <button onClick={handleLogout} className="exit">Выйти</button>
        ) : (
          <>
            <Link to="/login" style={{ marginRight: '10px' }}>Войти</Link>
            <Link to="/register">Регистрация</Link>
          </>
        )}
      </div>
    </nav>
  );
};

export default Navbar;
import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import './Navbar.css';
import { useAuth } from './context/AuthContext';

const Navbar = () => {
    const { logout } = useAuth();
    const navigate = useNavigate();

    const handleLogout = () => {
        logout();
        navigate('/login');
    };

    return (
        <nav className="navbar" style={{ fontFamily: 'Arial, sans-serif' }}>
            <div className="navbar-left">
                <Link to="/" className="logo">Страдать</Link>
            </div>
            <div className="navbar-center">
                <ul className="nav-links">
                    <li><Link to="/user">Профиль</Link></li>
                    <li><Link to="/">Сообщества</Link></li>
                    <li><Link to="/friends">Друзья</Link></li>
                </ul>
            </div>
            <div className="navbar-right">
                <button onClick={handleLogout} className="exit" style={{ background: 'none', border: 'none', color: 'white', cursor: 'pointer', fontSize: '1em' }}>Выйти</button>
            </div>
        </nav>
    );
};

export default Navbar;


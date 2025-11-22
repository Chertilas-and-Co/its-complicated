import React from 'react';
import { Link } from 'react-router-dom';
import './Navbar.css';

const Navbar = () => {
  return (
    <nav className="navbar" style={{fontFamily: 'Arial, sans-serif'}}>
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
        {/* This should probably be a button that triggers a logout function, 
            but for now it links to the login page. */}
        <Link to="/login" className="exit">Выйти</Link>
      </div>
    </nav>
  );
};

export default Navbar;
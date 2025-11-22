import React from 'react';
import { Link } from 'react-router-dom';
import './Navbar.css';

const Navbar = () => {
  return (
    <nav className="navbar" style={{fontFamily: 'Arial, sans-serif'}}>
      <div className="navbar-left">
        <a href="/" className="logo">Страдать</a>
      </div>
      <div className="navbar-center">
        <ul className="nav-links">
          <li><a href="/user">Профиль</a></li>
          <li><a href="/communities">Сообщества</a></li>
          <li><a href="/friends">Друзья</a></li>
        </ul>
      </div>
      <div className="navbar-right">
        <a href="/reg" className="exit">Выйти</a>
      </div>
    </nav>
  );
};

export default Navbar;

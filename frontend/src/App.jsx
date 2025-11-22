import React from 'react';
import { Outlet } from 'react-router-dom';
import Navbar from './Navbar';

function App() {
  return (
    <>
      <Navbar />
      <main style={{ paddingTop: '80px', padding: '20px' }}>
        {/* Outlet will render the child routes from main.jsx */}
        <Outlet />
      </main>
    </>
  );
}

export default App;

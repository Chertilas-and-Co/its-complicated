import React from "react";
import { Routes, Route, BrowserRouter } from 'react-router-dom';
import Navbar from './Navbar';
import Home from './Home';
import { User } from './User';
import { Communities } from './Communities';
import { Friends } from './Friends';

function App() {
  return (
    <React.StrictMode>
    <BrowserRouter>
      <Navbar />
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/user" element={<User />} />
        <Route path="/communities" element={<Communities />} />
        <Route path="/friends" element={<Friends />} />
      </Routes>
    </BrowserRouter>
    </React.StrictMode>
  );
}

export default App;

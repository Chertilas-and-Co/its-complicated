import React from "react";
import ReactDOM from "react-dom/client";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { createRoot } from 'react-dom/client'
import CommunitiesPage from './AllCommunities';
import Home from './Home';
import UserPage from './User';
import CommunityPage  from './Communities';
import Friends from './Friends';
createRoot(document.getElementById('root')).render(
 <React.StrictMode>
    <BrowserRouter>
   
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/user" element={<UserPage />} />
        <Route path="/communities" element={<CommunitiesPage />} />
        <Route path="/friends" element={<Friends />} />
        <Route path="/community" element={<CommunityPage />} />
      </Routes>
    </BrowserRouter>
  </React.StrictMode>
);

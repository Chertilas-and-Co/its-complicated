import React from "react";
import ReactDOM from "react-dom/client";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { createRoot } from 'react-dom/client'
import CommunitiesPage from './AllCommunities';
import Home from './Home';
import UserPage from './User';
import CommunityPage  from './Communities';
import Friends from './Friends';
import DiscussionPage from './Discussion';
import CreateCommunityPage from './CreateCommunity';
createRoot(document.getElementById('root')).render(
 <React.StrictMode>
    <BrowserRouter>
   
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/user" element={<UserPage />} />
        <Route path="/communities" element={<CommunitiesPage />} />
        <Route path="/friends" element={<Friends />} />
        <Route path="/community" element={<CommunityPage />} />
        <Route path="/discussion" element={<DiscussionPage />} />
        <Route path="/create_community" element={<CreateCommunityPage />} />
      </Routes>
    </BrowserRouter>
  </React.StrictMode>
);

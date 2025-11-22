import React from "react";
import ReactDOM from "react-dom/client";
import SignIn from "./SignIn.jsx"
import Register from "./Register.jsx"
import CommunityPage from "./Communities.jsx"
import CommunitiesPage from "./AllCommunities.jsx"
import ForceGraph from "./Graph.jsx"
import { BrowserRouter, Routes, Route } from "react-router-dom";
import Friends from "./Friends.jsx";
import UserPage from "./User.jsx";
ReactDOM.createRoot(document.getElementById("root")).render(
    <React.StrictMode>
        <BrowserRouter>
            <Routes>
                <Route path="/communities" element={<CommunitiesPage />} />
                <Route path="/graph" element={<ForceGraph />} />
                <Route path="/auth" element={<SignIn />} />
                <Route path="/register" element={<Register />} />
                <Route path="/friends" element={<Friends />} />
                <Route path="/user" element={<UserPage />} />
            </Routes>
        </BrowserRouter>
    </React.StrictMode>
);

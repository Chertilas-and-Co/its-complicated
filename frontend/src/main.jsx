import React from "react";
import ReactDOM from "react-dom/client";
import { BrowserRouter, Routes, Route } from "react-router-dom";

// Import global styles and the main App layout component
import './index.css';
import App from './App.jsx';

// Your page components
import SignIn from "./SignIn.jsx";
import Register from "./Register.jsx";
import CommunityPage from "./Communities.jsx";
import CommunitiesPage from "./AllCommunities.jsx";
import GraphPage from "./pages/GraphPage.jsx"; // Our styled page
import UserPage from "./User.jsx";
import Friends from "./Friends.jsx";
import DiscussionPage from "./Discussion.jsx";
import CreateCommunityPage from "./CreateCommunity.jsx";
import CreatePostPage from "./CreatePost.jsx";

ReactDOM.createRoot(document.getElementById("root")).render(
    <React.StrictMode>
        <BrowserRouter>
            <Routes>
                {/* All routes are now children of App to get the Navbar and common layout */}
                <Route path="/" element={<App />}>
                    {/* The index route is now the graph */}
                    <Route index element={<GraphPage />} /> 
                    <Route path="communities" element={<CommunitiesPage />} />
                    <Route path="graph" element={<GraphPage />} />
                    <Route path="auth" element={<SignIn />} />
                    <Route path="register" element={<Register />} />
                    <Route path="user" element={<UserPage />} />
                    <Route path="friends" element={<Friends />} />
                    <Route path="create_post" element={<CreatePostPage/>}/>
                    <Route path="discussion" element={<DiscussionPage />} />
                    <Route path="create_community" element={<CreateCommunityPage />} />
                    {/* The route for a single community we added earlier */}
                    <Route path="community/:id" element={<CommunityPage />} />
                </Route>
            </Routes>
        </BrowserRouter>
    </React.StrictMode>
);
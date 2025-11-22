import React from "react";
import {
  BrowserRouter as Router,
  Routes,
  Route,
  Navigate,
} from "react-router-dom";
import { AuthProvider, useAuth } from "./context/AuthContext";
import Register from "./Register";
import Home from "./Home";
import SignIn from "./SignIn";
import ProtectedRoute from "./ProtectedRoute";
import CommunityPage from "./Communities.jsx"
import CommunitiesPage from "./AllCommunities.jsx"
import ForceGraph from "./Graph.jsx"
import UserPage from "./User.jsx";
import Friends from "./Friends.jsx";
import DiscussionPage from "./Discussion.jsx";
import CreateCommunityPage from "./CreateCommunity.jsx";
import CreatePostPage from "./CreatePost.jsx";

const AppRoutes = () => {
  const { user, logout } = useAuth();

  return (
    <Routes>
      <Route
        path="/login"
        element={user ? <Navigate to="/home" /> : <SignIn />}
      />
      <Route
        path="/register"
        element={user ? <Navigate to="/home" /> : <Register />}
      />
      <Route
        path="/home"
        element={
          <ProtectedRoute>
            <Home onLogout={logout} />
          </ProtectedRoute>
        }
      />
        <Route path="/community/:id" element={<CommunityPage />} />
        <Route path="/communities" element={<CommunitiesPage />} />
        <Route path="/graph" element={<ForceGraph />} />
        <Route path="/user" element={<UserPage />} />
        <Route path="/friends" element={<Friends />} />
        <Route path="/create_post" element={<CreatePostPage/>}/>
        <Route path="/discussion" element={<DiscussionPage />} />
        <Route path="/create_community" element={<CreateCommunityPage />} />
      <Route path="*" element={<Navigate to="/login" />} />
    </Routes>
  );
};

const App = () => {
  return (
    <Router>
      <AuthProvider>
        <AppRoutes />
      </AuthProvider>
    </Router>
  );
};

export default App;
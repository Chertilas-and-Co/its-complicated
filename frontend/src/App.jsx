import React from 'react';
import { BrowserRouter as Router, Routes, Route, Outlet } from 'react-router-dom';
// Using absolute paths from the /src root
import Navbar from '/src/Navbar.jsx';
import LoginPage from '/src/pages/loginPage.jsx';
import RegisterPage from '/src/Register.jsx';
import GraphPage from '/src/pages/GraphPage.jsx';
import CommunityPage from '/src/Communities.jsx';
import CommunitiesPage from '/src/AllCommunities.jsx';
import UserPage from '/src/User.jsx';
import Friends from '/src/Friends.jsx';
import DiscussionPage from '/src/Discussion.jsx';
import CreateCommunityPage from '/src/CreateCommunity.jsx';
import CreatePostPage from '/src/CreatePost.jsx';
import ProtectedRoute from '/src/components/ProtectedRoute.jsx';
import { useAuth } from '/src/context/AuthContext.jsx';

const AppLayout = () => {
    const { user } = useAuth();
    return (
        <>
            {user && <Navbar />}
            <main style={{ paddingTop: user ? '80px' : '0', padding: '20px' }}>
                <Outlet />
            </main>
        </>
    );
};

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<AppLayout />}>
          <Route index element={<GraphPage />} />
          <Route path="graph" element={<GraphPage />} />
          <Route path="login" element={<LoginPage />} />
          <Route path="register" element={<RegisterPage />} />
          <Route path="community/:id" element={<CommunityPage />} />
          <Route path="communities" element={<CommunitiesPage />} />
          <Route element={<ProtectedRoute />}>
            <Route path="user" element={<UserPage />} />
            <Route path="friends" element={<Friends />} />
            <Route path="create_post" element={<CreatePostPage />} />
            <Route path="discussion" element={<DiscussionPage />} />
            <Route path="create_community" element={<CreateCommunityPage />} />
          </Route>
        </Route>
      </Routes>
    </Router>
  );
}

export default App;

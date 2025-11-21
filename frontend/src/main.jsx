import React from "react";
import ReactDOM from "react-dom/client";
import SignIn from "./SignIn.jsx"
import Register from "./Register.jsx"
import { BrowserRouter, Routes, Route } from "react-router-dom";

ReactDOM.createRoot(document.getElementById("root")).render(
  <React.StrictMode>
    <BrowserRouter>
      <Routes>
        <Route path="/auth" element={<SignIn />} />
        <Route path="/register" element={<Register />} />
      </Routes>
    </BrowserRouter>
  </React.StrictMode>
);

import React, { useState } from "react";
import AuthForm from "../components/AuthForm";
import { login } from "../api/auth";

function LoginPage() {
  const [error, setError] = useState(""); // для отображения ошибок

  const handleLogin = async ({ login: username, password }) => {
    setError("");
    try {
      await login({ login: username, password });
      alert("Успешный вход!"); // сообщение или редирект
    } catch (e) {
      setError(e.message);
    }
  };

  return (
    <div style={{ maxWidth: 350, margin: "100px auto" }}>
      <h2>Вход</h2>
      <AuthForm onSubmit={handleLogin} />
      {error && <div style={{ color: "red", marginTop: 12 }}>{error}</div>}
    </div>
  );
}

export default LoginPage;

import React, { useState } from "react";

const AuthPage = () => {
  const [login, setLogin] = useState("");
  const [password, setPassword] = useState("");
  const [message, setMessage] = useState("");

  const handleSubmit = async (e) => {
    e.preventDefault();
    setMessage("");

    try {
      const response = await fetch("http://localhost:8080/auth", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ login, password }),
      });

      if (response.ok) {
        setMessage("Авторизация успешна!");
      } else {
        const text = await response.text();
        setMessage(`Ошибка: ${text}`);
      }
    } catch (error) {
      setMessage(`Ошибка сети: ${error.message}`);
    }
  };

  return (
    <div style={{ maxWidth: 350, margin: "50px auto", padding: 20, border: "1px solid #eee", borderRadius: 5 }}>
      <h2>Авторизация</h2>
      <form onSubmit={handleSubmit}>
        <div>
          <label>
            Логин:<br />
            <input
              type="text"
              value={login}
              onChange={e => setLogin(e.target.value)}
              autoComplete="username"
              required
              style={{ width: "100%" }}
            />
          </label>
        </div>
        <div style={{ marginTop: 12 }}>
          <label>
            Пароль:<br />
            <input
              type="password"
              value={password}
              onChange={e => setPassword(e.target.value)}
              autoComplete="current-password"
              required
              style={{ width: "100%" }}
            />
          </label>
        </div>
        <button type="submit" style={{ marginTop: 16, width: "100%" }}>
          вфыоовлофы
        </button>
      </form>
      {message && <div style={{ marginTop: 20, color: "crimson" }}>{message}</div>}
    </div>
  );
};

export default AuthPage;

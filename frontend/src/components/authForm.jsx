import React, { useState } from "react";

function AuthForm({ onSubmit }) {
  const [login, setLogin] = useState("");
  const [password, setPassword] = useState("");

  const handleSubmit = (e) => {
    e.preventDefault(); // чтобы не перезагружать страницу
    onSubmit({ login, password });
  };

  return (
    <form onSubmit={handleSubmit}>
      <div>
        <label>Логин</label>
        <input
          type="text"
          value={login}
          onChange={(e) => setLogin(e.target.value)}
          required
          autoFocus
        />
      </div>
      <div>
        <label>Пароль</label>
        <input
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
        />
      </div>
      <button type="submit">Войвыоофв</button>
    </form>
  );
}

export default AuthForm;

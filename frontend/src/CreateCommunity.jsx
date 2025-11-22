import React, { useState } from "react";
import Navbar from "./Navbar";
const CreateCommunityPage = () => {
  // Поля формы
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [message, setMessage] = useState("");

  // Данные создателя и время создания — здесь примерные
  const creator = "Иван"; // сюда вставить логин текущего пользователя из контекста/состояния
  const createdAt = new Date().toLocaleString("ru-RU", {
    day: "numeric",
    month: "long",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit"
  });

  const handleSubmit = async (e) => {
    e.preventDefault();
    setMessage("");

    const newCommunity = {
      name,
      description,
      creator,
      createdAt
    };

    try {
      const response = await fetch("http://localhost:8080/communities", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(newCommunity),
      });

      if (response.ok) {
        setMessage("Сообщество успешно создано!");
        // Очистить форму при успешном создании
        setName("");
        setDescription("");
      } else {
        const text = await response.text();
        setMessage(`Ошибка: ${text}`);
      }
    } catch (error) {
      setMessage(`Ошибка сети: ${error.message}`);
    }
  };

  return (
    <>
    <Navbar/>
    <div style={{
      maxWidth: 520,
      margin: "50px auto",
      padding: 20,
      border: "1px solid #eee",
      borderRadius: 5,
      fontFamily: "Arial, sans-serif",
      backgroundColor: "#fff"
    }}>
      <h2>Создать сообщество</h2>
      <form onSubmit={handleSubmit}>
        <div>
          <label>
            Название:<br />
            <input
              type="text"
              value={name}
              onChange={e => setName(e.target.value)}
              required
              style={{ width: "100%", padding: "8px", marginTop: 6, boxSizing: "border-box" }}
              placeholder="Введите название сообщества"
            />
          </label>
        </div>
        <div style={{ marginTop: 16 }}>
          <label>
            Описание:<br />
            <textarea
              value={description}
              onChange={e => setDescription(e.target.value)}
              required
              rows={4}
              style={{ width: "100%", padding: "8px", marginTop: 6, boxSizing: "border-box" }}
              placeholder="Введите описание сообщества"
            />
          </label>
        </div>
        <div style={{
          marginTop: 16,
          color: "#555",
          fontSize: "14px",
          fontStyle: "italic"
        }}>
          Создатель: <b>{creator}</b><br />
          Дата создания: <b>{createdAt}</b>
        </div>
        <button
          type="submit"
          style={{
            marginTop: 24,
            width: "100%",
            padding: "10px",
            backgroundColor: "#282c34",
            color: "white",
            fontSize: "16px",
            border: "none",
            borderRadius: "6px",
            cursor: "pointer"
          }}
        >
          Создать сообщество
        </button>
      </form>
      {message && (
        <div style={{ marginTop: 20, color: message.startsWith("Ошибка") ? "crimson" : "green" }}>
          {message}
        </div>
      )}
    </div>
    </>
  );
};

export default CreateCommunityPage;

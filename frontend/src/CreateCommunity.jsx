import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import Navbar from "./Navbar";
import { useAuth } from "./context/AuthContext";
import "./CreateCommunity.css";

const CreateCommunityPage = () => {
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [message, setMessage] = useState("");
  const { user } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setMessage("");

    if (!user) {
      setMessage("Ошибка: Вы должны быть авторизованы, чтобы создать сообщество.");
      return;
    }

    const newCommunity = {
      name,
      description,
    };

    try {
      const response = await fetch("/api/communities", { // Changed to relative path
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include", // Отправляем cookie
        body: JSON.stringify(newCommunity),
      });

      const data = await response.json();
      console.log("Community creation response data:", data); // Added log

      if (response.ok) {
        setMessage("Сообщество успешно создано!");
        setName("");
        setDescription("");
        // Перенаправляем на страницу нового сообщества
        if (data.id) { // Ensure data.id exists before navigating
          navigate(`/community/${data.id}`);
        } else {
          console.error("Community ID not found in response data:", data);
          setMessage("Ошибка: ID нового сообщества не найден.");
        }
      } else {
        setMessage(`Ошибка: ${data.error || "Не удалось создать сообщество"}`);
      }
    } catch (error) {
      setMessage(`Ошибка сети: ${error.message}`);
    }
  };

  return (
    <>
      <Navbar />
      <div className="create-community-container">
        <h2>Создать сообщество</h2>
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label htmlFor="name">Название:</label>
            <input
              id="name"
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
              placeholder="Введите название сообщества"
            />
          </div>
          <div className="form-group">
            <label htmlFor="description">Описание:</label>
            <textarea
              id="description"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              required
              rows={4}
              placeholder="Введите описание сообщества"
            />
          </div>
          {user && (
            <div className="creator-info">
              Создатель: <b>{user.username}</b>
            </div>
          )}
          <button type="submit" className="submit-btn">
            Создать сообщество
          </button>
        </form>
        {message && (
          <div className={`message ${message.startsWith("Ошибка") ? "error" : "success"}`}>
            {message}
          </div>
        )}
      </div>
    </>
  );
};

export default CreateCommunityPage;


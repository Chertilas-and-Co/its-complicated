import React, { useState } from "react";
import Navbar from "./Navbar";

const CreatePostPage = () => {
  // Поля формы
  const [content, setContent] = useState("");
  const [image, setImage] = useState("");
  const [message, setMessage] = useState("");

  const creator = "Иван"; // Логин текущего пользователя
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

    const newPost = {
      content,
      image,
      creator,
      createdAt
    };

    try {
      const response = await fetch("http://localhost:8080/posts", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(newPost),
      });

      if (response.ok) {
        setMessage("Пост успешно создан!");
        setContent("");
        setImage("");
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
        <h2>Создать пост</h2>
        <form onSubmit={handleSubmit}>
          <div>
            <label>
              Текст поста:<br />
              <textarea
                value={content}
                onChange={e => setContent(e.target.value)}
                required
                rows={4}
                style={{ width: "100%", padding: "8px", marginTop: 6, boxSizing: "border-box" }}
                placeholder="Введите текст поста"
              />
            </label>
          </div>
          <div style={{ marginTop: 16 }}>
            <label>
              URL изображения (необязательно):<br />
              <input
                type="text"
                value={image}
                onChange={e => setImage(e.target.value)}
                placeholder="Введите URL картинки"
                style={{ width: "100%", padding: "8px", marginTop: 6, boxSizing: "border-box" }}
              />
            </label>
          </div>
          <div style={{
            marginTop: 16,
            color: "#555",
            fontSize: "14px",
            fontStyle: "italic"
          }}>
            Автор: <b>{creator}</b><br />
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
            Создать пост
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

export default CreatePostPage;

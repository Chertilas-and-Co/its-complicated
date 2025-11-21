export async function login({ login, password }) {
  const response = await fetch("http://localhost:8080/auth", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ login, password }),
  });

  if (!response.ok) {
    const text = await response.text();
    throw new Error(text || "Ошибка авторизации");
  }

  return true;
}

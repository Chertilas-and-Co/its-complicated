import React, { useEffect, useState } from 'react';
import Navbar from './Navbar';
import UserAvatar from '/user.jpg'
import { useParams } from 'react-router-dom';
import { useAuth } from './context/AuthContext'; // Import useAuth

function Search() {
    return (
        <div style={{
            width: '100%',
            maxWidth: '300px',
            display: 'flex',
            alignItems: 'center',
            gap: '10px',
        }}>
            <input
                type="text"
                placeholder="Поиск..."
                style={{
                    flex: 1,
                    padding: '8px 14px',
                    borderRadius: '6px',
                    border: '1px solid #ccc',
                    fontSize: '16px',
                    color: 'black',
                    background: 'white'
                }}
            />
            <button
                style={{
                    padding: '8px 16px',
                    background: '#222',
                    color: 'white',
                    border: 'none',
                    borderRadius: '6px',
                    cursor: 'pointer',
                    fontSize: '16px'
                }}
            >
                Найти
            </button>
        </div>
    );
}

function Sidebar({ title, children }) {
    return (
        <div style={{
            width: '280px',
            padding: '15px',
            boxSizing: 'border-box',
            border: '1px solid #ddd',
            borderRadius: '8px',
            backgroundColor: '#fafafa'
        }}>
            <h3>{title}</h3>
            {children}
        </div>
    );
}

function UserList({ users }) {
    return (
        <ul style={{ listStyle: 'none', padding: 0, margin: 0 }}>
            {users.map(user => (
                <li key={user.id} style={{
                    marginBottom: '12px',
                    display: 'flex',
                    alignItems: 'center'
                }}>
                    <img
                        src={UserAvatar}
                        style={{
                            width: '36px',
                            height: '36px',
                            borderRadius: '50%',
                            marginRight: '10px'
                        }}
                    />
                    <span>{user.username || user.email || "Без имени"}</span>
                </li>
            ))}
        </ul>
    );
}

function PostsList({ posts }) {
    return (
        <>
            {posts.map(post => (
                <div
                    key={post.id}
                    style={{
                        border: '1px solid #ccc',
                        padding: '15px',
                        marginBottom: '15px',
                        borderRadius: '6px',
                        backgroundColor: '#fff',
                        display: 'flex',
                        flexDirection: 'column',
                        color: 'black'
                    }}
                >
                    {/* Верхняя часть: текст и изображение*/}
                    <div style={{ display: 'flex', gap: '15px', alignItems: 'flex-start' }}>
                        <p style={{ margin: 0, flex: 1 }}>{post.content}</p>
                    </div>
                    <div >
                        {post.image && (
                            <img
                                src={post.image}
                                alt="Post"
                                style={{
                                    width: 'auto',
                                    height: 'auto',
                                    objectFit: 'cover',
                                    borderRadius: '6px',
                                }}
                            />
                        )}
                    </div>
                    {/* Нижняя строка: Автор слева, кнопка по центру, дата справа */}
                    <div style={{
                        display: 'flex',
                        alignItems: 'center',
                        minHeight: '40px',
                        marginTop: '10px',
                        gap: '10px',
                        width: '100%'
                    }}>
                        <div style={{ flex: 1 }}>{post.author}</div>
                        <button style={{
                            background: '#222',
                            color: '#fff',
                            border: 'none',
                            borderRadius: '8px',
                            padding: '7px 18px',
                            fontSize: '15px',
                            cursor: 'pointer',
                            outline: 'none',
                            margin: '0 auto',
                            display: 'block'
                        }}>
                            Обсудить
                        </button>
                        <div style={{
                            flex: 1,
                            textAlign: 'right',
                            color: '#666'
                        }}>
                            {post.date}
                        </div>
                    </div>
                </div>
            ))}
        </>
    );
}

function CommunityPage() {
    const { id } = useParams();
    const [community, setCommunity] = useState(null);
    const [isSubscribed, setIsSubscribed] = useState(false);
    const { user } = useAuth(); // Get current user from AuthContext

    // Function to fetch community details
    const fetchCommunityDetails = async () => {
        try {
            const response = await fetch(`http://localhost:8080/community/${id}`);
            if (!response.ok) throw new Error('Ошибка при загрузке сообщества');
            const data = await response.json();
            setCommunity(data);
        } catch (error) {
            console.error("Ошибка при загрузке сообщества:", error);
        }
    };

    // Function to fetch subscription status
    const fetchSubscriptionStatus = async () => {
        if (!user || !user.id) { // Only check if user is logged in
            setIsSubscribed(false);
            return;
        }
        try {
            const response = await fetch(`http://localhost:8080/api/community/${id}/is_subscribed`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include',
            });
            if (!response.ok) {
                // If not 200 OK, assume not subscribed or error
                setIsSubscribed(false);
                return;
            }
            const data = await response.json();
            setIsSubscribed(data.is_subscribed);
        } catch (error) {
            console.error("Ошибка при проверке статуса подписки:", error);
            setIsSubscribed(false);
        }
    };

    // Effect to fetch community details and subscription status
    useEffect(() => {
        fetchCommunityDetails();
        fetchSubscriptionStatus();
    }, [id, user]); // Re-run when community ID or user changes

    const handleSubscribe = async () => {
        if (!user || !user.id) {
            alert("Пожалуйста, войдите, чтобы подписаться.");
            return;
        }
        try {
            const response = await fetch(`http://localhost:8080/api/community/${id}/subscribe`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include',
            });
            if (response.ok) {
                setIsSubscribed(true);
                fetchCommunityDetails(); // Refresh community data to update subscriber count
            } else {
                const errorData = await response.json();
                alert(`Ошибка подписки: ${errorData.error || response.statusText}`);
            }
        } catch (error) {
            console.error("Ошибка при подписке:", error);
            alert("Произошла ошибка при подписке.");
        }
    };

    const handleUnsubscribe = async () => {
        if (!user || !user.id) {
            alert("Пожалуйста, войдите, чтобы отписаться.");
            return;
        }
        try {
            const response = await fetch(`http://localhost:8080/api/community/${id}/subscribe`, {
                method: 'DELETE',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include',
            });
            if (response.ok) {
                setIsSubscribed(false);
                fetchCommunityDetails(); // Refresh community data to update subscriber count
            } else {
                const errorData = await response.json();
                alert(`Ошибка отписки: ${errorData.error || response.statusText}`);
            }
        } catch (error) {
            console.error("Ошибка при отписке:", error);
            alert("Произошла ошибка при отписке.");
        }
    };

    if (!community) {
        return <div>Загрузка...</div>;
    }

    // Маппинг подписчиков на объект с нужными полями для UserList
    const users = (community.subscribers || []).map(u => ({ // Safely access subscribers
        id: u.id,
        username: u.username,
        email: u.email,
        avatar_url: u.avatar_url,
        bio: u.bio
    }));

    // Здесь posts оставим пустым или можно загрузить отдельно, по примеру вашего кода
    const posts = [];

    return (
        <>
            <Navbar />
            <div style={{
                display: 'flex',
                flexDirection: 'column',
                height: '100vh',
                fontFamily: 'Arial, sans-serif',
                color: 'black'
            }}>
                {/* Шапка: поиск слева, название по центру, кнопка справа */}
                <div style={{
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'space-between',
                    padding: '15px 30px 10px 30px',
                    backgroundColor: '#f0f2f5'
                }}>
                    <Search />
                    <div style={{
                        flex: 2,
                        textAlign: 'center',
                        fontSize: '28px',
                        fontWeight: 'bold',
                        color: '#282c34'
                    }}>
                        {community.name}
                    </div>
                    <div style={{ flex: 1, maxWidth: '340px', textAlign: 'right' }}>
                        <button
                            onClick={isSubscribed ? handleUnsubscribe : handleSubscribe}
                            disabled={!user || !user.id}
                            style={{
                                padding: '8px 16px',
                                background: isSubscribed ? '#dc3545' : '#28a745', // Red for leave, green for join
                                color: 'white',
                                border: 'none',
                                borderRadius: '6px',
                                cursor: 'pointer',
                                fontSize: '16px',
                                opacity: (!user || !user.id) ? 0.6 : 1,
                            }}
                        >
                            {isSubscribed ? 'Покинуть' : 'Присоединиться'}
                        </button>
                    </div>
                </div>
                {/* Основной контент */}
                <div style={{
                    width: '95vw',
                    height: '100vh',
                    flex: 1, display: 'flex', gap: '20px', padding: '20px', backgroundColor: '#f0f2f5', color: 'black'
                }}>
                    <Sidebar title="Описание сообщества">
                        <p style={{ fontSize: '14px', lineHeight: '1.4', color: 'black' }}>{community.description}</p>
                    </Sidebar>

                    <div style={{
                        flex: 1, overflowY: 'auto',
                        scrollbarWidth: 'none',
                        msOverflowStyle: 'none'
                    }}>
                        <PostsList posts={posts} />
                    </div>

                    <Sidebar title="Пользователи">
                        <UserList users={users} />
                    </Sidebar>
                </div>
            </div>
        </>
    );
}

export default CommunityPage;

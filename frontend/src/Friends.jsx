import React from 'react';
import Navbar from './Navbar';


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

// Компактный список для правого сайдбара
function FriendListCompact({ friends }) {
  return (
    <ul style={{ listStyle: 'none', padding: 0, margin: 0 }}>
      {friends.map(friend => (
        <li key={friend.id} style={{
          marginBottom: '12px',
          display: 'flex',
          alignItems: 'center'
        }}>
          <img
            src={friend.avatar}
            alt={friend.name}
            style={{
              width: '36px',
              height: '36px',
              borderRadius: '50%',
              marginRight: '10px'
            }}
          />
          <span>{friend.name}</span>
        </li>
      ))}
    </ul>
  );
}

// КАРТОЧКА друга в «формате поста»
function FriendCard({ friend }) {
  return (
    <div
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
      {/* Верхняя часть: аватар + базовая инфа */}
      <div style={{ display: 'flex', gap: '15px', alignItems: 'flex-start' }}>
        <img
          src={friend.avatar}
          alt={friend.name}
          style={{
            width: '72px',
            height: '72px',
            borderRadius: '50%',
            objectFit: 'cover'
          }}
        />
        <div style={{ flex: 1 }}>
          <h3 style={{ margin: '0 0 6px 0' }}>{friend.name}</h3>
          {friend.about && (
            <p style={{ margin: 0, color: '#444' }}>{friend.about}</p>
          )}
        </div>
      </div>

      {/* Нижняя строка: имя/доп.инфа слева, кнопка по центру, дата/статус справа */}
      <div style={{
        display: 'flex',
        alignItems: 'center',
        minHeight: '40px',
        marginTop: '10px',
        gap: '10px',
        width: '100%'
      }}>
        <div style={{ flex: 1 }}>
          {friend.status || 'В друзьях'}
        </div>

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
          Удалить из друзей
        </button>

        <div style={{
          flex: 1,
          textAlign: 'right',
          color: '#666'
        }}>
          {friend.addedAt}
        </div>
      </div>
    </div>
  );
}

// Список друзей в виде списка «постов»
function FriendsList({ friends }) {
  return (
    <>
      {friends.map(friend => (
        <FriendCard key={friend.id} friend={friend} />
      ))}
    </>
  );
}

function Friends() {
  const pageTitle = "Друзья";
  const description = "Список твоих друзей";
  const friends = [
    {
      id: 1,
      name: "Иван",
      avatar: "https://lh3.googleusercontent.com/a/ACg8ocK7OJThoaQ-AVkdzQt5WZE6ayNaFeSfEl-6Dw3SnfndDiIPeuQk=s288-c-no",
      about: "Любит раст и котиков",
      status: "Лучший друг",
      addedAt: "21 ноября 2025"
    },
    {
      id: 2,
      name: "Радомир",
      avatar: "https://lh3.googleusercontent.com/a/ACg8ocK7OJThoaQ-AVkdzQt5WZE6ayNaFeSfEl-6Dw3SnfndDiIPeuQk=s288-c-no",
      about: "Го разработчик",
      status: "Приятель",
      addedAt: "20 ноября 2025"
    },
    {
      id: 3,
      name: "Ирина",
      avatar: "https://lh3.googleusercontent.com/a/ACg8ocK7OJThoaQ-AVkdzQt5WZE6ayNaFeSfEl-6Dw3SnfndDiIPeuQk=s288-c-no",
      about: "Фронтенд и дизайн",
      status: "Друг",
      addedAt: "18 ноября 2025"
    },
    {
      id: 4,
      name: "Александр",
      avatar: "https://lh3.googleusercontent.com/a/ACg8ocK7OJThoaQ-AVkdzQt5WZE6ayNaFeSfEl-6Dw3SnfndDiIPeuQk=s288-c-no",
      about: "Любитель игр и C++",
      status: "Друг",
      addedAt: "15 ноября 2025"
    }
  ];

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
            {pageTitle}
          </div>
          <div style={{ flex: 1, maxWidth: '340px' }}></div>
        </div>

        <div style={{
          width: '95vw',
          height: '100vh',
          flex: 1,
          display: 'flex',
          gap: '20px',
          padding: '20px',
          backgroundColor: '#f0f2f5',
          color: 'black'
        }}>
          <Sidebar title="Информация">
            <p style={{
              fontSize: '14px',
              lineHeight: '1.4',
              color: 'black'
            }}>
              {description}
            </p>
          </Sidebar>

          <div style={{
            flex: 1,
            overflowY: 'auto',
            scrollbarWidth: 'none',
            msOverflowStyle: 'none'
          }}>
            <FriendsList friends={friends} />
          </div>

          <Sidebar title="Все друзья">
            <FriendListCompact friends={friends} />
          </Sidebar>
        </div>
      </div>
    </>
  );
}

export default Friends;

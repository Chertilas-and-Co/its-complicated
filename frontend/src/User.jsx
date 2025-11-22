import React from 'react';
import Navbar from './Navbar';

function Search() {
  return (
    <div style={{
      width: '100%',
      maxWidth: '340px',
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
      width: '290px',
      padding: '15px',
      minHeight: '95vh',
      color: '#000',
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

function FriendList({ friends }) {
  return (
    <ul style={{ listStyle: 'none', padding: 0 }}>
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

function CommunitiesList({ communities }) {
  return (
    <ul style={{ listStyle: 'none', padding: 0, margin: 0 }}>
      {communities.map(community => (
        <li key={community.id} style={{
          marginBottom: '10px',
          display: 'flex',
          alignItems: 'center'
        }}>
          <img
            src={community.avatar}
            alt={community.name}
            style={{
              width: '36px',
              height: '36px',
              borderRadius: '50%',
              marginRight: '10px'
            }}
          />
          <span>{community.name}</span>
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
            backgroundColor: '#FFF',
            display: 'flex',
            flexDirection: 'column',
            color: 'black'
          }}
        >
          <div style={{ width:'500px', display: 'flex', gap: '15px', alignItems: 'flex-start' }}>
            <p style={{ margin: 0, flex: 1 }}>{post.content}</p>
          </div>
          <div>
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
          <div style={{
            display: 'flex',
            alignItems: 'center',
            minHeight: '40px',
            marginTop: '10px',
            gap: '10px',
            width: '100%'
          }}>
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

function UserPage() {
  const user = {
    avatar: "https://lh3.googleusercontent.com/a/ACg8ocK7OJThoaQ-AVkdzQt5WZE6ayNaFeSfEl-6Dw3SnfndDiIPeuQk=s288-c-no",
    name: "Абоба",
    description: "aboba",
  };

  const friends = [
    { id: 1, name: "Иван", avatar: "https://lh3.googleusercontent.com/a/ACg8ocK7OJThoaQ-AVkdzQt5WZE6ayNaFeSfEl-6Dw3SnfndDiIPeuQk=s288-c-no"},
    { id: 2, name: "Оля", avatar:"https://lh3.googleusercontent.com/a/ACg8ocK7OJThoaQ-AVkdzQt5WZE6ayNaFeSfEl-6Dw3SnfndDiIPeuQk=s288-c-no" }
  ];

  const communities = [
    { id: 1, name: "ReactJS", avatar: "https://lh3.googleusercontent.com/a/ACg8ocK7OJThoaQ-AVkdzQt5WZE6ayNaFeSfEl-6Dw3SnfndDiIPeuQk=s288-c-no" },
    { id: 2, name: "Frontend", avatar: "https://lh3.googleusercontent.com/a/ACg8ocK7OJThoaQ-AVkdzQt5WZE6ayNaFeSfEl-6Dw3SnfndDiIPeuQk=s288-c-no" }
  ];

  const posts = [
    { id: 1, content: "Мой первый пост", date: "20.11.2025", image: "" },
    { id: 2, content: "Второй пост с картинкой", date: "21.11.2025", image: "https://lh3.googleusercontent.com/a/ACg8ocK7OJThoaQ-AVkdzQt5WZE6ayNaFeSfEl-6Dw3SnfndDiIPeuQk=s288-c-no"},
    { id: 3, content: "Мой первый пост", date: "20.11.2025", image: "" },
    { id: 4, content: "Мой первый пост", date: "20.11.2025", image: "" },
    { id: 5, content: "Мой первый пост", date: "20.11.2025", image:"" },
    { id: 6, content: "Мой первый пост", date: "20.11.2025", image: "" }
  ];

  return (
    <>
    <Navbar/>
    <div style={{
      maxHeight: '100vh',
      width: '95vw',
      display: 'flex',
      flexDirection: 'row',
      justifyContent: 'center',
      gap: '24px',
      fontFamily: 'Arial, sans-serif',
      background: '#f5f5f5',
      padding: '30px 20px'
    }}>
     <div>
      <Search />
      <Sidebar title="Друзья">
        
        <FriendList  friends={friends} />
      </Sidebar>
      </div>
      <div style={{
        flex: 1,
        width: 'auto',
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        overflowY: 'auto',
      scrollbarWidth: 'none',
      msOverflowStyle: 'none'
      }}>
        
        <img
          src={user.avatar}
          alt={user.name}
          style={{
            width: '128px',
            height: '128px',
            borderRadius: '50%',
            objectFit: 'cover',
            marginBottom: '16px',
            border: '3px solid #ececec'
          }}
        />
        <h2 style={{ margin: 0, color: '#000' }}>{user.name}</h2>
        <p style={{
          color: '#000',
          margin: '8px 0 20px 0',
          textAlign: 'center'
        }}>{user.description}</p>
         
        <PostsList posts={posts} />
      
      </div>

      

      <Sidebar title="Сообщества">
        <CommunitiesList communities={communities} />
      </Sidebar>
    </div>
    </>
  );
  
}

export default UserPage;

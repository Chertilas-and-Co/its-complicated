import React from 'react';
import Navbar from './Navbar';

function DiscussionPost({ post }) {
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
      {/* Верхняя часть: текст и изображение */}
      <div style={{ display: 'flex', gap: '15px', alignItems: 'flex-start' }}>
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
              borderRadius: '6px'
            }}
          />
        )}
      </div>

      {/* Нижняя строка: Автор слева, кнопка по центру, дата справа */}
      <div
        style={{
          display: 'flex',
          alignItems: 'center',
          minHeight: '40px',
          marginTop: '10px',
          gap: '10px',
          width: '100%'
        }}
      >
        <div style={{ flex: 1 }}>{post.author}</div>
    
        <div
          style={{
            flex: 1,
            textAlign: 'right',
            color: '#666'
          }}
        >
          {post.date}
        </div>
      </div>
    </div>
  );
}

function CommentsList({ comments }) {
  return (
    <div style={{ width: '400px' }}>
      {comments.map((comment) => (
        <div
          key={comment.id}
          style={{
            border: '1px solid #ddd',
            padding: '14px',
            marginBottom: '17px',
            borderRadius: '6px',
            backgroundColor: '#fafafa'
          }}
        >

            <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
                <img
                    src={comment.avatar}
                    alt={`${comment.author} avatar`}
                    style={{
                    width: '40px',
                    height: '40px',
                    borderRadius: '50%',
                    objectFit: 'cover',
                    flexShrink: 0,
                    }}
                />
                <span>{comment.author}</span>
            </div>

          <div style={{ fontSize: '15px', marginBottom: '10px' }}>
            {comment.text}
          </div>
          <div
            style={{
              display: 'flex',
              justifyContent: 'space-between',
              fontSize: '13px',
              color: '#666'
            }}
          >
            
            <span style={{ marginLeft: 'auto' }}>{comment.date}</span>
          </div>
        </div>
      ))}
    </div>
  );
}

function DiscussionPage() {
  const post = {
    id: 0,
    content: 'Основной вопрос или тема для обсуждения группы.',
    author: 'Иван',
    date: '21 ноября 2025',
    avatar: 'https://lh3.googleusercontent.com/a/ACg8ocK7OJThoaQ-AVkdzQt5WZE6ayNaFeSfEl-6Dw3SnfndDiIPeuQk=s288-c-no'
  };

  const comments = [
    {
      id: 1,
      text: 'Первый комментарий — моё мнение.',
      author: 'Радомир',
      date: '21 ноября 2025',
      avatar: 'https://lh3.googleusercontent.com/a/ACg8ocK7OJThoaQ-AVkdzQt5WZE6ayNaFeSfEl-6Dw3SnfndDiIPeuQk=s288-c-no'
    },
    {
      id: 2,
      text: 'Абоба!',
      author: 'Ирина',
      date: '20 ноября 2025',
      avatar: 'https://lh3.googleusercontent.com/a/ACg8ocK7OJThoaQ-AVkdzQt5WZE6ayNaFeSfEl-6Dw3SnfndDiIPeuQk=s288-c-no'
    },
    {
      id: 3,
      text: 'Согласен с темой.',
      author: 'Александр',
      date: '19 ноября 2025',
      avatar: 'https://lh3.googleusercontent.com/a/ACg8ocK7OJThoaQ-AVkdzQt5WZE6ayNaFeSfEl-6Dw3SnfndDiIPeuQk=s288-c-no'
    },
    {
      id: 4,
      text: 'Согласен с темой.',
      author: 'Александр',
      date: '19 ноября 2025',
      avatar: 'https://lh3.googleusercontent.com/a/ACg8ocK7OJThoaQ-AVkdzQt5WZE6ayNaFeSfEl-6Dw3SnfndDiIPeuQk=s288-c-no'
    },
    {
      id: 5,
      text: 'Согласен с темой.',
      author: 'Александр',
      date: '19 ноября 2025',
      avatar: 'https://lh3.googleusercontent.com/a/ACg8ocK7OJThoaQ-AVkdzQt5WZE6ayNaFeSfEl-6Dw3SnfndDiIPeuQk=s288-c-no'
    },
    {
      id: 6,
      text: 'Согласен с темой.',
      author: 'Александр',
      date: '19 ноября 2025',
      avatar: 'https://lh3.googleusercontent.com/a/ACg8ocK7OJThoaQ-AVkdzQt5WZE6ayNaFeSfEl-6Dw3SnfndDiIPeuQk=s288-c-no'
    },
    {
      id: 7,
      text: 'Согласен с темой.',
      author: 'Александр',
      date: '19 ноября 2025',
      avatar: 'https://lh3.googleusercontent.com/a/ACg8ocK7OJThoaQ-AVkdzQt5WZE6ayNaFeSfEl-6Dw3SnfndDiIPeuQk=s288-c-no'
    },
    {
      id: 8,
      text: 'Согласен с темой.',
      author: 'Александр',
      date: '19 ноября 2025',
      avatar: 'https://lh3.googleusercontent.com/a/ACg8ocK7OJThoaQ-AVkdzQt5WZE6ayNaFeSfEl-6Dw3SnfndDiIPeuQk=s288-c-no'
    }
  ];

  return (
    <>
      <Navbar />
      <div
        style={{
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          minHeight: '100vh',
          paddingTop: '25px',
          fontFamily: 'Arial, sans-serif',
          background: '#f0f2f5'
        }}
      >
        
        <DiscussionPost post={post} />
        <div style={{
            flex: 1,
            overflowY: 'auto',
            scrollbarWidth: 'none',
            msOverflowStyle: 'none',
            paddingRight: '10px',
          }}>
             <CommentsList comments={comments} />
          </div>
        
      </div>
    </>
  );
}

export default DiscussionPage;

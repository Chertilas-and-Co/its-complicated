import React from 'react';
import Navbar from './Navbar';

function Search() {
  return (
    <div style={{
      width: '100%',
      maxWidth: '290px',
      display: 'flex',
      alignItems: 'center',
      gap: '10px',
    }}>
      <input
        type="text"
        placeholder="Поиск..."
        style={{
          flex: 1,
          padding: '8px 10px',
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

function ButtonCreate(){
  return(
        <a href="/create_community">
           <div style={{ flex: 1, maxWidth: '340px' , alignItems: 'center'}}><button style={{
              background: '#222',
              width: '200px',
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
              Создать сообщество
            </button></div>
            </a>
    );
}

function ButtonCommunity(){
    return(
        <a href="/community">
            <button  style={{
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
            Перейти
            </button>
            </a>
    );
}

function CommunitiesListCompact({ communities }) {
  return (
    <ul style={{ listStyle: 'none', padding: 0 }}>
      {communities.map(community => (
        <li key={community.id} style={{
          marginBottom: '12px',
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


function CommunityCard({ community }) {
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

      <div style={{ display: 'flex', gap: '15px', alignItems: 'flex-start' }}>
        <img
          src={community.avatar}
          alt={community.name}
          style={{
            width: '72px',
            height: '72px',
            borderRadius: '50%',
            objectFit: 'cover'
          }}
        />
        <div style={{ flex: 1 }}>
          <h3 style={{ margin: '0 0 6px 0' }}>{community.name}</h3>
          {community.description && (
            <p style={{ margin: 0, color: '#444' }}>{community.description}</p>
          )}
        </div>
        
      </div>

      {/* Нижняя строка*/}
      <div style={{
        display: 'flex',
        alignItems: 'center',
        minHeight: '40px',
        marginTop: '10px',
        gap: '10px',
        width: '100%'
      }}>
        <div style={{ flex: 1 }}>
          {community.membersCount ? `${community.membersCount} участников` : ''}
        </div>

        <ButtonCommunity/>

        <div style={{
          flex: 1,
          textAlign: 'right',
          color: '#666'
        }}>
          {community.createdAt}
        </div>
      </div>
    </div>
  );
}

function CommunitiesList({ communities }) {
  return (
    <>
      {communities.map(community => (
        <CommunityCard key={community.id} community={community} />
      ))}
    </>
  );
}

function CommunitiesPage() {
  const pageTitle = "Сообщества";
  const description = "Все сообщества, в которых вы участвуете";
  const communities = [
    
  ];

  return (
    <>
      <Navbar />
      <div style={{
        display: 'flex',
        flexDirection: 'column',
        height: '100vh',
        fontFamily: 'Arial, sans-serif',
        color: 'black',
        padding: '20px',
        gap: '20px',
        backgroundColor: '#f0f2f5'
      }}>
        {/* Шапка с поиском и заголовком */}
        <div style={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          paddingBottom: '10px',
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
         <ButtonCreate/>
        </div>

        <div style={{
          display: 'flex',
          flex: 1,
          gap: '20px',
          overflow: 'hidden'
        }}>
          {/* Левая колонка с описанием */}
          <Sidebar title="Информация">
            <p style={{ fontSize: '14px', lineHeight: '1.4', color: 'black' }}>
              {description}
            </p>
          </Sidebar>

          {/* Центр - список сообществ в карточках */}
          <div style={{
            flex: 1,
            overflowY: 'auto',
            scrollbarWidth: 'none',
            msOverflowStyle: 'none',
            paddingRight: '10px',
            color: 'black',
            fontFamily: 'Arial, sans-serif'
          }}>
            {communities.length > 0 ? (
              <CommunitiesList communities={communities} />
            ) : (
              <div style={{ padding: '20px', textAlign: 'center', color: '#888', fontSize: '18px' }}>
                Сообществ нет
              </div>
            )}
          </div>


          {/* Правая колонка с компактным списком */}
          <Sidebar title="Все сообщества">
            <CommunitiesListCompact communities={communities} />
          </Sidebar>
        </div>
      </div>
    </>
  );
}

export default CommunitiesPage;

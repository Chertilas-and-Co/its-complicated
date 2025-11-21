-- Осовная таблица пользователей
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash bytea NOT NULL,
    salt bytea NOT NULL,
    bio TEXT,
    avatar BYTEA,
    avatar_url VARCHAR(255)
);

-- Таблица друзей (двусторонние отношения)
CREATE TABLE friendships (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    friend_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(20) DEFAULT 'pending', -- 'pending', 'accepted', 'blocked'
    UNIQUE(user_id, friend_id),
    CHECK (user_id != friend_id) -- Нельзя добавить самого себя
);

-- Таблица подписок (одностороннее отношение)
CREATE TABLE communities (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    is_private BOOLEAN DEFAULT FALSE,
    created_by BIGINT NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
    
-- Таблица постов в сообществе
CREATE TABLE posts (
    id BIGSERIAL PRIMARY KEY,
    community_id BIGINT NOT NULL REFERENCES communities(id),
    text TEXT, 
    pic BYTEA,
    pic_url VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT NOT NULL REFERENCES users(id)
);

-- Таблица подписчиков сообщества
CREATE TABLE community_subscriptions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    community_id BIGINT NOT NULL REFERENCES communities(id) ON DELETE CASCADE
);

-- Таблица редакторов сообщества
CREATE TABLE community_writer (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    community_id BIGINT NOT NULL REFERENCES communities(id) ON DELETE CASCADE
);

-- Таблица админов сообщества
CREATE TABLE community_admin (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    community_id BIGINT NOT NULL REFERENCES communities(id) ON DELETE CASCADE
);

CREATE INDEX idx_friendships_user_id ON friendships(user_id);
CREATE INDEX idx_friendships_friend_id ON friendships(friend_id);
CREATE INDEX idx_friendships_status ON friendships(status);
CREATE INDEX idx_friendships_both ON friendships(user_id, friend_id);

CREATE INDEX idx_posts_community_id ON posts(community_id);
CREATE INDEX idx_posts_created_by ON posts(created_by);
CREATE INDEX idx_posts_community_created ON posts(community_id, created_at DESC); -- Для paginate
CREATE INDEX idx_posts_created_at ON posts(created_at DESC);

CREATE INDEX idx_community_subscriptions_user_id ON community_subscriptions(user_id);
CREATE INDEX idx_community_subscriptions_community_id ON community_subscriptions(community_id);
CREATE INDEX idx_community_subscriptions_both ON community_subscriptions(user_id, community_id);

CREATE INDEX idx_communities_created_by ON communities(created_by);


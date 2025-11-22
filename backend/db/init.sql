-- Осовная таблица пользователей
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    salt TEXT NOT NULL,
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
    title VARCHAR(255) NOT NULL,
    text TEXT NOT NULL,
    pic_url VARCHAR(500),
    community_id BIGINT NOT NULL,
    author_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE post_likes (
    post_id BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    PRIMARY KEY (post_id, user_id),
    UNIQUE (post_id, user_id)
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

CREATE INDEX idx_post_likes_user_id ON post_likes(user_id);

CREATE INDEX idx_posts_community_id ON posts(community_id);
CREATE INDEX idx_posts_author_id ON posts(author_id);
CREATE INDEX idx_posts_created_at ON posts(created_at DESC);

CREATE INDEX idx_community_subscriptions_user_id ON community_subscriptions(user_id);
CREATE INDEX idx_community_subscriptions_community_id ON community_subscriptions(community_id);
CREATE INDEX idx_community_subscriptions_both ON community_subscriptions(user_id, community_id);

CREATE INDEX idx_communities_created_by ON communities(created_by);

-- --- SEED DATA ---

-- Create 4 users
-- Passwords and salts are placeholders as they are not needed for this test
INSERT INTO users (id, username, email, password_hash, salt) VALUES
(1, 'user1', 'user1@example.com', 'hash1', 'salt1'),
(2, 'user2', 'user2@example.com', 'hash2', 'salt2'),
(3, 'user3', 'user3@example.com', 'hash3', 'salt3'),
(4, 'user4', 'user4@example.com', 'hash4', 'salt4')
ON CONFLICT (id) DO NOTHING;

-- Create 3 communities
INSERT INTO communities (id, name, description, created_by) VALUES
(1, 'Любители кошек', 'Обсуждаем наших пушистых друзей', 1),
(2, 'Любители собак', 'Все о лучших друзьях человека', 2),
(3, 'Фанаты научной фантастики', 'От Азимова до Желязны', 1)
ON CONFLICT (id) DO NOTHING;

-- Create subscriptions to generate intersections
INSERT INTO community_subscriptions (user_id, community_id) VALUES
-- Cat Lovers (size: 2)
(1, 1),
(2, 1),
-- Dog Lovers (size: 3)
(2, 2),
(3, 2),
(4, 2),
-- Sci-Fi Fans (size: 3)
(1, 3),
(2, 3),
(3, 3)
ON CONFLICT DO NOTHING;

-- Reset sequence for correct auto-incrementing IDs if table was not empty
SELECT setval('users_id_seq', (SELECT MAX(id) FROM users));
SELECT setval('communities_id_seq', (SELECT MAX(id) FROM communities));
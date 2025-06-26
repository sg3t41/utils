-- リンク集テーブルの作成
CREATE TABLE IF NOT EXISTS links (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,                   -- リンクタイトル
    url TEXT NOT NULL,                             -- リンクURL
    description TEXT,                              -- 説明文
    platform VARCHAR(50) NOT NULL,                -- プラットフォーム名（twitter, instagram, github等）
    icon_name VARCHAR(50),                         -- アイコン名
    background_color VARCHAR(7) DEFAULT '#6B7280', -- 背景色（HEXコード）
    text_color VARCHAR(7) DEFAULT '#FFFFFF',       -- テキスト色
    order_index INTEGER DEFAULT 0,                -- 表示順序
    is_active BOOLEAN DEFAULT true,               -- 有効フラグ
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- インデックス作成
CREATE INDEX idx_links_user_id ON links(user_id);
CREATE INDEX idx_links_order_index ON links(order_index);
CREATE INDEX idx_links_platform ON links(platform);
CREATE INDEX idx_links_is_active ON links(is_active);

-- updated_atの自動更新用トリガー
CREATE OR REPLACE FUNCTION update_links_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER trigger_update_links_updated_at
    BEFORE UPDATE ON links
    FOR EACH ROW
    EXECUTE FUNCTION update_links_updated_at();

-- サンプルデータの挿入
INSERT INTO links (title, url, description, platform, icon_name, background_color, text_color, order_index, user_id) VALUES
('X (Twitter)', 'https://twitter.com/sg3t41', 'Follow me on X', 'twitter', 'twitter', '#1DA1F2', '#FFFFFF', 1, 1),
('Instagram', 'https://instagram.com/sg3t41', 'My Instagram profile', 'instagram', 'instagram', '#E4405F', '#FFFFFF', 2, 1),
('GitHub', 'https://github.com/sg3t41', 'Check out my repositories', 'github', 'github', '#333333', '#FFFFFF', 3, 1),
('LINE', 'https://line.me/ti/p/sg3t41', 'Add me on LINE', 'line', 'line', '#00C300', '#FFFFFF', 4, 1);
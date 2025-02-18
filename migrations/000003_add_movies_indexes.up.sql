CREATE INDEX IF NOT EXISTS idx_movies_title_search ON movies USING GIN (to_tsvector('simple', title));
CREATE INDEX IF NOT EXISTS idx_movies_genres ON movies USING GIN (genres);

-- Create pages table
CREATE TABLE IF NOT EXISTS pages (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create page_links table for tracking relationships
CREATE TABLE IF NOT EXISTS page_links (
    id SERIAL PRIMARY KEY,
    from_page_id INTEGER REFERENCES pages(id) ON DELETE CASCADE,
    to_page_id INTEGER REFERENCES pages(id) ON DELETE CASCADE,
    link_text VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(from_page_id, to_page_id)
);

CREATE INDEX IF NOT EXISTS idx_page_links_from ON page_links(from_page_id);
CREATE INDEX IF NOT EXISTS idx_page_links_to ON page_links(to_page_id);

-- Insert default Home page
INSERT INTO pages (name, content, created_at, updated_at)
VALUES (
    'Home',
    '<h1>Welcome to Wiki Notes</h1><p>This is your personal wiki-style note-taking application.</p><h2>Features</h2><ul><li><p><strong>Markdown Support</strong>: Write your notes in markdown format</p></li><li><p><strong>LaTeX Equations</strong>: Use inline math like $E = mc^2$ or display equations like $$\int_0^\infty e^{-x^2} dx = \frac{\sqrt{\pi}}{2}$$</p></li><li><p><strong>Wiki Links</strong>: Create links to other pages using [[Page Name]] syntax</p></li><li><p><strong>Search</strong>: Use the search bar at the top to find pages quickly</p></li></ul><h2>Getting Started</h2><p>Click on any [[Page Name]] to create a new page, or use the search bar to explore existing pages.</p><p>Happy note-taking!</p>',
    NOW(),
    NOW()
)
ON CONFLICT (name) DO NOTHING;


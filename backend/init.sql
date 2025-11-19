-- Create pages table
CREATE TABLE IF NOT EXISTS pages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) UNIQUE NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Insert default Home page
INSERT INTO pages (id, name, content, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    'Home',
    '<h1>Welcome to Wiki Notes</h1><p>This is your personal wiki-style note-taking application.</p><h2>Features</h2><ul><li><p><strong>Markdown Support</strong>: Write your notes in markdown format</p></li><li><p><strong>LaTeX Equations</strong>: Use inline math like $E = mc^2$ or display equations like $$\int_0^\infty e^{-x^2} dx = \frac{\sqrt{\pi}}{2}$$</p></li><li><p><strong>Wiki Links</strong>: Create links to other pages using [[Page Name]] syntax</p></li><li><p><strong>Search</strong>: Use the search bar at the top to find pages quickly</p></li></ul><h2>Getting Started</h2><p>Click on any [[Page Name]] to create a new page, or use the search bar to explore existing pages.</p><p>Happy note-taking!</p>',
    NOW(),
    NOW()
)
ON CONFLICT (name) DO NOTHING;


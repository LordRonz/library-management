CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE books (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    year INT NOT NULL,
    description TEXT,
    isbn VARCHAR(20),
    genre VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO books (title, author, year, description, isbn, genre) VALUES
('The Great Gatsby', 'F. Scott Fitzgerald', 1925, 'The story of the fabulously wealthy Jay Gatsby and his love for the beautiful Daisy Buchanan.', '978-0743273565', 'Tragedy'),
('To Kill a Mockingbird', 'Harper Lee', 1960, 'The story of a young girl, Scout Finch, and her lawyer father, Atticus, in the American South.', '978-0061120084', 'Southern Gothic'),
('1984', 'George Orwell', 1949, 'A dystopian novel set in Airstrip One, a province of the superstate Oceania in a world of perpetual war.', '978-0451524935', 'Dystopian');

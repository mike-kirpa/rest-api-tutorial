CREATE TABLE public.author (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name VARCHAR(100) NOT NULL
);

CREATE TABLE public.book (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name VARCHAR(100) NOT NULL,
	author_id UUID NOT NULL,
	CONSTRAINT author_fk FOREIGN KEY (author_id) REFERENCES public.author(id)
);

INSERT INTO author (name) VALUES ('Народ');
INSERT INTO author (name) VALUES ('Джек Лондон');
INSERT INTO author (name) VALUES ('Джоан Роулинг');

INSERT INTO book (name, author_id) VALUES ('колобок', '92eef194-c3b9-45c6-a8e2-fc17b99d2648');
INSERT INTO book (name, author_id) VALUES ('гарри поттер', 'e20862e3-b234-4d08-bdbb-9e2f2016a226');
INSERT INTO book (name, author_id) VALUES ('белый клык', 'abe3e1ed-18a2-4e9a-9dae-bf36453199b8');

//many to many

DROP TABLE IF EXISTS author CASCADE;
DROP TABLE IF EXISTS book CASCADE;
DROP TABLE IF EXISTS book_authors CASCADE;

CREATE TABLE public.author (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name VARCHAR(100) NOT NULL
);

CREATE TABLE public.book (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name VARCHAR(100) NOT NULL,
	age INT
);

CREATE TABLE public.book_authors (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	book_id UUID NOT NULL,
	author_id UUID NOT NULL,
	CONSTRAINT book_fk FOREIGN KEY (book_id) REFERENCES public.book(id),
	CONSTRAINT author_fk FOREIGN KEY (author_id) REFERENCES public.author(id),
	CONSTRAINT book_authors_unique UNIQUE (book_id, author_id)
);

INSERT INTO author (id, name) VALUES ('92eef194-c3b9-45c6-a8e2-fc17b99d2648', 'Народ');
INSERT INTO author (id, name) VALUES ('e20862e3-b234-4d08-bdbb-9e2f2016a226', 'Джек Лондон');
INSERT INTO author (id, name) VALUES ('abe3e1ed-18a2-4e9a-9dae-bf36453199b8', 'Джоан Роулинг');

INSERT INTO book (id, name, age) VALUES ('da8c677b-37e5-42a2-99c9-f02829105f77', 'колобок', 1000);
INSERT INTO book (id, name, age) VALUES ('2b6ece6a-cd3d-4871-91fa-167839d9281e', 'гарри поттер', 22);
INSERT INTO book (id, name) VALUES ('bbcbab77-d6d1-41dd-9f8d-aa0cc4234879', 'белый клык');

INSERT INTO book_authors (book_id, author_id) VALUES ('da8c677b-37e5-42a2-99c9-f02829105f77', '92eef194-c3b9-45c6-a8e2-fc17b99d2648');
INSERT INTO book_authors (book_id, author_id) VALUES ('da8c677b-37e5-42a2-99c9-f02829105f77', 'abe3e1ed-18a2-4e9a-9dae-bf36453199b8');

INSERT INTO book_authors (book_id, author_id) VALUES ('2b6ece6a-cd3d-4871-91fa-167839d9281e', '92eef194-c3b9-45c6-a8e2-fc17b99d2648');
INSERT INTO book_authors (book_id, author_id) VALUES ('2b6ece6a-cd3d-4871-91fa-167839d9281e', 'abe3e1ed-18a2-4e9a-9dae-bf36453199b8');

SELECT 
	b.id, b.name,
	array((SELECT ba.author_id FROM book_authors ba WHER ba.book_id = b.id)) AS authors
FROM book b
GROUP BY b.id, b.name 


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
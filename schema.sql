CREATE TABLE IF NOT EXISTS "pastes" (
	"id"	TEXT,
	"title"	TEXT,
	"content"	TEXT NOT NULL,
	"password"	TEXT,
	"expiry"	TEXT NOT NULL,
	"created"	TEXT NOT NULL,
	"burn"	INTEGER NOT NULL DEFAULT 0,
	"discuss"	INTEGER NOT NULL DEFAULT 0,
	PRIMARY KEY("id")
);
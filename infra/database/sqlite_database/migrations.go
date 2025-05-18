package sqlite_database

type migration struct {
	description string
	sql string
}

func getMigrations() []migration {
	return []migration{
		{
			description: "Create projects table",
			sql: `
			CREATE TABLE IF NOT EXISTS projects (
				uuid TEXT PRIMARY KEY NOT NULL,
				name TEXT NOT NULL UNIQUE,
				path TEXT NOT NULL
			);
			`,
		},
	}
}


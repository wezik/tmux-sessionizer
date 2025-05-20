package sqlite_database

import (
	"database/sql"
	"fmt"
	"os"
	"phopper/domain/errors"
	"phopper/domain/project"
	"phopper/domain/project/session_template"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type SqliteProjectRepository struct{}

func (_ SqliteProjectRepository) GetProjects() []project.Project {
	db := openDB()
	defer db.Close()

	rows, err := db.Query(`SELECT uuid, name, path FROM projects;`)
	errors.EnsureNotNil(err, "Could not query projects")

	defer rows.Close()

	projects := []project.Project{}
	for rows.Next() {
		var uuid string
		var name string
		var path string

		err = rows.Scan(&uuid, &name, &path)
		errors.EnsureNotNil(err, "Could not scan project")

		project := project.Project{
			UUID:    uuid,
			Session: session_template.SessionTemplate{Name: name, Path: path},
		}
		projects = append(projects, project)
	}

	return projects
}

func (_ SqliteProjectRepository) SaveProject(project project.Project) project.Project {
	db := openDB()
	defer db.Close()

	if project.UUID == "" {
		project.UUID = uuid.New().String()
	}

	project.Session.Name = renameIfExists(db, project.Session.Name)

	sql := `INSERT INTO projects (uuid, name, path) VALUES (?, ?, ?);`

	_, err := db.Exec(sql, project.UUID, project.Session.Name, project.Session.Path)
	errors.EnsureNotNil(err, "Could not save project")

	return project
}

// util function that updates the name with order numbers if identical one already exist
func renameIfExists(db *sql.DB, name string) string {
	var order = 0
	temp := name
	for true {
		var count int
		err := db.QueryRow(`SELECT COUNT(*) FROM projects WHERE name = ?;`, temp).Scan(&count)
		errors.EnsureNotNil(err, "Could not check if project name exists")

		if count == 0 {
			break
		}

		order++
		temp = fmt.Sprintf("%s(%d)", name, order)
	}
	return temp
}

func (_ SqliteProjectRepository) DeleteProject(uuid string) {
	db := openDB()
	defer db.Close()

	_, err := db.Exec(`DELETE FROM projects WHERE uuid = ?;`, uuid)
	errors.EnsureNotNil(err, "Could not delete project")
}

func (_ SqliteProjectRepository) PrepareTemplateFilePath(p project.Project) string {
	// TODO: this should prepare a temp editable file that will be deleted when the editor is closed
	// this is a bit tricky and to be figured out for now I am focusing on the yaml storage
	fmt.Println("This part of the code is not yet implemented")
	os.Exit(1)
	return ""
}

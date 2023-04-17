package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upCreateUsers, downCreateUsers)
}

func upCreateUsers(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE users
		(
			id               varbinary(255) NOT NULL,
			organization_id  varbinary(255) NOT NULL,
			level            varchar(255)   NOT NULL,
			name             varchar(255)   NOT NULL,
			email            varchar(255)   NOT NULL,
			created_at       datetime       NOT NULL,
			updated_at       datetime       DEFAULT NULL,
			deleted_at       datetime       DEFAULT NULL,
			PRIMARY KEY (id),
			CONSTRAINT users_organization_id_fk FOREIGN KEY (organization_id) REFERENCES organizations (id)
		) ENGINE = InnoDB
			DEFAULT CHARSET = utf8mb4
			COLLATE = utf8mb4_0900_ai_ci
	`)
	if err != nil {
		return err
	}
	return nil
}

func downCreateUsers(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE users`)
	if err != nil {
		return err
	}
	return nil
}

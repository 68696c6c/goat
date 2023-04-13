package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upCreateOrganizations, downCreateOrganizations)
}

func upCreateOrganizations(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE organizations
		(
			id         varbinary(255) NOT NULL,
			name       varchar(255)   NOT NULL,
			website    varchar(255)   NOT NULL,
			created_at datetime       NOT NULL,
			updated_at datetime       DEFAULT NULL,
			deleted_at datetime       DEFAULT NULL,
			PRIMARY KEY (id)
		) ENGINE = InnoDB
			DEFAULT CHARSET = utf8mb4
			COLLATE = utf8mb4_0900_ai_ci
	`)
	if err != nil {
		return err
	}
	return nil
}

func downCreateOrganizations(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE organizations`)
	if err != nil {
		return err
	}
	return nil
}

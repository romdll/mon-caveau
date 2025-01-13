package migrations

import (
	"database/sql"
	"fmt"
	"moncaveau/database/crypt"
)

type CustomMigrationFunc func(db *sql.DB) error

type Migration struct {
	Version         float64
	SQL             string
	CustomMigration CustomMigrationFunc
}

func ApplyMigrations(db *sql.DB) error {
	migrations := []Migration{
		{
			SQL: `
				CREATE TABLE accounts (
					id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
					account_key VARCHAR(255) NOT NULL UNIQUE,
					email VARCHAR(255) UNIQUE,
					password VARCHAR(255),
					name VARCHAR(255),
					surname VARCHAR(255),
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
				);
			`,
			Version: 1.0,
		},
		{
			SQL: `
				CREATE TABLE sessions (
					id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
					account_id INT NOT NULL,
					session_token VARCHAR(255) UNIQUE NOT NULL,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					expires_at TIMESTAMP,
					last_activity TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					FOREIGN KEY (account_id) REFERENCES accounts(id)
				);
			`,
			Version: 2.0,
		},
		{
			SQL: `
				CREATE TABLE wine_domains (
					id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
					name VARCHAR(255) NOT NULL UNIQUE 
				);
			`,
			Version: 3.1,
		},
		{
			SQL: `
				CREATE TABLE wine_regions (
					id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
					name VARCHAR(255) NOT NULL,
					country VARCHAR(255) NOT NULL,
					CONSTRAINT unique_name_country UNIQUE (name, country)
				);
			`,
			Version: 3.2,
		},
		{
			SQL: `
				CREATE TABLE wine_types (
					id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
					name VARCHAR(255) NOT NULL UNIQUE
				);
			`,
			Version: 3.3,
		},
		{
			SQL: `
				CREATE TABLE wine_bottle_sizes (
					id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
					size REAL NOT NULL UNIQUE,
					name REAL NOT NULL UNIQUE
				);
			`,
			Version: 3.4,
		},
		{
			SQL: `
				CREATE TABLE wine_wines (
					id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
					name VARCHAR(255) NOT NULL,

					domaine_id INTEGER NOT NULL,
					region_id INTEGER NOT NULL,
					type_id INTEGER NOT NULL,
					bottle_size_id INTEGER NOT NULL,

					vintage INTEGER NOT NULL,
					qantity INTEGER NOT NULL,

					buy_price REAL,
					description VARCHAR(3000),
					image LONGTEXT,

					account_id INT NOT NULL,

					FOREIGN KEY (account_id) REFERENCES accounts(id),
					FOREIGN KEY (domaine_id) REFERENCES wine_domains(id),
					FOREIGN KEY (region_id) REFERENCES wine_regions(id),
					FOREIGN KEY (type_id) REFERENCES wine_types(id),
					FOREIGN KEY (bottle_size_id) REFERENCES wine_bottle_sizes(id)
				);
			`,
			Version: 3.5,
		},
		{
			SQL: `
				CREATE TABLE wine_transactions (
					id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
					wine_id INTEGER NOT NULL,
					quantity INTEGER NOT NULL,
					type VARCHAR(255) NOT NULL,
					date DATETIME DEFAULT CURRENT_TIMESTAMP,
					FOREIGN KEY (wine_id) REFERENCES wine_wines(id)
				);
			`,
			Version: 3.6,
		},
		{
			SQL: `
				ALTER TABLE wine_bottle_sizes
				MODIFY COLUMN name VARCHAR(255) NOT NULL UNIQUE;
			`,
			Version: 4.1,
		},
		{
			SQL: `
				ALTER TABLE wine_wines
				CHANGE COLUMN domaine_id domain_id INTEGER NOT NULL;
			`,
			Version: 4.2,
		},
		{
			SQL: `
				ALTER TABLE wine_wines
				CHANGE COLUMN qantity quantity INTEGER NOT NULL;
			`,
			Version: 4.3,
		},
		{
			SQL: `
				ALTER TABLE wine_wines
				MODIFY COLUMN buy_price REAL DEFAULT 0;
			`,
			Version: 4.4,
		},
		{
			SQL: `
				ALTER TABLE wine_wines
				MODIFY COLUMN description VARCHAR(3000) DEFAULT '';
			`,
			Version: 4.5,
		},
		{
			SQL: `
				CREATE TRIGGER before_insert_wine_image
				BEFORE INSERT ON wine_wines
				FOR EACH ROW
				BEGIN
					IF NEW.image IS NULL THEN
						SET NEW.image = '/v1/images/logo.png';
					END IF;
				END;
			`,
			Version: 4.6,
		},
		{
			SQL: `
				CREATE TRIGGER after_insert_wine
				AFTER INSERT ON wine_wines
				FOR EACH ROW
				BEGIN
					INSERT INTO wine_transactions (wine_id, quantity, type, date)
					VALUES (NEW.id, NEW.quantity, 'added', UTC_TIMESTAMP(6));
				END;
			`,
			Version: 5.1,
		},
		{
			SQL: `
				CREATE TRIGGER after_update_wine_quantity
				AFTER UPDATE ON wine_wines
				FOR EACH ROW
				BEGIN
					IF NEW.quantity > OLD.quantity THEN
						INSERT INTO wine_transactions (wine_id, quantity, type, date)
						VALUES (NEW.id, NEW.quantity - OLD.quantity, 'added', UTC_TIMESTAMP(6));
					END IF;

					IF NEW.quantity < OLD.quantity THEN
						INSERT INTO wine_transactions (wine_id, quantity, type, date)
						VALUES (NEW.id, OLD.quantity - NEW.quantity, 'drank', UTC_TIMESTAMP(6)); 
					END IF;
				END;
			`,
			Version: 5.2,
		},
		{
			SQL: `
				ALTER TABLE wine_transactions
				MODIFY COLUMN date DATETIME(6) NOT NULL;
			`,
			Version: 5.3,
		},
		{
			SQL:     `DROP TRIGGER before_insert_wine_image;`,
			Version: 5.4,
		},
		{
			SQL: `
				CREATE TRIGGER before_insert_wine_image
				BEFORE INSERT ON wine_wines
				FOR EACH ROW
				BEGIN
					IF NEW.image IS NULL THEN
						SET NEW.image = '/v1/images/no_photo_generic.svg';
					END IF;
				END;
			`,
			Version: 5.5,
		},
		{
			SQL: `
				ALTER TABLE accounts
				DROP INDEX email;
			`,
			Version: 6.1,
		},
		{
			SQL: `
				ALTER TABLE accounts
				MODIFY email VARCHAR(255) DEFAULT '',
				MODIFY password VARCHAR(255) DEFAULT '',
				MODIFY name VARCHAR(255) DEFAULT '',
				MODIFY surname VARCHAR(255) DEFAULT '';
			`,
			Version: 6.2,
		},
		{
			SQL: `
				UPDATE accounts
				SET 
					email = COALESCE(email, ''),
					password = COALESCE(password, ''),
					name = COALESCE(name, ''),
					surname = COALESCE(surname, '');
			`,
			Version: 6.3,
		},
		{
			CustomMigration: crypt.CustomHashAccountKeysMigration,
			Version:         7.0,
		},
		{
			SQL: `
				ALTER TABLE wine_wines
				ADD COLUMN preferred_start_date DATE NULL,
				ADD COLUMN preferred_end_date DATE NULL;
			`,
			Version: 8.0,
		},
	}

	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (version FLOAT PRIMARY KEY)`)
	if err != nil {
		return fmt.Errorf("failed to ensure migrations table exists: %w", err)
	}

	var currentVersion float64
	err = db.QueryRow(`SELECT IFNULL(MAX(version), 0) FROM schema_migrations`).Scan(&currentVersion)
	if err != nil {
		return fmt.Errorf("failed to fetch current schema version: %w", err)
	}

	logger.Infof("Current schema version: %.01f", currentVersion)

	var migrationsToApply []Migration
	for _, migration := range migrations {
		if migration.Version > currentVersion {
			migrationsToApply = append(migrationsToApply, migration)
		}
	}

	if len(migrationsToApply) == 0 {
		logger.Info("No migrations to apply")
		return nil
	}

	logger.Infof("Found %d migrations to apply", len(migrationsToApply))

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	migrationApplied := false
	for _, migration := range migrationsToApply {
		migrationApplied = true
		logger.Infof("Applying migration version %.01f", migration.Version)

		if migration.CustomMigration != nil {
			logger.Infof("Executing custom function for migration version %.01f", migration.Version)
			if err := migration.CustomMigration(db); err != nil {
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					return fmt.Errorf("failed to apply migration %.01f, and rollback failed: %w", migration.Version, rollbackErr)
				}
				return fmt.Errorf("failed to apply custom migration %.01f: %w", migration.Version, err)
			}
		} else {
			logger.Infof("Executing SQL for migration version %.01f", migration.Version)
			if _, err := tx.Exec(migration.SQL); err != nil {
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					return fmt.Errorf("failed to apply migration %.01f, and rollback failed: %w", migration.Version, rollbackErr)
				}
				return fmt.Errorf("failed to apply migration %.01f: %w", migration.Version, err)
			}
		}

		if _, err := tx.Exec(`INSERT INTO schema_migrations (version) VALUES (?)`, migration.Version); err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return fmt.Errorf("failed to record migration %.01f, and rollback failed: %w", migration.Version, rollbackErr)
			}
			return fmt.Errorf("failed to record migration %.01f: %w", migration.Version, err)
		}

		logger.Infof("Migration version %.01f applied successfully", migration.Version)
	}

	if migrationApplied {
		err = tx.Commit()
		if err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}
		logger.Info("All migrations applied successfully")
	}
	return nil

}

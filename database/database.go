package database

import (
	"database/sql"
	"fmt"
)

const (
	driver   = "mysql"
	user     = "root"
	password = ""
	protocol = "tcp"
	address  = "127.0.0.1"
	port     = "3306"
)

// DB is the database connection handler.
var DB *sql.DB

// Init intializes the database.
func Init(env string) (*sql.DB, error) {
	db, err := sql.Open(driver, fmt.Sprintf("%v:%v@%v(%v:%v)/", user, password, protocol, address, port))
	if err != nil {
		return nil, err
	}

	if env == "test" {
		_, err = db.Exec(
			"CREATE DATABASE IF NOT EXISTS trapAdvisor_test DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci",
		)
		if err != nil {
			return nil, err
		}

		_, err = db.Exec("USE trapAdvisor_test")
		if err != nil {
			return nil, err
		}
	} else {
		_, err = db.Exec(
			"CREATE DATABASE IF NOT EXISTS trapAdvisor DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci",
		)
		if err != nil {
			return nil, err
		}

		_, err = db.Exec("USE trapAdvisor")
		if err != nil {
			return nil, err
		}
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS traveller (
			id BIGINT UNSIGNED NOT NULL,
			name VARCHAR(250) NOT NULL,
			PRIMARY KEY (id),
			UNIQUE INDEX traveller_id_unique (id ASC)
		) ENGINE = InnoDB
	`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS friendship (
			travellerId BIGINT UNSIGNED NOT NULL,
			friendId BIGINT UNSIGNED NOT NULL,
			similarity FLOAT NULL DEFAULT NULL,
			PRIMARY KEY (travellerId, friendId),
			INDEX fk_friendship_traveller_idx (travellerId ASC),
			INDEX fk_friendship_friend_idx (friendId ASC),
			CONSTRAINT fk_friendship_traveller
				FOREIGN KEY (travellerId)
				REFERENCES traveller (id)
				ON DELETE NO ACTION
				ON UPDATE NO ACTION,
			CONSTRAINT fk_friendship_friend
				FOREIGN KEY (friendId)
				REFERENCES traveller (id)
				ON DELETE NO ACTION
				ON UPDATE NO ACTION
		) ENGINE = InnoDB
	`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS trip (
			id INT NOT NULL AUTO_INCREMENT,
			name VARCHAR(250) NOT NULL,
			startDate DATE NOT NULL,
			endDate DATE NOT NULL,
			rating FLOAT NOT NULL,
			review TEXT NULL,
			travellerId BIGINT UNSIGNED NOT NULL,
			PRIMARY KEY (id),
			INDEX fk_trip_traveller_idx (travellerId ASC),
			CONSTRAINT fk_trip_traveller
				FOREIGN KEY (travellerId)
				REFERENCES traveller (id)
				ON DELETE NO ACTION
				ON UPDATE NO ACTION
		) ENGINE = InnoDB
	`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS touristAttraction (
			id INT NOT NULL AUTO_INCREMENT,
			name VARCHAR(250) NOT NULL,
			location VARCHAR(250) NOT NULL,
			visitDate DATE NOT NULL,
			rating FLOAT NOT NULL DEFAULT 0,
			pros TEXT NULL,
			cons TEXT NULL,
			tripId INT NOT NULL,
			PRIMARY KEY (id),
			INDEX fk_touristAttraction_trip_idx (tripId ASC),
			CONSTRAINT fk_touristAttraction_trip
				FOREIGN KEY (tripId)
				REFERENCES trip (id)
				ON DELETE NO ACTION
				ON UPDATE NO ACTION
		) ENGINE = InnoDB
	`)
	if err != nil {
		return nil, err
	}

	if env == "test" {
		_, err = db.Exec("INSERT INTO traveller VALUES (123, 'Morisquinho')")
		if err != nil {
			return nil, err
		}

		_, err = db.Exec("INSERT INTO traveller VALUES (1234, 'Ludwig')")
		if err != nil {
			return nil, err
		}

		_, err = db.Exec(
			"INSERT INTO trip VALUES (1, 'Roland Garros', '2017-05-30', '2017-05-30', 9, 'legal', 123)",
		)
		if err != nil {
			return nil, err
		}

		_, err = db.Exec(`
			INSERT INTO touristAttraction
			VALUES (1, 'Torre Eiffel', 'Paris', '2017-05-30', 7.5, 'Não fica no Brasil', 'Fica na França', 1)
		`)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

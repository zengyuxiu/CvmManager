package db_api

import (
	"database/sql"
	"log"
	"os"
)

func Init_DB() {
	_ = os.Remove("stats-database.db")
	// SQLite is a file based database.
	file, err := os.Create("stats-database.db") // Create SQLite file
	if err != nil {
		log.Fatal(err)
	}
	err = file.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func CreateTable(db *sql.DB) {
	createDockerTable(db)
	createKvmTable(db)
	createLxcTable(db)
}

func createDockerTable(db *sql.DB) {
	CreateDockerTable := `create table docker_stats
	(
			ID                text
				constraint docker_stats_pk
					primary key,
			memory_usage      integer,
			memory_max_usage  integer,
			memory_limit      integer,
			cpu_usage         integer,
			cpu_online        integer,
			cpu_usage_percent float,
			storage_write     integer,
			storage_read      integer,
			time_stamp        timestamp
				);` // SQL Statement for Create Table

	log.Println("Create Docker table...")
	statement, err := db.Prepare(CreateDockerTable) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = statement.Exec() // Execute SQL Statements
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("Docker table created")
}

func createKvmTable(db *sql.DB) {
	CreateKvmTable := `create table kvm_stats
(
    Name             text
        constraint kvm_stats_pk
            primary key,
    memory_usage     int,
    memory_max_usage int,
    cpu_usage        int,
    cpu_online       int,
    time_stamp       timestamp
);
` // SQL Statement for Create Table

	log.Println("Create Kvm table...")
	statement, err := db.Prepare(CreateKvmTable) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = statement.Exec() // Execute SQL Statements
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("Kvm table created")
}

func createLxcTable(db *sql.DB) {
	CreateLxcTable := `create table lxc_stats
		(
			name             text not null
				constraint lxc_stats_pk
					primary key,
			cpu_usage        integer,
			memory_usage     integer,
			memory_max_usage integer,
			time_stamp       timestamp
		);` // SQL Statement for Create Table

	log.Println("Create Lxc table...")
	statement, err := db.Prepare(CreateLxcTable) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = statement.Exec() // Execute SQL Statements
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("Lxc table created")
}

package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// CREATE TABLE `userinfo` (
//		`uid` INT(10) NOT NULL AUTO_INCREMENT,
//      `username` VARCHAR(64) NULL DEFAULT NULL,
//      `department` VARCHAR(64) NULL DEFAULT NULL,
//		`created` DATE NULL DEFAULT NULL,
//		PRIMARY KEY (`uid`)
// );

// CREATE TABLE `userdetail` (
//		`uid` INT(10) NOT NULL DEFAULT '0',
//		`intro` TEXT NULL,
//		`profile` TEXT NULL,
//		PRIMARY KEY (`uid`)
// );

// Result ...
type Result struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Default sql.NullString
	Extra   string
}

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test?charset=utf8")
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxIdleConns(100)
	db.SetMaxOpenConns(600)

	res, err := db.Exec(`DROP TABLE IF EXISTS userinfo;`)
	if err != nil {
		log.Fatal(err)
	}
	res, err = db.Exec("DROP TABLE IF EXISTS userdetail;")
	if err != nil {
		log.Fatal(err)
	}

	res, err = db.Exec(`CREATE TABLE IF NOT EXISTS userinfo(
		uid INT(10) NOT NULL AUTO_INCREMENT,
		username VARCHAR(64) NULL DEFAULT NULL,
		department VARCHAR(64) NULL DEFAULT NULL,
		created DATE NULL DEFAULT NULL,
		PRIMARY KEY(` + "`uid`" + `)
	) ENGINE=InnoDB charset=utf8mb4;`)
	if err != nil {
		log.Fatal(err)
	}

	res, err = db.Exec(`CREATE TABLE IF NOT EXISTS userdetail(
		uid INT(10) NOT NULL DEFAULT '0',
		intro TEXT NULL,
		profile TEXT NULL,
		PRIMARY KEY(uid)
	) ENGINE=InnoDB charset=utf8mb4`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)

	rows, err := db.Query("SHOW TABLES")
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var table string
		err = rows.Scan(&table)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(table)
	}

	rows, err = db.Query("SHOW CREATE TABLE userinfo")
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var tablename string
		var createSQL string
		err = rows.Scan(&tablename, &createSQL)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("SHOW CREATE TABLE: %s, (%s)\n", tablename, createSQL)
	}

	rows, err = db.Query("DESCRIBE userinfo")
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var result Result
		err = rows.Scan(&result.Field, &result.Type, &result.Null, &result.Key, &(result.Default), &result.Extra)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result)
	}

	stmt, err := db.Prepare("INSERT userinfo SET username=?,department=?,created=?")
	defer stmt.Close()
	if err != nil {
		log.Fatal(err)
	}

	res, err = stmt.Exec("test", "test", "2012-01-01")
	if err != nil {
		log.Fatal(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(id)

	stmt, err = db.Prepare("UPDATE userinfo SET username=? where uid=?")
	defer stmt.Close()
	if err != nil {
		log.Fatal(err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(affected)

	rows, err = db.Query("SELECT * FROM userinfo")
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var uid int
		var username string
		var department string
		var created string
		err = rows.Scan(&uid, &username, &department, &created)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("uid=%d, username=%s, department=%s, created=%s\n", uid, username, department, created)
	}

	stmt, err = db.Prepare("DELETE FROM userinfo WHERE uid=?")
	defer stmt.Close()
	if err != nil {
		log.Fatal(err)
	}

	res, err = stmt.Exec(id)
	if err != nil {
		log.Fatal(err)
	}

	affected, err = res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("DELETE: ", affected)
	db.Close()
}

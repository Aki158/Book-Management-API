package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Aki158/School-API/pkg/structs"
	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	UseDb *sql.DB
}

func (db *Database) Connect() {
	database, err := sql.Open("mysql", "testUser:testPassword@tcp(db:3306)/school?charset=utf8mb4")
	if err != nil {
		log.Println(err)
		return
	}

	db.UseDb = database
}

func GetSqlQuery(facilitatorId int, page int, limit int, sort string, order string, name_like string, loginId_like string) (string, []interface{}) {
	var args []interface{}
	skipRows := (page - 1) * limit
	sortTarget := "s.student_id"
	sqlQuery := `
		SELECT 
			s.student_id,
			s.student_name,
			s.login_id,
			c.class_id,
			c.class_name
		FROM 
			teachers t
		LEFT JOIN 
			classes c ON t.teacher_id = c.teacher_id
		LEFT JOIN
			students s ON c.class_id = s.class_id
		WHERE
			t.teacher_id = ? 
	`
	args = append(args, facilitatorId)

	// {key}_like により部分一致検索できるように設定する
	if name_like != "" {
		sqlQuery += "AND s.student_name LIKE ? "
		args = append(args, "%"+name_like+"%")
	}
	
	if loginId_like != "" {
		sqlQuery += "AND s.login_id LIKE ? "
		args = append(args, "%"+loginId_like+"%")
	}

	// sort によりソートのキーを設定する
	switch sort {
	case "name":
		sortTarget = "s.student_name"
	case "loginid":
		sortTarget = "s.login_id"
	}

	sqlQuery += fmt.Sprintf("ORDER BY %s %s LIMIT %d OFFSET %d;", sortTarget, order, limit, skipRows)

	return sqlQuery, args
}

func (db *Database) Read(facilitatorId int, page int, limit int, sort string, order string, name_like string, loginId_like string) ([]structs.Student) {
	var responseStudents []structs.Student

	// クエリ実行に必要な値を取得する
	query, args := GetSqlQuery(facilitatorId, page, limit, sort, order, name_like, loginId_like)
	
	res, err := db.UseDb.Query(query, args...)
	if err != nil {
		log.Println(err)
	}

	// resの取得結果をresponseStudentsに入れる
	for res.Next() {
		var student structs.Student

		err = res.Scan(&student.Id, &student.Name, &student.LoginId, &student.ClassRoom.Id, &student.ClassRoom.Name)
		if err != nil {
			log.Println(err)
		}

		responseStudents = append(responseStudents, student)
	}
	return responseStudents
}

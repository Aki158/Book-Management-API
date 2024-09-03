package db_test

import (
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/Aki158/School-API/pkg/db"
	"github.com/Aki158/School-API/pkg/structs"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func Test_Connect(t *testing.T) {
	t.Run(
		"Successful connection",
		func(t *testing.T) {
			// sqlmockを使ってモックデータベースと期待する振る舞いを設定する
			mockDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
			}
			defer mockDB.Close()

			// Database.Connectメソッドをテストする
			mydb := &db.Database{UseDb: mockDB}
			mydb.Connect()

			// モックの期待した振る舞いがすべて満たされたかを検証する
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		},
	)
}

func Test_Read(t *testing.T) {
	// sqlmockを使用して、データベースのモックを作成する
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	// モックDBラッパーの作成する
	mydb := &db.Database{UseDb: mockDB}

	// クエリと期待される引数を設定する
	expectedQuery := `
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
		ORDER BY s.student_id asc LIMIT 5 OFFSET 0;
	`

	// モックDBのクエリの期待値を設定する
	rows := sqlmock.NewRows([]string{"student_id", "student_name", "login_id", "class_id", "class_name"}).
		AddRow(1, "佐藤", "foo123", 1, "クラスA").
		AddRow(3, "田中", "baz789", 3, "クラスC").
		AddRow(4, "加藤", "hoge0000", 1, "クラスA").
		AddRow(6, "佐々木", "piyo5678", 3, "クラスC").
		AddRow(7, "小田原", "fizz9999", 1, "クラスA")

	mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(1).WillReturnRows(rows)

	// 実際の関数呼び出し
	students := mydb.Read(1, 1, 5, "id", "asc", "", "")

	// 期待される結果を確認する
	expectedStudents := []structs.Student{
		{Id: 1, Name: "佐藤", LoginId: "foo123", ClassRoom: structs.Class{Id: 1, Name: "クラスA"}},
		{Id: 3, Name: "田中", LoginId: "baz789", ClassRoom: structs.Class{Id: 3, Name: "クラスC"}},
		{Id: 4, Name: "加藤", LoginId: "hoge0000", ClassRoom: structs.Class{Id: 1, Name: "クラスA"}},
		{Id: 6, Name: "佐々木", LoginId: "piyo5678", ClassRoom: structs.Class{Id: 3, Name: "クラスC"}},
		{Id: 7, Name: "小田原", LoginId: "fizz9999", ClassRoom: structs.Class{Id: 1, Name: "クラスA"}},
	}
	
	// 期待される結果と実際の結果が一致するか確認する
	assert.Equal(t, expectedStudents, students)

	// モックが期待したクエリをすべて実行したかを確認する
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func Test_GetSqlQuery(t *testing.T) {
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

	tests := []struct {
		name           string
		facilitatorId  int
		page           int
		limit          int
		sort           string
		order          string
		name_like      string
		loginId_like   string
		expectedQuery  string
		expectedArgs   []interface{}
	}{
		{
			name:          "Basic query",
			facilitatorId: 1,
			page:          1,
			limit:         5,
			sort:          "id",
			order:         "asc",
			name_like:     "",
			loginId_like:  "",
			expectedQuery: sqlQuery + "ORDER BY s.student_id asc LIMIT 5 OFFSET 0;",
			expectedArgs: []interface{}{1},
		},
		{
			name:          "Check Pagination",
			facilitatorId: 1,
			page:          2,
			limit:         5,
			sort:          "id",
			order:         "asc",
			name_like:     "",
			loginId_like:  "",
			expectedQuery: sqlQuery + "ORDER BY s.student_id asc LIMIT 5 OFFSET 5;",
			expectedArgs: []interface{}{1},
		},
		{
			name:          "Query with name_like filter",
			facilitatorId: 1,
			page:          1,
			limit:         5,
			sort:          "id",
			order:         "asc",
			name_like:     "Johnson",
			loginId_like:  "",
			expectedQuery: sqlQuery + "AND s.student_name LIKE ? ORDER BY s.student_id asc LIMIT 5 OFFSET 0;",
			expectedArgs: []interface{}{1, "%"+"Johnson"+"%"},
		},
		{
			name:          "Query with loginId_like filter",
			facilitatorId: 1,
			page:          1,
			limit:         5,
			sort:          "id",
			order:         "asc",
			name_like:     "",
			loginId_like:  "foo123",
			expectedQuery: sqlQuery + "AND s.login_id LIKE ? ORDER BY s.student_id asc LIMIT 5 OFFSET 0;",
			expectedArgs: []interface{}{1, "%"+"foo123"+"%"},
		},
		{
			name:          "Check sort name",
			facilitatorId: 1,
			page:          1,
			limit:         5,
			sort:          "name",
			order:         "asc",
			name_like:     "",
			loginId_like:  "",
			expectedQuery: sqlQuery + "ORDER BY s.student_name asc LIMIT 5 OFFSET 0;",
			expectedArgs: []interface{}{1},
		},
		{
			name:          "Check sort loginId",
			facilitatorId: 1,
			page:          1,
			limit:         5,
			sort:          "loginid",
			order:         "asc",
			name_like:     "",
			loginId_like:  "",
			expectedQuery: sqlQuery + "ORDER BY s.login_id asc LIMIT 5 OFFSET 0;",
			expectedArgs: []interface{}{1},
		},
		{
			name:          "Query with all filters",
			facilitatorId: 2,
			page:          4,
			limit:         15,
			sort:          "name",
			order:         "desc",
			name_like:     "鈴木",
			loginId_like:  "baz789",
			expectedQuery: `
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
				AND s.student_name LIKE ? 
				AND s.login_id LIKE ? 
				ORDER BY s.student_name desc LIMIT 15 OFFSET 45;`,
			expectedArgs: []interface{}{2, "%"+"鈴木"+"%", "%"+"baz789"+"%"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, args := db.GetSqlQuery(tt.facilitatorId, tt.page, tt.limit, tt.sort, tt.order, tt.name_like, tt.loginId_like)

			compactQuery := strings.Join(strings.Fields(query), " ")
			expectedCompactQuery := strings.Join(strings.Fields(tt.expectedQuery), " ")

			if compactQuery != expectedCompactQuery {
				t.Errorf("Expected query:\n%s\nGot:\n%s", expectedCompactQuery, query)
			}

			if !reflect.DeepEqual(args, tt.expectedArgs) {
				t.Errorf("Expected args:\n%v\nGot:\n%v", tt.expectedArgs, args)
			}
		})
	}
}

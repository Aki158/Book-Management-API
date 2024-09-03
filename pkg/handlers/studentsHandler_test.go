package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/Aki158/School-API/pkg/db"
	"github.com/Aki158/School-API/pkg/handlers"
	"github.com/DATA-DOG/go-sqlmock"
)

func Test_StudentsHandler(t *testing.T) {
    // sqlmockを使用して、データベースのモックを作成
    mockDB, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
    }
    defer mockDB.Close()

    // モックDBラッパーの作成
    mydb := &db.Database{UseDb: mockDB}

    handler := handlers.StudentsHandler(mydb)

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

    tests := []struct {
        name               string
        query              string
        expectedStatusCode int
        setupMock          func()
    }{
        {
            name:               "Valid request with results",
            query:              "facilitator_id=1&page=1",
            expectedStatusCode: http.StatusOK,
            setupMock: func() {
                // クエリと期待される結果を設定する
                rows := sqlmock.NewRows([]string{"student_id", "student_name", "login_id", "class_id", "class_name"}).
                    AddRow(1, "佐藤", "foo123", 1, "クラスA").
                    AddRow(3, "田中", "baz789", 3, "クラスC").
                    AddRow(4, "加藤", "hoge0000", 1, "クラスA").
                    AddRow(6, "佐々木", "piyo5678", 3, "クラスC").
                    AddRow(7, "小田原", "fizz9999", 1, "クラスA")
                mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(1).WillReturnRows(rows)
            },
        },
        {
            name:               "Invalid facilitator_id",
            query:              "facilitator_id=abc",
            expectedStatusCode: http.StatusBadRequest,
            setupMock:          func() {},
        },
		{
            name:               "Invalid page parameter",
            query:              "facilitator_id=1&page=0",
            expectedStatusCode: http.StatusBadRequest,
            setupMock:          func() {},
        },
		{
            name:               "Invalid limit parameter",
            query:              "facilitator_id=1&limit=0",
            expectedStatusCode: http.StatusBadRequest,
            setupMock:          func() {},
        },
        {
            name:               "Invalid sort parameter",
            query:              "facilitator_id=1&sort=unknown",
            expectedStatusCode: http.StatusBadRequest,
            setupMock:          func() {},
        },
        {
            name:               "Invalid order parameter",
            query:              "facilitator_id=1&order=invalid",
            expectedStatusCode: http.StatusBadRequest,
            setupMock:          func() {},
        },
        {
            name:               "Facilitator with no students",
            query:              "facilitator_id=100",
            expectedStatusCode: http.StatusNotFound,
            setupMock: func() {
                // 空の結果を返す
                rows := sqlmock.NewRows([]string{"student_id", "student_name", "login_id", "class_id", "class_name"})
                mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(100).WillReturnRows(rows)
            },
        },
	}

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // テストケースごとにモックを設定する
            tt.setupMock()

            req := httptest.NewRequest(http.MethodGet, "/students?"+tt.query, nil)
            rr := httptest.NewRecorder()

            handler.ServeHTTP(rr, req)

            if status := rr.Code; status != tt.expectedStatusCode {
                t.Errorf("Handler returned wrong status code: got %v want %v",
                    status, tt.expectedStatusCode)
            }

            if tt.expectedStatusCode == http.StatusOK {
                if !strings.Contains(rr.Header().Get("Content-Type"), "application/json") {
                    t.Errorf("Handler did not set Content-Type application/json")
                }
            }
        })
    }

    // モックが期待したクエリをすべて実行したかを確認する
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("There were unfulfilled expectations: %s", err)
    }
}

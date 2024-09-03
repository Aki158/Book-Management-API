```mermaid
erDiagram
  teachers ||--o{ students : "has"
  teachers ||--o{ classes: "has"
  classes ||--o{ students: "has"

  teachers {
    int id "PK"
    timestamp created_at
    timestamp updated_at
  }

  classes {
    int id "PK"
    int teacher_id "FK"
    string class_name "UNIQUE"
    timestamp created_at
    timestamp updated_at
  }

  students {
    int id "PK"
    int teacher_id "FK"
    int class_id "FK"
    string student_name
    string login_id "UNIQUE"
    timestamp created_at
    timestamp updated_at
  }
```
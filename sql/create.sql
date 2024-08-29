CREATE SCHEMA school;
USE school;

CREATE TABLE teachers (
    teacher_id INT PRIMARY KEY
);

CREATE TABLE classes (
    class_id INT PRIMARY KEY,
    teacher_id INT,
    class_name VARCHAR(255) UNIQUE NOT NULL,
    FOREIGN KEY (teacher_id) REFERENCES teachers(teacher_id)
);

CREATE TABLE students (
    student_id INT PRIMARY KEY,
    teacher_id INT,
    class_id INT,
    student_name VARCHAR(255) NOT NULL,
    login_id VARCHAR(255) UNIQUE NOT NULL,
    FOREIGN KEY (teacher_id) REFERENCES teachers(teacher_id),
    FOREIGN KEY (class_id) REFERENCES classes(class_id)
);

INSERT INTO teachers (teacher_id) VALUES (1);
INSERT INTO teachers (teacher_id) VALUES (2);

INSERT INTO classes(class_id,teacher_id,class_name) VALUES (1,1,"クラスA");
INSERT INTO classes(class_id,teacher_id,class_name) VALUES (2,2,"クラスB");
INSERT INTO classes(class_id,teacher_id,class_name) VALUES (3,1,"クラスC");

INSERT INTO students(student_id,teacher_id,class_id,student_name,login_id) VALUES (1,1,1,"佐藤","foo123");
INSERT INTO students(student_id,teacher_id,class_id,student_name,login_id) VALUES (2,2,2,"鈴木","bar456");
INSERT INTO students(student_id,teacher_id,class_id,student_name,login_id) VALUES (3,1,3,"田中","baz789");
INSERT INTO students(student_id,teacher_id,class_id,student_name,login_id) VALUES (4,1,1,"加藤","hoge0000");
INSERT INTO students(student_id,teacher_id,class_id,student_name,login_id) VALUES (5,2,2,"太田","fuga1234");
INSERT INTO students(student_id,teacher_id,class_id,student_name,login_id) VALUES (6,1,3,"佐々木","piyo5678");
INSERT INTO students(student_id,teacher_id,class_id,student_name,login_id) VALUES (7,1,1,"小田原","fizz9999");
INSERT INTO students(student_id,teacher_id,class_id,student_name,login_id) VALUES (8,2,2,"Smith","buzz777");
INSERT INTO students(student_id,teacher_id,class_id,student_name,login_id) VALUES (9,1,3,"Johnson","fizzbuzz#123");
INSERT INTO students(student_id,teacher_id,class_id,student_name,login_id) VALUES (10,1,1,"Williams","xxx:42");

CREATE TABLE users (
    ID SERIAL PRIMARY KEY,
    Name VARCHAR(50),
    Email VARCHAR(100) UNIQUE NOT NULL,
    Password VARCHAR(100) NOT NULL,
    Role VARCHAR(50) NOT NULL  -- Roles like 'Teacher', 'Student', 'Parent'
);

CREATE TABLE subjects (
    ID SERIAL PRIMARY KEY,
    Name VARCHAR(100) NOT NULL
);

CREATE TABLE topics (
    ID SERIAL PRIMARY KEY,
    SubjectID INTEGER NOT NULL,
    Name VARCHAR(100) NOT NULL,
    FOREIGN KEY (SubjectID) REFERENCES subjects(ID)
);

CREATE TABLE assignments (
    ID SERIAL PRIMARY KEY,
    Name VARCHAR(100) NOT NULL,
    TopicID INTEGER NOT NULL,
    Weight INTEGER NOT NULL,
    TeacherID INTEGER NOT NULL,    
    FOREIGN KEY (TopicID) REFERENCES topics(ID)
);

CREATE TABLE questions (
    ID SERIAL PRIMARY KEY,
    AssignmentID INTEGER NOT NULL,
    Text TEXT NOT NULL,
    FOREIGN KEY (AssignmentID) REFERENCES assignments(ID)
);

CREATE TABLE answers (
    ID SERIAL PRIMARY KEY,
    QuestionID INTEGER NOT NULL,
    Text TEXT NOT NULL,
    IsCorrect BOOLEAN NOT NULL,
    FOREIGN KEY (QuestionID) REFERENCES questions(ID)
);

CREATE TABLE user_activity_logs (
    ID SERIAL PRIMARY KEY,
    UserID INTEGER NOT NULL,
    Timestamp TIMESTAMP NOT NULL DEFAULT,
    Route TEXT NOT NULL,
    FOREIGN KEY (UserID) REFERENCES users(ID)
);

CREATE TABLE student_assignments (
    ID SERIAL PRIMARY KEY,
    AssignmentID INTEGER NOT NULL,
    StudentID INTEGER NOT NULL,
    Score INTEGER,  -- Nullable initially, becomes non-null when assignment is completed
    FOREIGN KEY (AssignmentID) REFERENCES assignments(ID),
    FOREIGN KEY (StudentID) REFERENCES users(ID)
);

CREATE TABLE student_performance_by_subject (
    SubjectID INTEGER NOT NULL,
    StudentID INTEGER NOT NULL,
    OverallScore INTEGER,  -- Computed and stored for performance reasons
    FOREIGN KEY (SubjectID) REFERENCES subjects(ID),
    FOREIGN KEY (StudentID) REFERENCES users(ID),
    PRIMARY KEY (SubjectID, StudentID)
);


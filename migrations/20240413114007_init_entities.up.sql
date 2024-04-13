CREATE TABLE IF NOT EXISTS users(
  ID SERIAL PRIMARY KEY,
  Name VARCHAR(50),
  Email VARCHAR(100) NOT NULL UNIQUE,
  Password VARCHAR(100) NOT NULL,
  Role VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS reports(
  ID SERIAL PRIMARY KEY,
  UserID INTEGER NOT NULL,
  FOREIGN KEY (UserID) REFERENCES users(ID)
);

CREATE TABLE IF NOT EXISTS userActivity(
  ID SERIAL PRIMARY KEY,
  UserID INTEGER NOT NULL,
  Route TEXT NOT NULL,
  Time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (UserID) REFERENCES users(ID)
);

CREATE TABLE IF NOT EXISTS topics(
  ID SERIAL PRIMARY KEY,
  Name VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS subjects(
  ID SERIAL PRIMARY KEY,
  Name VARCHAR(100) NOT NULL,
  OverallScore INTEGER NOT NULL, 
  TeacherID INTEGER NOT NULL,
  StudentID INTEGER NOT NULL,
  FOREIGN KEY (TeacherID) REFERENCES users(ID),
  FOREIGN KEY (StudentID) REFERENCES users(ID)
);

CREATE TABLE IF NOT EXISTS assignments(
  ID SERIAL PRIMARY KEY,
  Name VARCHAR(100) NOT NULL,
  Score INTEGER NOT NULL,
  TopicID INTEGER NOT NULL,  
  SubjectID INTEGER NOT NULL,
  FOREIGN KEY (TopicID) REFERENCES topics(ID),
  FOREIGN KEY (SubjectID) REFERENCES subjects(ID)
);

CREATE TABLE IF NOT EXISTS questions(
  ID SERIAL PRIMARY KEY,
  Name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS answers(
  ID SERIAL PRIMARY KEY,
  Name TEXT NOT NULL,
  isCorrect BOOLEAN NOT NULL,
  QuestionID INTEGER NOT NULL,
  FOREIGN KEY (QuestionID) REFERENCES questions(ID)
);

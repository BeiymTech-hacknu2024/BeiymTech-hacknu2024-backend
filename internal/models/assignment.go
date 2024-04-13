package models

type Topic struct {
	ID   int
	Name string
}

type Subject struct {
	ID           int
	Name         string
	OverallScore int
	TeacherID    int
	StudentID    int
}

type Assignment struct {
	ID        int
	Name      string
	Score     int
	TopicID   int
	SubjectID int
}

type Question struct {
	ID   int
	Name string
}

type Answer struct {
	ID         int
	Name       string
	IsCorrect  bool
	QuestionID int
}

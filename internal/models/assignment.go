package models

type Subject struct {
	ID   int
	Name string
}

type Topic struct {
	ID        int
	Name      string
	SubjectID int
}

type Assignment struct {
	ID        int
	Name      string
	TopicID   int
	Weight    int
	TeacherID int
}

type Question struct {
	ID           int
	Text         string
	AssignmentID int
}

type Answer struct {
	ID         int
	Text       string
	IsCorrect  bool
	QuestionID int
}


type CreateAssignmentRequest struct {
	Assignment Assignment
	StudentIDs []int
}

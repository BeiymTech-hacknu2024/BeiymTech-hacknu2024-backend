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
	ID        int    `db:"id"`
	Name      string `db:"name"`
	TopicID   int    `db:"topicid"`
	Weight    int    `db:"weight"`
	TeacherID int    `db:"teacherid"`
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

type AnswerSubmission struct {
	QuestionID       int `json:"questionId"`
	SelectedAnswerID int `json:"selectedAnswerId"`
}

type SubmitRequest struct {
	StudentID    int                `json:"studentId"`
	AssignmentID int                `json:"assignmentId"`
	Answers      []AnswerSubmission `json:"answers"`
}

type Performance struct {
	ID              int
	Kinematics      int
	Dynamics        int
	Electrodynamics int
	Acids           int
	ChemicalBonding int
	Trigonometry    int
	LinearAlgebra   int
	Geometry        int
	Probability     int
}

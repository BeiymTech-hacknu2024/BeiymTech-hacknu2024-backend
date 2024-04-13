-- Drop foreign key constraints before dropping tables

-- Drop the foreign key constraint from the "student_performance_by_subject" table
ALTER TABLE student_performance_by_subject DROP CONSTRAINT IF EXISTS student_performance_by_subject_assignmentid_fkey;

-- Drop the foreign key constraint from the "student_assignments" table
ALTER TABLE student_assignments DROP CONSTRAINT IF EXISTS student_assignments_assignmentid_fkey;

-- Drop the foreign key constraint from the "student_assignments" table
ALTER TABLE student_assignments DROP CONSTRAINT IF EXISTS student_assignments_studentid_fkey;

-- Drop the foreign key constraint from the "answers" table
ALTER TABLE answers DROP CONSTRAINT IF EXISTS answers_questionid_fkey;

-- Drop the foreign key constraint from the "questions" table
ALTER TABLE questions DROP CONSTRAINT IF EXISTS questions_assignmentid_fkey;

-- Drop the foreign key constraint from the "assignments" table
ALTER TABLE assignments DROP CONSTRAINT IF EXISTS assignments_topicid_fkey;

-- Drop the foreign key constraint from the "student_performance_by_subject" table
ALTER TABLE student_performance_by_subject DROP CONSTRAINT IF EXISTS student_performance_by_subject_subjectid_fkey;

-- Drop the foreign key constraint from the "user_activity_logs" table
ALTER TABLE user_activity_logs DROP CONSTRAINT IF EXISTS user_activity_logs_userid_fkey;

-- Now drop the tables

-- Drop the "answers" table
DROP TABLE IF EXISTS answers;

-- Drop the "questions" table
DROP TABLE IF EXISTS questions;

-- Drop the "student_assignments" table
DROP TABLE IF EXISTS student_assignments;

-- Drop the "student_performance_by_subject" table
DROP TABLE IF EXISTS student_performance_by_subject;

-- Drop the "user_activity_logs" table
DROP TABLE IF EXISTS user_activity_logs;

-- Drop the "assignments" table
DROP TABLE IF EXISTS assignments;

-- Drop the "topics" table
DROP TABLE IF EXISTS topics;

-- Drop the "subjects" table
DROP TABLE IF EXISTS subjects;

-- Drop the "users" table
DROP TABLE IF EXISTS users;

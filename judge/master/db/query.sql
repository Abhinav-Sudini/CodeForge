-- name: CreateQuestion :exec
INSERT INTO questions (
    question_id,
    question_category,
    question_name,
    question_description,
    input_description,
    output_description,
    constraints_description,
    time_constraint,
    mem_constraint,
    example_inputs,
    example_outputs
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
);

-- name: GetTimeAndMemConstraints :one
SELECT 
  time_constraint,mem_constraint
FROM questions
WHERE question_id = $1
LIMIT 1;

-- name: QuestionExists :one
SELECT EXISTS (
    SELECT 1
    FROM questions
    WHERE question_id = $1
);

-- name: GetAllQuestionsMinimalDetails :many
SELECT 
  question_id,
  question_name,
  question_category
FROM questions ;


-- name: GetQuestion :one
SELECT *
FROM questions 
WHERE question_id = $1;


-- name: CreateUserAndReturnId :one
INSERT INTO users (
  username,
  email,
  encrypted_password,
  role,
  created_at,
  last_online
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING id;

-- name: GetUser :one
SELECT *
FROM users 
WHERE email = $1;

-- name: UserNameExists :one
SELECT EXISTS (
    SELECT 1
    FROM users
    WHERE email = $1
);

-- name: CreateSubmissionAndReturnId :one
INSERT INTO submissions (
  user_id,
  question_id,
  submission_time,
  submited_code,
  code_runtime,
  verdict
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING submission_id;

-- name: UpdateVerdictForSubmition :exec
UPDATE submissions 
SET verdict = $2
WHERE submission_id = $1;




-- name: CreateVerdictStatsRecord :exec
INSERT INTO verdict_stats (
  submission_id,
  mem_usage,
  time_ms,
  not_accepted_test_case,
  not_accepted_test_case_stdout,
  stderr
) VALUES (
    $1, $2, $3, $4, $5, $6
);


-- name: GetSubmissionVerdict :one
SELECT 
s.submission_id,
s.submited_code,
s.question_id,
s.verdict,
v.mem_usage,
v.time_ms, v.not_accepted_test_case, v.not_accepted_test_case_stdout, v.stderr, s.submission_time
FROM submissions s 
LEFT JOIN verdict_stats v ON s.submission_id = v.submission_id
WHERE s.submission_id = $1 ;

-- name: GetVerdictStats :one
SELECT * FROM verdict_stats
WHERE submission_id = $1
LIMIT 1;

-- name: GetSubmissionVerdictAndQuestionid :one
SELECT verdict,question_id FROM submissions
WHERE submission_id = $1
LIMIT 1;


-- name: GetAllSubmissionOfQuestion :many
SELECT 
s.submission_id,
s.submited_code,
s.question_id,
s.verdict,
v.mem_usage,
v.time_ms, v.not_accepted_test_case, v.not_accepted_test_case_stdout, v.stderr, s.submission_time
FROM submissions s 
LEFT JOIN verdict_stats v ON s.submission_id = v.submission_id
WHERE s.question_id = $1 AND s.user_id = $2;


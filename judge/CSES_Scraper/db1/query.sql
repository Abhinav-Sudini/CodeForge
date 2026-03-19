-- name: CreateQuestion :exec
INSERT INTO questions (
    question_id,
    question_catagory,
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



-- name: UpdateVerdictForSubmition :exec
UPDATE submissions 
SET verdict = $2
WHERE submission_id = $1;




-- name: GetTimeAndMemConstraints :one
SELECT 
  time_constraint,mem_constraint
FROM questions
WHERE question_id = $1
LIMIT 1;


-- name: GetVerdictStats :one
SELECT * FROM verdict_stats
WHERE submission_id = $1
LIMIT 1;

-- name: GetSubmissionVerdictAndQuestionid :one
SELECT verdict,question_id FROM submissions
WHERE submission_id = $1
LIMIT 1;

-- name: QuestionExists :one
SELECT EXISTS (
    SELECT 1
    FROM questions
    WHERE question_id = $1
);

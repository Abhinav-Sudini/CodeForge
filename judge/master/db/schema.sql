
CREATE TABLE users (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    username VARCHAR,
    email VARCHAR,
    encrypted_password VARCHAR,
    role VARCHAR,
    created_at TIMESTAMP,
    last_online TIMESTAMP
);

CREATE TABLE questions (
    question_id INTEGER PRIMARY KEY,
    question_category TEXT,
    question_name TEXT,
    question_description TEXT,
    input_description TEXT,
    output_description TEXT,
    constraints_description TEXT,
    time_constraint INTEGER NOT NULL,
    mem_constraint INTEGER NOT NULL,
    example_inputs TEXT[],
    example_outputs TEXT[]
);


CREATE TABLE submissions (
    submission_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id INTEGER NOT NULL,
    question_id INTEGER NOT NULL,
    submission_time TIMESTAMP,
    submited_code TEXT,
    code_runtime TEXT,
    verdict TEXT,

    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
        REFERENCES users(id),

    CONSTRAINT fk_question
        FOREIGN KEY (question_id)
        REFERENCES questions(question_id)
);

CREATE TABLE verdict_stats (
    submission_id INTEGER PRIMARY KEY,
    mem_usage INT NOT NULL,
    time_ms INT NOT NULL,
    not_accepted_test_case INT,
    not_accepted_test_case_stdout TEXT,
    stderr TEXT,

    CONSTRAINT fk_submission
        FOREIGN KEY (submission_id)
        REFERENCES submissions(submission_id)
);

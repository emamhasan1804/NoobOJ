-- ============================================================= --
-- Schema for Competitive Programming Platform                   --
-- Generated for Go project                                      --
-- ============================================================= --

-- Users
CREATE TABLE IF NOT EXISTS users (
    username    VARCHAR(50)  PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    email       VARCHAR(100) NOT NULL UNIQUE,
    password    MEDIUMTEXT NOT NULL,
    user_type   ENUM('user','admin')  NOT NULL DEFAULT 'user',
    created_on  DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Profiles (avatar / picture path per user)
CREATE TABLE IF NOT EXISTS profiles (
    username    VARCHAR(50),
    path        VARCHAR(255) NOT NULL,
    PRIMARY KEY(username),
    FOREIGN KEY(username)
    REFERENCES users(username) 
    ON DELETE CASCADE 
    ON UPDATE CASCADE
);

-- Problems
DROP TABLE problems;
CREATE TABLE IF NOT EXISTS problems (
    id          	SERIAL       PRIMARY KEY,
    title       	VARCHAR(200) NOT NULL,
    statement   	TEXT         NOT NULL,
    input       	TEXT,
    output      	TEXT,
    constraints 	TEXT,
    author 			VARCHAR(50)	NOT NULL,
    time_limit 	 	INT			DEFAULT 1,
    memory_limit 	INT			DEFAULT 512,
    editorial 		TEXT,
    code 			TEXT,
    visibility 		BOOLEAN		DEFAULT FALSE,
    CONSTRAINT problems_users
    FOREIGN KEY(author) 
    REFERENCES users(username)
    ON DELETE CASCADE 
    ON UPDATE CASCADE
);

-- Test Cases
CREATE TABLE IF NOT EXISTS test_cases (
    id          SERIAL  PRIMARY KEY,
    problem_id  BIGINT UNSIGNED NOT NULL,
    input       TEXT    NOT NULL,
    output      TEXT    NOT NULL,
    type		ENUM('sample','hidden') DEFAULT 'hidden',
    CONSTRAINT test_cases_problems
    FOREIGN KEY(problem_id)
    REFERENCES problems(id)
    ON DELETE CASCADE 
    ON UPDATE CASCADE
);

-- Tags
CREATE TABLE IF NOT EXISts tags(
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL
);
INSERT IGNORE INTO tags (name) VALUES 
('2-sat'),
('binary search'),
('bitmasks'),
('brute force'),
('chinese remainder theorem'),
('combinatorics'),
('constructive algorithms'),
('data structures'),
('dfs and similar'),
('divide and conquer'),
('dp'),
('dsu'),
('expression parsing'),
('fft'),
('flows'),
('games'),
('geometry'),
('graph matchings'),
('graphs'),
('greedy'),
('hashing'),
('implementation'),
('interactive'),
('math'),
('matrices'),
('meet-in-the-middle'),
('number theory'),
('probabilities'),
('schedules'),
('shortest paths'),
('sortings'),
('string suffix structures'),
('strings'),
('ternary search'),
('trees'),
('two pointers');

-- Tags of a Problem
CREATE TABLE IF NOT EXISTS problem_tags (
    id          SERIAL      PRIMARY KEY,
    problem_id  BIGINT UNSIGNED NOT NULL,        
    tag_id BIGINT UNSIGNED NOT NULL,
    CONSTRAINT tags_problems
    FOREIGN KEY(problem_id)
    REFERENCES problems(id) 
    ON DELETE CASCADE 
    ON UPDATE CASCADE,
    CONSTRAINT problem_tags_tags
    FOREIGN KEY(tag_id)
    REFERENCES tags(id)
    ON DELETE CASCADE 
    ON UPDATE CASCADE
);

-- Ratings
CREATE TABLE IF NOT EXISTS ratings (
    id          SERIAL      PRIMARY KEY,
    problem_id  BIGINT UNSIGNED NOT NULL,
    rating      INT DEFAULT 800,
    CONSTRAINT ratings_problems
    FOREIGN KEY(problem_id)
    REFERENCES problems(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE
);

-- Contests
CREATE TABLE IF NOT EXISTS contests (
    id          SERIAL       PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    start_time  TIMESTAMP    NOT NULL,
    end_time    TIMESTAMP    NOT NULL
);

-- Tasks (problems assigned to a contest with a score)
CREATE TABLE IF NOT EXISTS tasks (
    id          SERIAL PRIMARY KEY,
    problem_id  INT    NOT NULL REFERENCES problems(id) ON DELETE CASCADE ON UPDATE CASCADE,
    contest_id  INT    NOT NULL REFERENCES contests(id) ON DELETE CASCADE ON UPDATE CASCADE,
    score       INT    NOT NULL DEFAULT 0,
    UNIQUE (problem_id, contest_id)
);

-- Participants (users enrolled in a contest)
CREATE TABLE IF NOT EXISTS participants (
    contest_id  INT         NOT NULL REFERENCES contests(id) ON DELETE CASCADE ON UPDATE CASCADE,
    username    VARCHAR(50) NOT NULL REFERENCES users(username) ON DELETE CASCADE ON UPDATE CASCADE,
    PRIMARY KEY (contest_id, username)
);

-- Authors (users who authored / set a contest)
CREATE TABLE IF NOT EXISTS authors (
    contest_id  INT         NOT NULL REFERENCES contests(id) ON DELETE CASCADE ON UPDATE CASCADE,
    username    VARCHAR(50) NOT NULL REFERENCES users(username) ON DELETE CASCADE ON UPDATE CASCADE,
    PRIMARY KEY (contest_id, username)
);

-- Submissions
CREATE TABLE IF NOT EXISTS submissions (
    id           SERIAL      PRIMARY KEY,
    problem_id   INT         NOT NULL REFERENCES problems(id) ON DELETE CASCADE ON UPDATE CASCADE,
    username     VARCHAR(50) NOT NULL REFERENCES users(username) ON DELETE CASCADE ON UPDATE CASCADE,
    code         TEXT        NOT NULL,
    submitted_at TIMESTAMP   NOT NULL DEFAULT NOW(),
    verdict      VARCHAR(20) NOT NULL DEFAULT 'pending'
);

-- Streaks (daily problem-solving streaks per user)
CREATE TABLE IF NOT EXISTS streaks (
    username    VARCHAR(50) NOT NULL REFERENCES users(username) ON DELETE CASCADE ON UPDATE CASCADE,
    problem_id  INT         NOT NULL REFERENCES problems(id) ON DELETE CASCADE ON UPDATE CASCADE,
    date        DATE        NOT NULL,
    PRIMARY KEY (username, problem_id)
);

-- =============================================================
-- Indexes for common query patterns
-- =============================================================

CREATE INDEX IF NOT EXISTS idx_submissions_username    ON submissions(username);
CREATE INDEX IF NOT EXISTS idx_submissions_problem_id  ON submissions(problem_id);
CREATE INDEX IF NOT EXISTS idx_submissions_verdict     ON submissions(verdict);
CREATE INDEX IF NOT EXISTS idx_tasks_contest_id        ON tasks(contest_id);
CREATE INDEX IF NOT EXISTS idx_test_cases_problem_id   ON test_cases(problem_id);
CREATE INDEX IF NOT EXISTS idx_tags_problem_id         ON tags(problem_id);
CREATE INDEX IF NOT EXISTS idx_ratings_problem_id      ON ratings(problem_id);
CREATE INDEX IF NOT EXISTS idx_streaks_username        ON streaks(username);
CREATE INDEX IF NOT EXISTS idx_streaks_date            ON streaks(date);
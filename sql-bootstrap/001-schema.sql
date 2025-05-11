-- Day schedules enum
CREATE TYPE day_of_week AS ENUM ('monday', 'tuesday', 'wednesday', 'thursday', 'friday');

-- Teachers table
CREATE TABLE teachers (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Subjects table
CREATE TABLE subjects (
    id SERIAL PRIMARY KEY,
    code VARCHAR(20) NOT NULL, -- e.g., "CHEM101"
    name VARCHAR(255) NOT NULL, -- e.g., "General Chemistry 1"
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Classrooms table
CREATE TABLE classrooms (
    id SERIAL PRIMARY KEY,
    building VARCHAR(100) NOT NULL,
    room_number VARCHAR(20) NOT NULL,
    capacity INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(building, room_number)
);

-- Students table
CREATE TABLE students (
    id SERIAL PRIMARY KEY,
    student_id VARCHAR(50) UNIQUE NOT NULL, -- university student ID
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Sections table (main join table)
CREATE TABLE sections (
    id SERIAL PRIMARY KEY,
    subject_id INTEGER NOT NULL REFERENCES subjects(id),
    teacher_id INTEGER NOT NULL REFERENCES teachers(id),
    classroom_id INTEGER NOT NULL REFERENCES classrooms(id),
    section_code VARCHAR(20) NOT NULL, -- e.g., "001", "002"
    start_time TIME NOT NULL,
    duration_minutes INTEGER NOT NULL DEFAULT 50,
    max_enrollment INTEGER NOT NULL,
    current_enrollment INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(subject_id, section_code),
    CHECK (start_time >= '07:30:00'),
    CHECK (start_time + (duration_minutes || ' minutes')::INTERVAL <= '22:00:00'),
    CHECK (duration_minutes IN (50, 80)),
    CHECK (current_enrollment <= max_enrollment),
    CHECK (current_enrollment >= 0)
);

-- Section days (many-to-many relationship for days)
CREATE TABLE section_days (
    section_id INTEGER NOT NULL REFERENCES sections(id) ON DELETE CASCADE,
    day day_of_week NOT NULL,
    PRIMARY KEY (section_id, day)
);

-- Student enrollments
CREATE TABLE enrollments (
    id SERIAL PRIMARY KEY,
    student_id INTEGER NOT NULL REFERENCES students(id),
    section_id INTEGER NOT NULL REFERENCES sections(id),
    enrollment_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(student_id, section_id)
);

-- Indexes for performance
CREATE INDEX idx_sections_subject_id ON sections(subject_id);
CREATE INDEX idx_sections_teacher_id ON sections(teacher_id);
CREATE INDEX idx_sections_classroom_id ON sections(classroom_id);
CREATE INDEX idx_enrollments_student_id ON enrollments(student_id);
CREATE INDEX idx_enrollments_section_id ON enrollments(section_id);
CREATE INDEX idx_section_days_section_id ON section_days(section_id);

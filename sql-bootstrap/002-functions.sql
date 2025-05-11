-- Function to check for schedule conflicts
CREATE OR REPLACE FUNCTION check_schedule_conflict(
    p_student_id INTEGER,
    p_section_id INTEGER
) RETURNS BOOLEAN AS $$
DECLARE
    v_conflict_count INTEGER;
BEGIN
    WITH new_section_info AS (
        SELECT
            s.start_time,
            s.start_time + (s.duration_minutes || ' minutes')::INTERVAL as end_time,
            array_agg(sd.day) as days
        FROM sections s
        JOIN section_days sd ON s.id = sd.section_id
        WHERE s.id = p_section_id
        GROUP BY s.id, s.start_time, s.duration_minutes
    ),
    enrolled_sections AS (
        SELECT
            s.id,
            s.start_time,
            s.start_time + (s.duration_minutes || ' minutes')::INTERVAL as end_time,
            array_agg(sd.day) as days
        FROM sections s
        JOIN section_days sd ON s.id = sd.section_id
        JOIN enrollments e ON s.id = e.section_id
        WHERE e.student_id = p_student_id
        GROUP BY s.id, s.start_time, s.duration_minutes
    )
    SELECT COUNT(*)
    INTO v_conflict_count
    FROM new_section_info nsi, enrolled_sections es
    WHERE
        nsi.days && es.days  -- arrays have common elements (days overlap)
        AND (
            (nsi.start_time >= es.start_time AND nsi.start_time < es.end_time)
            OR (nsi.end_time > es.start_time AND nsi.end_time <= es.end_time)
            OR (nsi.start_time <= es.start_time AND nsi.end_time >= es.end_time)
        );

    RETURN v_conflict_count > 0;
END;
$$ LANGUAGE plpgsql;

-- Function to prevent enrollment conflicts (returns trigger)
CREATE OR REPLACE FUNCTION prevent_enrollment_conflicts()
RETURNS TRIGGER AS $$
BEGIN
    IF check_schedule_conflict(NEW.student_id, NEW.section_id) THEN
        RAISE EXCEPTION 'Schedule conflict detected. Cannot enroll in this section.';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Function to update current enrollment count (returns trigger)
CREATE OR REPLACE FUNCTION update_enrollment_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE sections
        SET current_enrollment = current_enrollment + 1
        WHERE id = NEW.section_id;

        -- Check if we exceed max enrollment
        IF (SELECT current_enrollment FROM sections WHERE id = NEW.section_id) >
           (SELECT max_enrollment FROM sections WHERE id = NEW.section_id) THEN
            RAISE EXCEPTION 'Section is full. Cannot enroll.';
        END IF;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE sections
        SET current_enrollment = current_enrollment - 1
        WHERE id = OLD.section_id;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Function to update timestamp (returns trigger)
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

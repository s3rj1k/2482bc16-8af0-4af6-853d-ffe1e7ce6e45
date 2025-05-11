-- View for student schedules (for PDF generation)
CREATE VIEW student_schedule_view AS
SELECT
    e.student_id,
    s.id as section_id,
    sub.code as subject_code,
    sub.name as subject_name,
    sec.section_code,
    t.first_name as teacher_first_name,
    t.last_name as teacher_last_name,
    c.building,
    c.room_number,
    sec.start_time,
    sec.start_time + (sec.duration_minutes || ' minutes')::INTERVAL as end_time,
    sec.duration_minutes,
    array_agg(sd.day ORDER BY sd.day) as days
FROM enrollments e
JOIN sections sec ON e.section_id = sec.id
JOIN subjects sub ON sec.subject_id = sub.id
JOIN teachers t ON sec.teacher_id = t.id
JOIN classrooms c ON sec.classroom_id = c.id
JOIN section_days sd ON sec.id = sd.section_id
JOIN students s ON e.student_id = s.id
GROUP BY
    e.student_id, s.id, sec.id, sub.id, t.id, c.id;

-- Trigger to prevent enrollment conflicts
CREATE TRIGGER trg_prevent_enrollment_conflicts
BEFORE INSERT ON enrollments
FOR EACH ROW
EXECUTE FUNCTION prevent_enrollment_conflicts();

-- Trigger to update current enrollment count
CREATE TRIGGER trg_update_enrollment_count
AFTER INSERT OR DELETE ON enrollments
FOR EACH ROW
EXECUTE FUNCTION update_enrollment_count();

-- Update timestamp triggers
CREATE TRIGGER update_teachers_updated_at
BEFORE UPDATE ON teachers
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_subjects_updated_at
BEFORE UPDATE ON subjects
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_classrooms_updated_at
BEFORE UPDATE ON classrooms
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_students_updated_at
BEFORE UPDATE ON students
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_sections_updated_at
BEFORE UPDATE ON sections
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

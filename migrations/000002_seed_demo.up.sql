INSERT INTO users (email, password_hash, role, first_name, last_name, is_active)
VALUES
('admin@clinic.local', crypt('Admin123!', gen_salt('bf')), 'ADMIN', 'Main', 'Admin', TRUE),
('doctor1@clinic.local', crypt('Doctor123!', gen_salt('bf')), 'DOCTOR', 'John', 'Doctor', TRUE)
ON CONFLICT (email) DO NOTHING;

INSERT INTO doctor_profiles (user_id, specialization, phone)
SELECT u.id, 'Therapist', '+77000000000'
FROM users u
WHERE u.email = 'doctor1@clinic.local'
ON CONFLICT (user_id) DO NOTHING;

INSERT INTO patients (first_name, last_name, gender, phone, email, address, comment)
SELECT v.first_name, v.last_name, v.gender, v.phone, v.email, v.address, v.comment
FROM (
    VALUES
        ('Aruzhan', 'Saparova', 'F', '+77010000001', 'aruzhan@mail.local', 'Almaty', 'Demo patient'),
        ('Alikhan', 'Nurpeisov', 'M', '+77010000002', 'alikhan@mail.local', 'Almaty', 'Demo patient')
) AS v(first_name, last_name, gender, phone, email, address, comment)
WHERE NOT EXISTS (
    SELECT 1 FROM patients p WHERE p.email = v.email
);

INSERT INTO doctor_schedules (doctor_id, weekday, start_time, end_time, slot_minutes)
SELECT s.doctor_id, s.weekday, s.start_time, s.end_time, s.slot_minutes
FROM (
    SELECT u.id AS doctor_id, 0 AS weekday, TIME '09:00' AS start_time, TIME '17:00' AS end_time, 30 AS slot_minutes
    FROM users u WHERE u.email = 'doctor1@clinic.local'
    UNION ALL
    SELECT u.id, 1, TIME '09:00', TIME '17:00', 30 FROM users u WHERE u.email = 'doctor1@clinic.local'
    UNION ALL
    SELECT u.id, 2, TIME '09:00', TIME '17:00', 30 FROM users u WHERE u.email = 'doctor1@clinic.local'
    UNION ALL
    SELECT u.id, 3, TIME '09:00', TIME '17:00', 30 FROM users u WHERE u.email = 'doctor1@clinic.local'
    UNION ALL
    SELECT u.id, 4, TIME '09:00', TIME '17:00', 30 FROM users u WHERE u.email = 'doctor1@clinic.local'
) AS s
WHERE NOT EXISTS (
    SELECT 1
    FROM doctor_schedules ds
    WHERE ds.doctor_id = s.doctor_id
      AND ds.weekday = s.weekday
      AND ds.start_time = s.start_time
      AND ds.end_time = s.end_time
      AND ds.slot_minutes = s.slot_minutes
);

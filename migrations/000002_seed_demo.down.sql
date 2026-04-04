DELETE FROM doctor_schedules WHERE doctor_id IN (SELECT id FROM users WHERE email = 'doctor1@clinic.local');
DELETE FROM doctor_profiles WHERE user_id IN (SELECT id FROM users WHERE email = 'doctor1@clinic.local');
DELETE FROM appointments WHERE doctor_id IN (SELECT id FROM users WHERE email IN ('admin@clinic.local', 'doctor1@clinic.local'));
DELETE FROM patients WHERE email IN ('aruzhan@mail.local', 'alikhan@mail.local');
DELETE FROM users WHERE email IN ('admin@clinic.local', 'doctor1@clinic.local');

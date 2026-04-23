CREATE TABLE IF NOT EXISTS notifications (
    id             BIGSERIAL PRIMARY KEY,
    appointment_id BIGINT NOT NULL REFERENCES appointments(id) ON DELETE CASCADE,
    message        TEXT NOT NULL,
    sent_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_notifications_appointment ON notifications(appointment_id);

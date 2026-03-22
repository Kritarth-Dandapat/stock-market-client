-- Migration: 002_create_kyc_records.sql

CREATE TABLE IF NOT EXISTS kyc_records (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    encrypted_pan TEXT NOT NULL,
    pan_masked    TEXT NOT NULL,
    full_name     TEXT NOT NULL,
    date_of_birth DATE,
    status        TEXT NOT NULL DEFAULT 'pending',
    ckyc_number   TEXT,
    verified_at   TIMESTAMPTZ,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT kyc_records_status_check CHECK (
        status IN ('pending', 'in_progress', 'verified', 'rejected')
    )
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_kyc_records_user_id ON kyc_records (user_id);

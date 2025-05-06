CREATE TABLE IF NOT EXISTS public.anomalies (
  id          BIGSERIAL PRIMARY KEY,
  field       TEXT    NOT NULL,
  value       REAL    NOT NULL,
  timestamp   BIGINT  NOT NULL,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE IF NOT EXISTS public.turion_data_packets (
  packet_id        INTEGER    NOT NULL,
  packet_seq_ctrl  INTEGER    NOT NULL,
  packet_length    INTEGER    NOT NULL,
  ts               BIGINT     NOT NULL,
  subsystem_id     INTEGER    NOT NULL,
  temperature      REAL       NOT NULL,
  battery          REAL       NOT NULL,
  altitude         REAL       NOT NULL,
  signal           REAL       NOT NULL
);

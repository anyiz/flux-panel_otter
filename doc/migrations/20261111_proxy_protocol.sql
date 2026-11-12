-- Apply once before starting the upgraded panel. Existing rules remain off.
ALTER TABLE forward ADD COLUMN proxy_protocol_send INT NOT NULL DEFAULT 0;

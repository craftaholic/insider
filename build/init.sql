CREATE TABLE IF NOT EXISTS messages (
    id BIGSERIAL PRIMARY KEY,
    phone_number VARCHAR(20) NOT NULL,
    content TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'sent', 'failed')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    sent_at TIMESTAMP WITH TIME ZONE NULL,
    message_id VARCHAR(255) NULL,
    error_message TEXT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NULL
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_messages_status_created ON messages (status, created_at);
CREATE INDEX IF NOT EXISTS idx_messages_phone_number ON messages (phone_number);
CREATE INDEX IF NOT EXISTS idx_messages_message_id ON messages (message_id);
CREATE INDEX IF NOT EXISTS idx_messages_sent_at ON messages (sent_at);
CREATE INDEX IF NOT EXISTS idx_messages_updated_at ON messages (updated_at);
CREATE INDEX IF NOT EXISTS idx_messages_processing_stuck ON messages (status, updated_at) WHERE status = 'processing';

-- Insert sample data for testing
INSERT INTO messages (phone_number, content, status) VALUES 
    ('+905551111111', 'Test message 1 - Insider Project', 'pending'),
    ('+905551111112', 'Test message 2 - Hello World', 'pending'),
    ('+905551111113', 'Test message 3 - Sample Content', 'pending'),
    ('+905551111114', 'Test message 4 - Demo Message', 'pending'),
    ('+905551111115', 'Test message 5 - Testing System', 'pending'),
    ('+905551111116', 'Test message 6 - Additional Test', 'pending'),
    ('+905551111117', 'Test message 7 - Load Testing', 'pending'),
    ('+905551111118', 'Test message 8 - Batch Processing', 'pending');

-- Create a view for sent messages (optional, for easier querying)
CREATE VIEW sent_messages AS 
SELECT 
    id,
    phone_number,
    content,
    status,
    created_at,
    sent_at,
    message_id,
    EXTRACT(EPOCH FROM (sent_at - created_at)) as processing_time_seconds
FROM messages 
WHERE status = 'sent'
ORDER BY sent_at DESC;

-- Create a view for monitoring processing messages
CREATE VIEW processing_messages AS 
SELECT 
    id,
    phone_number,
    content,
    created_at,
    updated_at,
    EXTRACT(EPOCH FROM (NOW() - updated_at)) as seconds_processing
FROM messages 
WHERE status = 'processing'
ORDER BY updated_at ASC;

-- Create a view for stuck messages (processing > 10 minutes)
CREATE VIEW stuck_messages AS 
SELECT 
    id,
    phone_number,
    content,
    updated_at,
    EXTRACT(EPOCH FROM (NOW() - updated_at)) as stuck_duration_seconds
FROM messages 
WHERE status = 'processing' 
  AND updated_at < NOW() - INTERVAL '10 minutes'
ORDER BY updated_at ASC;

-- Function for getting_unsent_messages atomicly  
-- (For avoid the go application getting the same messages)
CREATE OR REPLACE FUNCTION get_unsent_messages(batch_size INTEGER DEFAULT 2)
RETURNS TABLE (
    id BIGINT,
    phone_number VARCHAR(20),
    content TEXT,
    created_at TIMESTAMP WITH TIME ZONE
) AS $$
BEGIN
    RETURN QUERY
    UPDATE messages 
    SET status = 'processing',
        updated_at = CURRENT_TIMESTAMP
    WHERE messages.id IN (
        SELECT m.id
        FROM messages m
        WHERE m.status = 'pending'
        ORDER BY m.created_at ASC
        LIMIT batch_size
        FOR UPDATE SKIP LOCKED
    )
    RETURNING messages.id, messages.phone_number, messages.content, messages.created_at;
END;
$$ LANGUAGE plpgsql;

-- Function to mark message as sent
CREATE OR REPLACE FUNCTION mark_message_sent(
    msg_id BIGINT,
    webhook_message_id VARCHAR(255)
) RETURNS VOID AS $$
BEGIN
    UPDATE messages 
    SET 
        status = 'sent',
        sent_at = CURRENT_TIMESTAMP,
        message_id = webhook_message_id,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = msg_id;
END;
$$ LANGUAGE plpgsql;

-- Function to mark message as failed
CREATE OR REPLACE FUNCTION mark_message_failed(
    msg_id BIGINT,
    error_msg TEXT
) RETURNS VOID AS $$
BEGIN
    UPDATE messages 
    SET 
        status = 'failed',
        error_message = error_msg,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = msg_id;
END;
$$ LANGUAGE plpgsql;

-- Function to reset stuck messages (processing > specified minutes)
CREATE OR REPLACE FUNCTION reset_stuck_messages(stuck_minutes INTEGER DEFAULT 10)
RETURNS INTEGER AS $$
DECLARE
    affected_count INTEGER;
BEGIN
    UPDATE messages 
    SET status = 'pending', 
        updated_at = NULL,
        error_message = CONCAT('Reset from stuck processing after ', stuck_minutes, ' minutes')
    WHERE status = 'processing' 
      AND updated_at < NOW() - (stuck_minutes || ' minutes')::INTERVAL;
    
    GET DIAGNOSTICS affected_count = ROW_COUNT;
    RETURN affected_count;
END;
$$ LANGUAGE plpgsql;

-- Function to get processing statistics
CREATE OR REPLACE FUNCTION get_processing_stats()
RETURNS TABLE (
    total_messages BIGINT,
    pending_count BIGINT,
    processing_count BIGINT,
    sent_count BIGINT,
    failed_count BIGINT,
    stuck_count BIGINT,
    avg_processing_seconds NUMERIC
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        COUNT(*) as total_messages,
        COUNT(*) FILTER (WHERE status = 'pending') as pending_count,
        COUNT(*) FILTER (WHERE status = 'processing') as processing_count,
        COUNT(*) FILTER (WHERE status = 'sent') as sent_count,
        COUNT(*) FILTER (WHERE status = 'failed') as failed_count,
        COUNT(*) FILTER (WHERE status = 'processing' AND updated_at < NOW() - INTERVAL '10 minutes') as stuck_count,
        AVG(EXTRACT(EPOCH FROM (NOW() - updated_at))) FILTER (WHERE status = 'processing') as avg_processing_seconds
    FROM messages;
END;
$$ LANGUAGE plpgsql;

-- Create a scheduler status table for tracking system state
CREATE TABLE IF NOT EXISTS scheduler_status (
    id SERIAL PRIMARY KEY,
    is_running BOOLEAN DEFAULT FALSE,
    started_at TIMESTAMP WITH TIME ZONE NULL,
    stopped_at TIMESTAMP WITH TIME ZONE NULL,
    last_run_at TIMESTAMP WITH TIME ZONE NULL,
    last_cleanup_at TIMESTAMP WITH TIME ZONE NULL,
    messages_processed INTEGER DEFAULT 0,
    total_sent INTEGER DEFAULT 0,
    total_failed INTEGER DEFAULT 0,
    total_stuck_reset INTEGER DEFAULT 0
);

-- Insert initial scheduler status record
INSERT INTO scheduler_status (id, is_running) VALUES (1, FALSE)
ON CONFLICT (id) DO NOTHING;

-- Function to update scheduler status
CREATE OR REPLACE FUNCTION update_scheduler_status(
    running BOOLEAN,
    processed INTEGER DEFAULT 0,
    sent INTEGER DEFAULT 0,
    failed INTEGER DEFAULT 0
) RETURNS VOID AS $$
BEGIN
    UPDATE scheduler_status 
    SET 
        is_running = running,
        started_at = CASE WHEN running AND NOT is_running THEN NOW() ELSE started_at END,
        stopped_at = CASE WHEN NOT running AND is_running THEN NOW() ELSE stopped_at END,
        last_run_at = CASE WHEN processed > 0 THEN NOW() ELSE last_run_at END,
        messages_processed = messages_processed + processed,
        total_sent = total_sent + sent,
        total_failed = total_failed + failed
    WHERE id = 1;
END;
$$ LANGUAGE plpgsql;

-- Cleanup function to be run periodically
CREATE OR REPLACE FUNCTION cleanup_stuck_messages()
RETURNS INTEGER AS $$
DECLARE
    reset_count INTEGER;
BEGIN
    -- Reset stuck messages
    SELECT reset_stuck_messages(10) INTO reset_count;
    
    -- Update scheduler status
    IF reset_count > 0 THEN
        UPDATE scheduler_status 
        SET 
            last_cleanup_at = NOW(),
            total_stuck_reset = total_stuck_reset + reset_count
        WHERE id = 1;
    END IF;
    
    RETURN reset_count;
END;
$$ LANGUAGE plpgsql;

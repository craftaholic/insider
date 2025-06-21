-- init.sql - Database initialization script

-- Create messages table
CREATE TABLE IF NOT EXISTS messages (
    id BIGSERIAL PRIMARY KEY,
    phone_number VARCHAR(20) NOT NULL,
    content TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'sent', 'failed')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    sent_at TIMESTAMP WITH TIME ZONE NULL,
    message_id VARCHAR(255) NULL,
    error_message TEXT NULL
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_messages_status_created ON messages (status, created_at);
CREATE INDEX IF NOT EXISTS idx_messages_phone_number ON messages (phone_number);
CREATE INDEX IF NOT EXISTS idx_messages_message_id ON messages (message_id);
CREATE INDEX IF NOT EXISTS idx_messages_sent_at ON messages (sent_at);

-- Insert sample data for testing
INSERT INTO messages (phone_number, content, status) VALUES 
    ('+905551111111', 'Test message 1 - Insider Project', 'pending'),
    ('+905551111112', 'Test message 2 - Hello World', 'pending'),
    ('+905551111113', 'Test message 3 - Sample Content', 'pending'),
    ('+905551111114', 'Test message 4 - Demo Message', 'pending'),
    ('+905551111115', 'Test message 5 - Testing System', 'pending');

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

-- Create a function to get unsent messages (optional, can be useful)
CREATE OR REPLACE FUNCTION get_unsent_messages(batch_size INTEGER DEFAULT 2)
RETURNS TABLE (
    id BIGINT,
    phone_number VARCHAR(20),
    content TEXT,
    created_at TIMESTAMP WITH TIME ZONE
) AS $$
BEGIN
    RETURN QUERY
    SELECT m.id, m.phone_number, m.content, m.created_at
    FROM messages m
    WHERE m.status = 'pending'
    ORDER BY m.created_at ASC
    LIMIT batch_size;
END;
$$ LANGUAGE plpgsql;

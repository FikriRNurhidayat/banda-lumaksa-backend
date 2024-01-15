CREATE TABLE subscriptions (
       id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
       name VARCHAR(255) NOT NULL UNIQUE,
       fee INTEGER NOT NULL,
       subscription_type VARCHAR(255) NOT NULL,
       started_at TIMESTAMP WITH TIME ZONE NOT NULL,
       ended_at TIMESTAMP WITH TIME ZONE,
       due_at TIMESTAMP WITH TIME ZONE NOT NULL,
       created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
       updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

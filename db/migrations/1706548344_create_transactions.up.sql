CREATE TABLE transactions (
       id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
       description VARCHAR(255) NOT NULL,
       amount INTEGER NOT NULL,
       created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
       updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

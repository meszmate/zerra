
CREATE TABLE users (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid (),
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT,
    created_at TIMESTAMPTZ DEFAULT now (),
    updated_at TIMESTAMPTZ DEFAULT now ()
);

CREATE TABLE organizations (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid (),
    name VARCHAR(255) NOT NULL,
    owner UUID REFERENCES users (id) NOT NULL,
    icon UUID REFERENCES images(id) NOT NULL,
    
    created_at TIMESTAMPTZ DEFAULT now (),
    updated_at TIMESTAMPTZ DEFAULT now ()
);

CREATE TYPE file_owner_type AS ENUM ('user', 'organization');

CREATE TABLE files (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid (),
    parent_type file_owner_type NOT NULL,
    parent_id UUID NOT NULL,
    file_key TEXT NOT NULL,
    file_type TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
);

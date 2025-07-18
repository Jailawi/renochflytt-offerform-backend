CREATE SCHEMA IF NOT EXISTS renochflytt;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS renochflytt.customers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    ssn VARCHAR(11) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL,
    phone VARCHAR(15) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS renochflytt.bookings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    customer_id UUID REFERENCES renochflytt.customers(id) ON DELETE CASCADE,
    services TEXT[] NOT NULL,
    moving_date DATE,
    cleaning_date DATE,
    flexible BOOLEAN NOT NULL,
    from_address_id UUID NOT NULL,
    to_address_id UUID NOT NULL,
    message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS renochflytt.addresses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    adress VARCHAR(255) NOT NULL,
    residence_type VARCHAR(50) NOT NULL,
    floor VARCHAR(50),
    accessibility VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
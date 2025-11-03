CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  name VARCHAR(255),
  role VARCHAR(50) DEFAULT 'user',
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TYPE IF NOT EXISTS property_category AS ENUM ('buy', 'rent', 'shortlet');

CREATE TABLE IF NOT EXISTS properties (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  title VARCHAR(255) NOT NULL,
  description TEXT,
  category property_category NOT NULL,
  price NUMERIC NOT NULL,
  currency VARCHAR(10) DEFAULT 'NGN',
  address TEXT,
  area VARCHAR(100),
  bedrooms INT,
  bathrooms INT,
  furnished BOOLEAN DEFAULT false,
  party_allowed BOOLEAN DEFAULT false,
  instant_book BOOLEAN DEFAULT false,
  owner_id UUID,
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS property_images (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  property_id UUID REFERENCES properties(id) ON DELETE CASCADE,
  url TEXT NOT NULL,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS bookings (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  property_id UUID REFERENCES properties(id) ON DELETE CASCADE,
  user_id UUID,
  checkin DATE,
  checkout DATE,
  nights INT,
  guests INT,
  total_amount NUMERIC,
  status VARCHAR(50) DEFAULT 'pending',
  payment_ref TEXT,
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_properties_category_area ON properties (category, area);
CREATE INDEX IF NOT EXISTS idx_bookings_property_dates ON bookings (property_id, checkin, checkout);

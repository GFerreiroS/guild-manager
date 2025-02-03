DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_extension WHERE extname = 'uuid-ossp') THEN
        CREATE EXTENSION "uuid-ossp";
    END IF;
END $$;

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    battle_net_id VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE,
    role VARCHAR(50) NOT NULL DEFAULT 'member',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE guilds (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    realm VARCHAR(255) NOT NULL,
    faction VARCHAR(50) NOT NULL CHECK (faction IN ('alliance', 'horde')),
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE characters (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    realm VARCHAR(255) NOT NULL,
    class VARCHAR(50) NOT NULL CHECK (class IN (
        'warrior', 'paladin', 'hunter', 'rogue', 'priest', 
        'death-knight', 'shaman', 'mage', 'warlock', 
        'monk', 'druid', 'demon-hunter', 'evoker'
    )),
    spec VARCHAR(50),
    ilvl INTEGER NOT NULL,
    last_synced TIMESTAMP WITH TIME ZONE,
    user_id UUID REFERENCES users(id),
    guild_id UUID REFERENCES guilds(id) ON DELETE CASCADE,
    raid_group_id UUID,
    UNIQUE(name, realm)
);

CREATE TABLE guild_members (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    guild_id UUID REFERENCES guilds(id) ON DELETE CASCADE,
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    role VARCHAR(50),
    PRIMARY KEY (user_id, guild_id)
);

CREATE TABLE raid_group_characters (
    character_id UUID REFERENCES characters(id) ON DELETE CASCADE,
    raid_group_id UUID REFERENCES raid_groups(id) ON DELETE CASCADE,
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (character_id, raid_group_id)
);

CREATE INDEX idx_characters_guild ON characters(guild_id);

CREATE TABLE raid_groups (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    guild_id UUID REFERENCES guilds(id) ON DELETE CASCADE,
    schedule JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    raid_name VARCHAR(255) NOT NULL,
    difficulty VARCHAR(50) CHECK (difficulty IN ('normal', 'heroic', 'mythic')),
    scheduled_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_by UUID REFERENCES users(id),
    guild_id UUID REFERENCES guilds(id) ON DELETE CASCADE
);

CREATE TABLE confirmations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_id UUID REFERENCES events(id) ON DELETE CASCADE,
    character_id UUID REFERENCES characters(id) ON DELETE CASCADE,
    status VARCHAR(50) CHECK (status IN ('confirmed', 'declined', 'tentative')),
    reason TEXT,
    responded_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
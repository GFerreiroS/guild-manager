CREATE INDEX IF NOT EXISTS idx_characters_guild ON characters(guild_id);
CREATE INDEX IF NOT EXISTS idx_events_guild ON events(guild_id);
CREATE INDEX IF NOT EXISTS idx_confirmations_event ON confirmations(event_id);
CREATE INDEX IF NOT EXISTS idx_guild_members_user ON guild_members(user_id);
CREATE INDEX IF NOT EXISTS idx_guild_members_guild ON guild_members(guild_id);
CREATE INDEX IF NOT EXISTS idx_raid_group_members_char ON raid_group_characters(character_id);
CREATE INDEX IF NOT EXISTS idx_raid_group_members_group ON raid_group_characters(raid_group_id);
CREATE INDEX IF NOT EXISTS idx_events_scheduled ON events(scheduled_at);

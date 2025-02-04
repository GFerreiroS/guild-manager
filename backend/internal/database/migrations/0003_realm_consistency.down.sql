ALTER TABLE characters
DROP CONSTRAINT IF EXISTS fk_characters_realm;

ALTER TABLE guilds
DROP CONSTRAINT IF EXISTS unique_guild_realm;
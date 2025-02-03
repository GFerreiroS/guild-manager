ALTER TABLE guilds
ADD CONSTRAINT unique_guild_realm UNIQUE (realm);

ALTER TABLE characters
ADD CONSTRAINT fk_characters_realm
FOREIGN KEY (realm) REFERENCES guilds(realm)
ON DELETE CASCADE;
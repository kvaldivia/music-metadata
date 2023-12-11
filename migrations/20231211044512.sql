-- Create "artists" table
CREATE TABLE "artists" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "name" text NULL,
  "spotify_id" text NULL,
  PRIMARY KEY ("id")
);
-- Create index "artists_spotify_id_key" to table: "artists"
CREATE UNIQUE INDEX "artists_spotify_id_key" ON "artists" ("spotify_id");
-- Create index "idx_artists_deleted_at" to table: "artists"
CREATE INDEX "idx_artists_deleted_at" ON "artists" ("deleted_at");
-- Create "tracks" table
CREATE TABLE "tracks" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "spotify_id" text NULL,
  "isrc" text NULL,
  "image_uri" text NULL,
  "title" text NULL,
  "artist_id" bigint NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_artists_tracks" FOREIGN KEY ("artist_id") REFERENCES "artists" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_tracks_deleted_at" to table: "tracks"
CREATE INDEX "idx_tracks_deleted_at" ON "tracks" ("deleted_at");

-- Create "artists" table
CREATE TABLE "artists" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "name" text NULL,
  PRIMARY KEY ("id")
);
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
  "artist_ref" bigint NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_artists_songs" FOREIGN KEY ("artist_ref") REFERENCES "artists" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_tracks_deleted_at" to table: "tracks"
CREATE INDEX "idx_tracks_deleted_at" ON "tracks" ("deleted_at");

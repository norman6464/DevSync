table "spotify_recent_tracks" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigint
  }
  column "track_name" {
    null = false
    type = text
  }
  column "artist_name" {
    null = false
    type = text
  }
  column "album_name" {
    null = true
    type = text
  }
  column "album_image" {
    null = true
    type = text
  }
  column "track_url" {
    null = true
    type = text
  }
  column "played_at" {
    null = false
    type = timestamptz
  }
  column "created_at" {
    null = true
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_spotify_recent_tracks_user_id" {
    columns = [column.user_id]
  }
}

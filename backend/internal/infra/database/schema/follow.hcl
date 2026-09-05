table "follows" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "follower_id" {
    null = false
    type = bigint
  }
  column "followee_id" {
    null = false
    type = bigint
  }
  column "created_at" {
    null = true
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_follows_followee" {
    columns     = [column.followee_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  foreign_key "fk_follows_follower" {
    columns     = [column.follower_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  index "idx_follower_following" {
    unique  = true
    columns = [column.follower_id, column.followee_id]
  }
}

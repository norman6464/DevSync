-- name: DeleteSpotifyRecentTracksByUser :exec
DELETE FROM spotify_recent_tracks WHERE user_id = $1;

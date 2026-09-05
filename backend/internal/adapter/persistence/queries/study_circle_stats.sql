-- name: CountStudyCircleMembersByCircle :one
SELECT count(*) FROM study_circle_members
WHERE circle_id = $1;

-- name: CountStudyCircleCheckinsByCircle :one
SELECT count(*) FROM study_circle_checkins
WHERE circle_id = $1;

-- name: CountStudyCircleStepsByCircle :one
SELECT count(*) FROM study_circle_steps
WHERE circle_id = $1;

-- name: CountStudyCircleCompletedStepsByCircle :one
SELECT count(*) FROM study_circle_member_progresses
WHERE circle_id = $1 AND is_completed = true;

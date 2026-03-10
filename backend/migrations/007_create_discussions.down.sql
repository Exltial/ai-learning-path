-- Migration: Drop discussions system
-- Created: 2026-03-10
-- Description: Rolls back the community discussion system

-- Drop indexes first
DROP INDEX IF EXISTS idx_discussions_popular;
DROP INDEX IF EXISTS idx_discussions_hot;
DROP INDEX IF EXISTS idx_discussions_course_created;
DROP INDEX IF EXISTS idx_discussions_parent_depth;
DROP INDEX IF EXISTS idx_discussion_mentions_is_read;
DROP INDEX IF EXISTS idx_discussion_mentions_user_id;
DROP INDEX IF EXISTS idx_discussion_mentions_discussion_id;
DROP INDEX IF EXISTS idx_discussion_favorites_user_id;
DROP INDEX IF EXISTS idx_discussion_favorites_discussion_id;
DROP INDEX IF EXISTS idx_discussion_likes_user_id;
DROP INDEX IF EXISTS idx_discussion_likes_discussion_id;
DROP INDEX IF EXISTS idx_discussions_deleted_at;
DROP INDEX IF EXISTS idx_discussions_is_pinned;
DROP INDEX IF EXISTS idx_discussions_is_resolved;
DROP INDEX IF EXISTS idx_discussions_reply_count;
DROP INDEX IF EXISTS idx_discussions_upvotes;
DROP INDEX IF EXISTS idx_discussions_updated_at;
DROP INDEX IF EXISTS idx_discussions_created_at;
DROP INDEX IF EXISTS idx_discussions_parent_id;
DROP INDEX IF EXISTS idx_discussions_user_id;
DROP INDEX IF EXISTS idx_discussions_lesson_id;
DROP INDEX IF EXISTS idx_discussions_course_id;

-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS discussion_tag_mapping;
DROP TABLE IF EXISTS discussion_tags;
DROP TABLE IF EXISTS discussion_mentions;
DROP TABLE IF EXISTS discussion_favorites;
DROP TABLE IF EXISTS discussion_likes;
DROP TABLE IF EXISTS discussions;

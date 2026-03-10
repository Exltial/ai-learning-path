-- Migration: Create discussions system
-- Created: 2026-03-10
-- Description: Implements community discussion system with courses, posts, replies, likes, and favorites

-- Main discussions table (stores both threads and replies)
CREATE TABLE IF NOT EXISTS discussions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    lesson_id UUID REFERENCES lessons(id) ON DELETE SET NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    parent_id UUID REFERENCES discussions(id) ON DELETE CASCADE,
    title VARCHAR(500),
    content TEXT NOT NULL,
    content_html TEXT,
    is_resolved BOOLEAN DEFAULT FALSE,
    is_locked BOOLEAN DEFAULT FALSE,
    is_pinned BOOLEAN DEFAULT FALSE,
    is_anonymous BOOLEAN DEFAULT FALSE,
    upvotes INT DEFAULT 0,
    downvotes INT DEFAULT 0,
    reply_count INT DEFAULT 0,
    view_count INT DEFAULT 0,
    depth INT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT valid_depth CHECK (depth >= 0 AND depth <= 10),
    CONSTRAINT has_title_if_root CHECK (parent_id IS NOT NULL OR title IS NOT NULL)
);

-- Discussion likes table
CREATE TABLE IF NOT EXISTS discussion_likes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    discussion_id UUID NOT NULL REFERENCES discussions(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    like_type VARCHAR(10) NOT NULL DEFAULT 'upvote' CHECK (like_type IN ('upvote', 'downvote')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(discussion_id, user_id, like_type)
);

-- Discussion favorites/bookmarks table
CREATE TABLE IF NOT EXISTS discussion_favorites (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    discussion_id UUID NOT NULL REFERENCES discussions(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(discussion_id, user_id)
);

-- Discussion mentions table (for @user functionality)
CREATE TABLE IF NOT EXISTS discussion_mentions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    discussion_id UUID NOT NULL REFERENCES discussions(id) ON DELETE CASCADE,
    mentioned_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    mentioned_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Discussion tags table
CREATE TABLE IF NOT EXISTS discussion_tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL UNIQUE,
    color VARCHAR(7) DEFAULT '#3B82F6',
    description TEXT,
    usage_count INT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Discussion to tags mapping
CREATE TABLE IF NOT EXISTS discussion_tag_mapping (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    discussion_id UUID NOT NULL REFERENCES discussions(id) ON DELETE CASCADE,
    tag_id UUID NOT NULL REFERENCES discussion_tags(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(discussion_id, tag_id)
);

-- Indexes for common queries
CREATE INDEX idx_discussions_course_id ON discussions(course_id);
CREATE INDEX idx_discussions_lesson_id ON discussions(lesson_id);
CREATE INDEX idx_discussions_user_id ON discussions(user_id);
CREATE INDEX idx_discussions_parent_id ON discussions(parent_id);
CREATE INDEX idx_discussions_created_at ON discussions(created_at DESC);
CREATE INDEX idx_discussions_updated_at ON discussions(updated_at DESC);
CREATE INDEX idx_discussions_upvotes ON discussions(upvotes DESC);
CREATE INDEX idx_discussions_reply_count ON discussions(reply_count DESC);
CREATE INDEX idx_discussions_is_resolved ON discussions(is_resolved);
CREATE INDEX idx_discussions_is_pinned ON discussions(is_pinned DESC);
CREATE INDEX idx_discussions_deleted_at ON discussions(deleted_at) WHERE deleted_at IS NOT NULL;
CREATE INDEX idx_discussions_course_created ON discussions(course_id, created_at DESC);
CREATE INDEX idx_discussions_parent_depth ON discussions(parent_id, depth);

-- Indexes for likes and favorites
CREATE INDEX idx_discussion_likes_discussion_id ON discussion_likes(discussion_id);
CREATE INDEX idx_discussion_likes_user_id ON discussion_likes(user_id);
CREATE INDEX idx_discussion_favorites_discussion_id ON discussion_favorites(discussion_id);
CREATE INDEX idx_discussion_favorites_user_id ON discussion_favorites(user_id);

-- Indexes for mentions
CREATE INDEX idx_discussion_mentions_discussion_id ON discussion_mentions(discussion_id);
CREATE INDEX idx_discussion_mentions_user_id ON discussion_mentions(mentioned_user_id);
CREATE INDEX idx_discussion_mentions_is_read ON discussion_mentions(is_read) WHERE is_read = FALSE;

-- Composite indexes for popular discussions
CREATE INDEX idx_discussions_popular ON discussions(course_id, is_pinned DESC, is_resolved ASC, upvotes DESC, created_at DESC);
CREATE INDEX idx_discussions_hot ON discussions(course_id, created_at DESC, upvotes DESC, reply_count DESC);

-- Insert some default tags
INSERT INTO discussion_tags (name, color, description, usage_count) VALUES
    ('提问', '#EF4444', '课程相关问题', 0),
    ('分享', '#10B981', '学习心得分享', 0),
    ('讨论', '#3B82F6', '一般讨论', 0),
    ('建议', '#F59E0B', '课程改进建议', 0),
    ('已解决', '#10B981', '问题已解决', 0),
    ('精华', '#8B5CF6', '优质内容', 0)
ON CONFLICT (name) DO NOTHING;

-- Comment on tables
COMMENT ON TABLE discussions IS 'Community discussion threads and replies (nested structure)';
COMMENT ON COLUMN discussions.parent_id IS 'NULL for top-level threads, references parent discussion for replies';
COMMENT ON COLUMN discussions.depth IS 'Nesting depth (0 for top-level, 1-10 for replies)';
COMMENT ON COLUMN discussions.content_html IS 'HTML rendered version of content (for Markdown)';
COMMENT ON COLUMN discussions.is_anonymous IS 'Whether the discussion is posted anonymously';
COMMENT ON TABLE discussion_likes IS 'User likes/upvotes/downvotes on discussions';
COMMENT ON TABLE discussion_favorites IS 'User bookmarked/favorited discussions';
COMMENT ON TABLE discussion_mentions IS 'User mentions in discussions (@username)';

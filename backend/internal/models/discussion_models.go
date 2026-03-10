package models

import (
	"time"

	"github.com/google/uuid"
)

// Discussion represents a discussion thread or reply
type Discussion struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	CourseID      uuid.UUID  `json:"course_id" db:"course_id"`
	LessonID      *uuid.UUID `json:"lesson_id,omitempty" db:"lesson_id"`
	UserID        uuid.UUID  `json:"user_id" db:"user_id"`
	ParentID      *uuid.UUID `json:"parent_id,omitempty" db:"parent_id"`
	Title         *string    `json:"title,omitempty" db:"title"`
	Content       string     `json:"content" db:"content"`
	ContentHTML   *string    `json:"content_html,omitempty" db:"content_html"`
	IsResolved    bool       `json:"is_resolved" db:"is_resolved"`
	IsLocked      bool       `json:"is_locked" db:"is_locked"`
	IsPinned      bool       `json:"is_pinned" db:"is_pinned"`
	IsAnonymous   bool       `json:"is_anonymous" db:"is_anonymous"`
	Upvotes       int        `json:"upvotes" db:"upvotes"`
	Downvotes     int        `json:"downvotes" db:"downvotes"`
	ReplyCount    int        `json:"reply_count" db:"reply_count"`
	ViewCount     int        `json:"view_count" db:"view_count"`
	Depth         int        `json:"depth" db:"depth"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
	
	// Computed fields (not stored in DB)
	User         *User        `json:"user,omitempty" db:"-"`
	Replies      []*Discussion `json:"replies,omitempty" db:"-"`
	IsLiked      *bool        `json:"is_liked,omitempty" db:"-"` // Current user's like status
	IsFavorited  *bool        `json:"is_favorited,omitempty" db:"-"` // Current user's favorite status
	Tags         []*DiscussionTag `json:"tags,omitempty" db:"-"`
}

// DiscussionLike represents a like/upvote/downvote on a discussion
type DiscussionLike struct {
	ID           uuid.UUID `json:"id" db:"id"`
	DiscussionID uuid.UUID `json:"discussion_id" db:"discussion_id"`
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	LikeType     string    `json:"like_type" db:"like_type"` // upvote, downvote
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// DiscussionFavorite represents a user's favorite/bookmark of a discussion
type DiscussionFavorite struct {
	ID           uuid.UUID `json:"id" db:"id"`
	DiscussionID uuid.UUID `json:"discussion_id" db:"discussion_id"`
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// DiscussionMention represents a user mention in a discussion
type DiscussionMention struct {
	ID            uuid.UUID `json:"id" db:"id"`
	DiscussionID  uuid.UUID `json:"discussion_id" db:"discussion_id"`
	MentionedUserID uuid.UUID `json:"mentioned_user_id" db:"mentioned_user_id"`
	MentionedBy   uuid.UUID `json:"mentioned_by" db:"mentioned_by"`
	IsRead        bool      `json:"is_read" db:"is_read"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

// DiscussionTag represents a tag for categorizing discussions
type DiscussionTag struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Color       string    `json:"color" db:"color"`
	Description string    `json:"description,omitempty" db:"description"`
	UsageCount  int       `json:"usage_count" db:"usage_count"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// DiscussionTagMapping represents the many-to-many relationship between discussions and tags
type DiscussionTagMapping struct {
	ID           uuid.UUID `json:"id" db:"id"`
	DiscussionID uuid.UUID `json:"discussion_id" db:"discussion_id"`
	TagID        uuid.UUID `json:"tag_id" db:"tag_id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// DiscussionWithStats includes additional statistics for a discussion
type DiscussionWithStats struct {
	*Discussion
	TotalReplies int     `json:"total_replies" db:"total_replies"`
	LatestReplyAt *time.Time `json:"latest_reply_at,omitempty" db:"latest_reply_at"`
	ParticipantCount int   `json:"participant_count" db:"participant_count"`
}

// CreateDiscussionRequest represents a request to create a new discussion
type CreateDiscussionRequest struct {
	CourseID    uuid.UUID `json:"course_id" binding:"required"`
	LessonID    *uuid.UUID `json:"lesson_id,omitempty"`
	Title       *string   `json:"title,omitempty" binding:"required_without=ParentID"`
	Content     string    `json:"content" binding:"required,min=1,max=50000"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
	IsAnonymous bool      `json:"is_anonymous,omitempty"`
	TagIDs      []uuid.UUID `json:"tag_ids,omitempty"`
}

// UpdateDiscussionRequest represents a request to update a discussion
type UpdateDiscussionRequest struct {
	Title       *string   `json:"title,omitempty"`
	Content     *string   `json:"content,omitempty" binding:"omitempty,min=1,max=50000"`
	IsResolved  *bool     `json:"is_resolved,omitempty"`
	IsLocked    *bool     `json:"is_locked,omitempty"`
	IsPinned    *bool     `json:"is_pinned,omitempty"`
	TagIDs      []uuid.UUID `json:"tag_ids,omitempty"`
}

// DiscussionListOptions represents options for listing discussions
type DiscussionListOptions struct {
	CourseID   uuid.UUID `json:"course_id"`
	LessonID   *uuid.UUID `json:"lesson_id,omitempty"`
	ParentID   *uuid.UUID `json:"parent_id,omitempty"` // Filter by parent (for replies)
	UserID     *uuid.UUID `json:"user_id,omitempty"`   // Filter by author
	TagIDs     []uuid.UUID `json:"tag_ids,omitempty"`
	Search     string    `json:"search,omitempty"`
	SortBy     string    `json:"sort_by,omitempty"` // created_at, updated_at, upvotes, reply_count, hot
	SortOrder  string    `json:"sort_order,omitempty"` // asc, desc
	Page       int       `json:"page"`
	Limit      int       `json:"limit"`
	WithReplies bool     `json:"with_replies"` // Include nested replies
	MaxDepth   int       `json:"max_depth"`    // Maximum reply depth to fetch
}

// DiscussionResponse represents the API response for a discussion
type DiscussionResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// DiscussionListResponse represents the API response for listing discussions
type DiscussionListResponse struct {
	Success    bool                   `json:"success"`
	Data       []*Discussion          `json:"data"`
	Pagination map[string]interface{} `json:"pagination,omitempty"`
}

// HotDiscussion represents a hot/trending discussion with calculated score
type HotDiscussion struct {
	*Discussion
	HotScore float64 `json:"hot_score" db:"hot_score"`
}

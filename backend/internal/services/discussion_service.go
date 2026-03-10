package services

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	"ai-learning-platform/internal/models"
	"ai-learning-platform/internal/repository"

	"github.com/google/uuid"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// Error definitions for discussion service
var (
	ErrDiscussionNotFound     = errors.New("discussion not found")
	ErrDiscussionLocked       = errors.New("discussion is locked")
	ErrDiscussionDeleted      = errors.New("discussion has been deleted")
	ErrInvalidDepth           = errors.New("invalid reply depth")
	ErrUnauthorized           = errors.New("unauthorized action")
	// ErrCourseNotFound is defined in course_service.go
	ErrContentTooLong         = errors.New("content is too long")
	ErrTitleRequired          = errors.New("title is required for top-level discussions")
	ErrParentNotFound         = errors.New("parent discussion not found")
	ErrDuplicateFavorite      = errors.New("discussion already favorited")
)

// MaxDepth is the maximum nesting depth for replies
const MaxDepth = 10

// MaxContentLength is the maximum content length in characters
const MaxContentLength = 50000

// DiscussionService handles discussion-related business logic
type DiscussionService struct {
	discussionRepo *repository.DiscussionRepository
	courseRepo     *repository.CourseRepository
	userRepo       *repository.UserRepository
}

// NewDiscussionService creates a new DiscussionService
func NewDiscussionService(
	discussionRepo *repository.DiscussionRepository,
	courseRepo *repository.CourseRepository,
	userRepo *repository.UserRepository,
) *DiscussionService {
	return &DiscussionService{
		discussionRepo: discussionRepo,
		courseRepo:     courseRepo,
		userRepo:       userRepo,
	}
}

// CreateDiscussion creates a new discussion or reply
func (s *DiscussionService) CreateDiscussion(ctx context.Context, req *models.CreateDiscussionRequest, userID uuid.UUID) (*models.Discussion, error) {
	// Validate content length
	if len(req.Content) > MaxContentLength {
		return nil, ErrContentTooLong
	}

	// Verify course exists
	_, err := s.courseRepo.GetByID(ctx, req.CourseID)
	if err != nil {
		return nil, ErrCourseNotFound
	}

	discussion := &models.Discussion{
		ID:          uuid.New(),
		CourseID:    req.CourseID,
		LessonID:    req.LessonID,
		UserID:      userID,
		Content:     req.Content,
		IsAnonymous: req.IsAnonymous,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// If this is a reply
	if req.ParentID != nil {
		parent, err := s.discussionRepo.GetByID(ctx, *req.ParentID)
		if err != nil {
			return nil, ErrParentNotFound
		}
		if parent.DeletedAt != nil {
			return nil, ErrDiscussionDeleted
		}
		if parent.IsLocked {
			return nil, ErrDiscussionLocked
		}

		discussion.ParentID = req.ParentID
		discussion.Depth = parent.Depth + 1

		if discussion.Depth > MaxDepth {
			return nil, ErrInvalidDepth
		}

		// Increment parent's reply count
		s.discussionRepo.IncrementReplyCount(ctx, *req.ParentID)
	} else {
		// Top-level discussion requires title
		if req.Title == nil || strings.TrimSpace(*req.Title) == "" {
			return nil, ErrTitleRequired
		}
		discussion.Title = req.Title
		discussion.Depth = 0
	}

	// Render Markdown to HTML
	discussion.ContentHTML = s.renderMarkdown(req.Content)

	// Create discussion
	if err := s.discussionRepo.Create(ctx, discussion); err != nil {
		return nil, err
	}

	// Process mentions
	s.processMentions(ctx, discussion.ID, userID, req.Content)

	// Set tags if provided
	if len(req.TagIDs) > 0 {
		s.discussionRepo.SetDiscussionTags(ctx, discussion.ID, req.TagIDs)
	}

	// Load user info
	s.loadUserInfo(ctx, discussion)

	return discussion, nil
}

// GetDiscussion retrieves a discussion by ID
func (s *DiscussionService) GetDiscussion(ctx context.Context, id uuid.UUID, currentUserID uuid.UUID) (*models.Discussion, error) {
	discussion, err := s.discussionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrDiscussionNotFound
	}
	if discussion.DeletedAt != nil {
		return nil, ErrDiscussionDeleted
	}

	// Increment view count (non-blocking)
	go func() {
		updateCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		s.discussionRepo.IncrementViewCount(updateCtx, id)
	}()

	// Load user info
	s.loadUserInfo(ctx, discussion)

	// Load tags
	tags, err := s.discussionRepo.GetDiscussionTags(ctx, id)
	if err == nil {
		discussion.Tags = tags
	}

	// Load user's like status
	if currentUserID != uuid.Nil {
		likeType, err := s.discussionRepo.GetUserLike(ctx, id, currentUserID)
		if err == nil {
			isLiked := likeType == "upvote"
			discussion.IsLiked = &isLiked
		}

		// Load user's favorite status
		isFavorited, err := s.discussionRepo.IsFavorited(ctx, id, currentUserID)
		if err == nil {
			discussion.IsFavorited = &isFavorited
		}
	}

	return discussion, nil
}

// GetDiscussionWithReplies retrieves a discussion with its replies
func (s *DiscussionService) GetDiscussionWithReplies(ctx context.Context, id uuid.UUID, maxDepth int, currentUserID uuid.UUID) (*models.Discussion, error) {
	discussion, err := s.GetDiscussion(ctx, id, currentUserID)
	if err != nil {
		return nil, err
	}

	// Load replies recursively
	replies, err := s.loadReplies(ctx, id, 1, maxDepth, currentUserID)
	if err != nil {
		return nil, err
	}
	discussion.Replies = replies

	return discussion, nil
}

// loadReplies recursively loads replies for a discussion
func (s *DiscussionService) loadReplies(ctx context.Context, parentID uuid.UUID, currentDepth, maxDepth int, currentUserID uuid.UUID) ([]*models.Discussion, error) {
	if currentDepth > maxDepth {
		return []*models.Discussion{}, nil
	}

	replies, err := s.discussionRepo.GetReplies(ctx, parentID, maxDepth)
	if err != nil {
		return nil, err
	}

	// Load additional info for each reply
	for _, reply := range replies {
		s.loadUserInfo(ctx, reply)

		// Load user's like status
		if currentUserID != uuid.Nil {
			likeType, err := s.discussionRepo.GetUserLike(ctx, reply.ID, currentUserID)
			if err == nil {
				isLiked := likeType == "upvote"
				reply.IsLiked = &isLiked
			}

			isFavorited, err := s.discussionRepo.IsFavorited(ctx, reply.ID, currentUserID)
			if err == nil {
				reply.IsFavorited = &isFavorited
			}
		}

		// Recursively load nested replies
		if reply.Depth < maxDepth {
			nestedReplies, err := s.loadReplies(ctx, reply.ID, currentDepth+1, maxDepth, currentUserID)
			if err == nil && len(nestedReplies) > 0 {
				reply.Replies = nestedReplies
			}
		}
	}

	return replies, nil
}

// ListDiscussions retrieves discussions with pagination and filters
func (s *DiscussionService) ListDiscussions(ctx context.Context, opts *models.DiscussionListOptions, currentUserID uuid.UUID) ([]*models.Discussion, int, error) {
	discussions, total, err := s.discussionRepo.List(ctx, opts)
	if err != nil {
		return nil, 0, err
	}

	// Load additional info for each discussion
	for _, discussion := range discussions {
		s.loadUserInfo(ctx, discussion)

		// Load tags
		tags, err := s.discussionRepo.GetDiscussionTags(ctx, discussion.ID)
		if err == nil {
			discussion.Tags = tags
		}

		// Load user's like and favorite status
		if currentUserID != uuid.Nil {
			likeType, err := s.discussionRepo.GetUserLike(ctx, discussion.ID, currentUserID)
			if err == nil {
				isLiked := likeType == "upvote"
				discussion.IsLiked = &isLiked
			}

			isFavorited, err := s.discussionRepo.IsFavorited(ctx, discussion.ID, currentUserID)
			if err == nil {
				discussion.IsFavorited = &isFavorited
			}
		}
	}

	return discussions, total, nil
}

// GetHotDiscussions retrieves hot/trending discussions for a course
func (s *DiscussionService) GetHotDiscussions(ctx context.Context, courseID uuid.UUID, limit int, currentUserID uuid.UUID) ([]*models.HotDiscussion, error) {
	hotDiscussions, err := s.discussionRepo.GetHotDiscussions(ctx, courseID, limit)
	if err != nil {
		return nil, err
	}

	// Load additional info
	for _, hot := range hotDiscussions {
		s.loadUserInfo(ctx, hot.Discussion)

		tags, err := s.discussionRepo.GetDiscussionTags(ctx, hot.ID)
		if err == nil {
			hot.Tags = tags
		}

		if currentUserID != uuid.Nil {
			likeType, err := s.discussionRepo.GetUserLike(ctx, hot.ID, currentUserID)
			if err == nil {
				isLiked := likeType == "upvote"
				hot.IsLiked = &isLiked
			}

			isFavorited, err := s.discussionRepo.IsFavorited(ctx, hot.ID, currentUserID)
			if err == nil {
				hot.IsFavorited = &isFavorited
			}
		}
	}

	return hotDiscussions, nil
}

// UpdateDiscussion updates an existing discussion
func (s *DiscussionService) UpdateDiscussion(ctx context.Context, id uuid.UUID, req *models.UpdateDiscussionRequest, userID uuid.UUID) (*models.Discussion, error) {
	discussion, err := s.discussionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrDiscussionNotFound
	}
	if discussion.DeletedAt != nil {
		return nil, ErrDiscussionDeleted
	}

	// Check ownership
	if discussion.UserID != userID {
		return nil, ErrUnauthorized
	}

	// Update fields if provided
	if req.Title != nil {
		discussion.Title = req.Title
	}
	if req.Content != nil {
		if len(*req.Content) > MaxContentLength {
			return nil, ErrContentTooLong
		}
		discussion.Content = *req.Content
		discussion.ContentHTML = s.renderMarkdown(*req.Content)
		
		// Process mentions in updated content
		s.processMentions(ctx, discussion.ID, userID, *req.Content)
	}
	if req.IsResolved != nil {
		discussion.IsResolved = *req.IsResolved
	}
	if req.IsLocked != nil {
		discussion.IsLocked = *req.IsLocked
	}
	if req.IsPinned != nil {
		discussion.IsPinned = *req.IsPinned
	}

	discussion.UpdatedAt = time.Now()

	if err := s.discussionRepo.Update(ctx, discussion); err != nil {
		return nil, err
	}

	// Update tags if provided
	if req.TagIDs != nil {
		s.discussionRepo.SetDiscussionTags(ctx, discussion.ID, req.TagIDs)
	}

	s.loadUserInfo(ctx, discussion)
	return discussion, nil
}

// DeleteDiscussion soft-deletes a discussion
func (s *DiscussionService) DeleteDiscussion(ctx context.Context, id uuid.UUID, userID uuid.UUID, isAdmin bool) error {
	discussion, err := s.discussionRepo.GetByID(ctx, id)
	if err != nil {
		return ErrDiscussionNotFound
	}

	// Check ownership or admin
	if discussion.UserID != userID && !isAdmin {
		return ErrUnauthorized
	}

	// If this is a reply, decrement parent's reply count
	if discussion.ParentID != nil {
		s.discussionRepo.DecrementReplyCount(ctx, *discussion.ParentID)
	}

	return s.discussionRepo.Delete(ctx, id)
}

// ToggleLike toggles a like/upvote/downvote on a discussion
func (s *DiscussionService) ToggleLike(ctx context.Context, discussionID, userID uuid.UUID, likeType string) error {
	discussion, err := s.discussionRepo.GetByID(ctx, discussionID)
	if err != nil {
		return ErrDiscussionNotFound
	}
	if discussion.DeletedAt != nil {
		return ErrDiscussionDeleted
	}

	// Validate like type
	if likeType != "upvote" && likeType != "downvote" {
		return errors.New("invalid like type")
	}

	// Check if user already liked with this type
	existingLike, err := s.discussionRepo.GetUserLike(ctx, discussionID, userID)
	if err == nil && existingLike == likeType {
		// Remove like (toggle off)
		return s.discussionRepo.RemoveLike(ctx, discussionID, userID, likeType)
	}

	// Remove existing like of different type
	if existingLike != "" {
		s.discussionRepo.RemoveLike(ctx, discussionID, userID, existingLike)
	}

	// Add new like
	like := &models.DiscussionLike{
		ID:           uuid.New(),
		DiscussionID: discussionID,
		UserID:       userID,
		LikeType:     likeType,
		CreatedAt:    time.Now(),
	}

	if err := s.discussionRepo.AddLike(ctx, like); err != nil {
		return err
	}

	// Update vote counts
	return s.discussionRepo.UpdateVoteCounts(ctx, discussionID)
}

// ToggleFavorite toggles a discussion favorite
func (s *DiscussionService) ToggleFavorite(ctx context.Context, discussionID, userID uuid.UUID) (bool, error) {
	discussion, err := s.discussionRepo.GetByID(ctx, discussionID)
	if err != nil {
		return false, ErrDiscussionNotFound
	}
	if discussion.DeletedAt != nil {
		return false, ErrDiscussionDeleted
	}

	// Check if already favorited
	isFavorited, err := s.discussionRepo.IsFavorited(ctx, discussionID, userID)
	if err != nil {
		return false, err
	}

	if isFavorited {
		// Remove favorite
		err := s.discussionRepo.RemoveFavorite(ctx, discussionID, userID)
		return false, err
	} else {
		// Add favorite
		favorite := &models.DiscussionFavorite{
			ID:           uuid.New(),
			DiscussionID: discussionID,
			UserID:       userID,
			CreatedAt:    time.Now(),
		}
		err := s.discussionRepo.AddFavorite(ctx, favorite)
		return true, err
	}
}

// GetUserFavorites retrieves a user's favorited discussions
func (s *DiscussionService) GetUserFavorites(ctx context.Context, userID uuid.UUID, page, limit int) ([]*models.Discussion, int, error) {
	discussions, total, err := s.discussionRepo.GetUserFavorites(ctx, userID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	// Load additional info
	for _, discussion := range discussions {
		s.loadUserInfo(ctx, discussion)
	}

	return discussions, total, nil
}

// GetUnreadMentions retrieves unread mentions for a user
func (s *DiscussionService) GetUnreadMentions(ctx context.Context, userID uuid.UUID) ([]*models.DiscussionMention, error) {
	return s.discussionRepo.GetUnreadMentions(ctx, userID)
}

// MarkMentionAsRead marks a mention as read
func (s *DiscussionService) MarkMentionAsRead(ctx context.Context, mentionID uuid.UUID) error {
	return s.discussionRepo.MarkMentionAsRead(ctx, mentionID)
}

// GetAllTags retrieves all available tags
func (s *DiscussionService) GetAllTags(ctx context.Context) ([]*models.DiscussionTag, error) {
	return s.discussionRepo.GetAllTags(ctx)
}

// loadUserInfo loads user information for a discussion
func (s *DiscussionService) loadUserInfo(ctx context.Context, discussion *models.Discussion) {
	if discussion.IsAnonymous {
		discussion.User = &models.User{
			ID:        uuid.Nil,
			Username:  "匿名用户",
			AvatarURL: "",
		}
		return
	}

	user, err := s.userRepo.GetByID(ctx, discussion.UserID)
	if err == nil {
		discussion.User = user
	}
}

// renderMarkdown converts Markdown to HTML
func (s *DiscussionService) renderMarkdown(content string) *string {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,
			extension.Linkify,
			extension.Strikethrough,
			extension.Table,
			extension.TaskList,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)

	var buf strings.Builder
	if err := md.Convert([]byte(content), &buf); err != nil {
		// If conversion fails, return original content
		return &content
	}

	htmlContent := buf.String()
	return &htmlContent
}

// processMentions extracts and processes @mentions from content
func (s *DiscussionService) processMentions(ctx context.Context, discussionID, mentionedBy uuid.UUID, content string) {
	// Regex to match @username patterns
	mentionRegex := regexp.MustCompile(`@([\w\u4e00-\u9fa5]+)`)
	matches := mentionRegex.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		username := match[1]
		
		// Find user by username
		user, err := s.userRepo.GetByUsername(ctx, username)
		if err != nil {
			continue // User not found, skip
		}

		// Create mention record
		mention := &models.DiscussionMention{
			ID:              uuid.New(),
			DiscussionID:    discussionID,
			MentionedUserID: user.ID,
			MentionedBy:     mentionedBy,
			CreatedAt:       time.Now(),
		}

		s.discussionRepo.AddMention(ctx, mention)
	}
}

// ResolveDiscussion marks a discussion as resolved
func (s *DiscussionService) ResolveDiscussion(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*models.Discussion, error) {
	discussion, err := s.discussionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrDiscussionNotFound
	}

	// Only author or admin can resolve
	if discussion.UserID != userID {
		return nil, ErrUnauthorized
	}

	discussion.IsResolved = true
	discussion.UpdatedAt = time.Now()

	if err := s.discussionRepo.Update(ctx, discussion); err != nil {
		return nil, err
	}

	return discussion, nil
}

// PinDiscussion pins/unpins a discussion (admin only)
func (s *DiscussionService) PinDiscussion(ctx context.Context, id uuid.UUID, pinned bool) (*models.Discussion, error) {
	discussion, err := s.discussionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrDiscussionNotFound
	}

	discussion.IsPinned = pinned
	discussion.UpdatedAt = time.Now()

	if err := s.discussionRepo.Update(ctx, discussion); err != nil {
		return nil, err
	}

	return discussion, nil
}

// LockDiscussion locks/unlocks a discussion (admin only)
func (s *DiscussionService) LockDiscussion(ctx context.Context, id uuid.UUID, locked bool) (*models.Discussion, error) {
	discussion, err := s.discussionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrDiscussionNotFound
	}

	discussion.IsLocked = locked
	discussion.UpdatedAt = time.Now()

	if err := s.discussionRepo.Update(ctx, discussion); err != nil {
		return nil, err
	}

	return discussion, nil
}

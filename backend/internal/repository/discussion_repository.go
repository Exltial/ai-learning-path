package repository

import (
	"context"
	"time"

	"ai-learning-platform/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DiscussionRepository handles database operations for discussions
type DiscussionRepository struct {
	db *pgxpool.Pool
}

// NewDiscussionRepository creates a new DiscussionRepository
func NewDiscussionRepository(db *pgxpool.Pool) *DiscussionRepository {
	return &DiscussionRepository{db: db}
}

// Create creates a new discussion
func (r *DiscussionRepository) Create(ctx context.Context, discussion *models.Discussion) error {
	query := `
		INSERT INTO discussions (id, course_id, lesson_id, user_id, parent_id, title, content, 
			content_html, is_anonymous, depth, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err := r.db.Exec(ctx, query,
		discussion.ID,
		discussion.CourseID,
		discussion.LessonID,
		discussion.UserID,
		discussion.ParentID,
		discussion.Title,
		discussion.Content,
		discussion.ContentHTML,
		discussion.IsAnonymous,
		discussion.Depth,
		discussion.CreatedAt,
		discussion.UpdatedAt,
	)
	return err
}

// GetByID retrieves a discussion by ID
func (r *DiscussionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Discussion, error) {
	query := `
		SELECT id, course_id, lesson_id, user_id, parent_id, title, content, content_html,
			is_resolved, is_locked, is_pinned, is_anonymous, upvotes, downvotes, reply_count,
			view_count, depth, created_at, updated_at, deleted_at
		FROM discussions WHERE id = $1 AND deleted_at IS NULL
	`
	discussion := &models.Discussion{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&discussion.ID,
		&discussion.CourseID,
		&discussion.LessonID,
		&discussion.UserID,
		&discussion.ParentID,
		&discussion.Title,
		&discussion.Content,
		&discussion.ContentHTML,
		&discussion.IsResolved,
		&discussion.IsLocked,
		&discussion.IsPinned,
		&discussion.IsAnonymous,
		&discussion.Upvotes,
		&discussion.Downvotes,
		&discussion.ReplyCount,
		&discussion.ViewCount,
		&discussion.Depth,
		&discussion.CreatedAt,
		&discussion.UpdatedAt,
		&discussion.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return discussion, nil
}

// List retrieves discussions with pagination and filters
func (r *DiscussionRepository) List(ctx context.Context, opts *models.DiscussionListOptions) ([]*models.Discussion, int, error) {
	// Build count query
	countQuery := `SELECT COUNT(*) FROM discussions WHERE deleted_at IS NULL AND parent_id IS NULL`
	var args []interface{}
	argIndex := 1

	if opts.CourseID != uuid.Nil {
		countQuery += ` AND course_id = $` + string(rune(argIndex+'0'-1))
		args = append(args, opts.CourseID)
		argIndex++
	}
	if opts.LessonID != nil && *opts.LessonID != uuid.Nil {
		countQuery += ` AND lesson_id = $` + string(rune(argIndex+'0'-1))
		args = append(args, *opts.LessonID)
		argIndex++
	}
	if opts.Search != "" {
		countQuery += ` AND (title ILIKE $` + string(rune(argIndex+'0'-1)) + ` OR content ILIKE $` + string(rune(argIndex+'0'-1)) + `)`
		args = append(args, "%"+opts.Search+"%", "%"+opts.Search+"%")
		argIndex++
	}

	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Build data query
	offset := (opts.Page - 1) * opts.Limit
	query := `
		SELECT id, course_id, lesson_id, user_id, parent_id, title, content, content_html,
			is_resolved, is_locked, is_pinned, is_anonymous, upvotes, downvotes, reply_count,
			view_count, depth, created_at, updated_at, deleted_at
		FROM discussions WHERE deleted_at IS NULL AND parent_id IS NULL
	`

	if opts.CourseID != uuid.Nil {
		query += ` AND course_id = $` + string(rune(argIndex+'0'-1))
		args = append(args, opts.CourseID)
		argIndex++
	}
	if opts.LessonID != nil && *opts.LessonID != uuid.Nil {
		query += ` AND lesson_id = $` + string(rune(argIndex+'0'-1))
		args = append(args, *opts.LessonID)
		argIndex++
	}
	if opts.Search != "" {
		query += ` AND (title ILIKE $` + string(rune(argIndex+'0'-1)) + ` OR content ILIKE $` + string(rune(argIndex+'0'-1)) + `)`
		args = append(args, "%"+opts.Search+"%", "%"+opts.Search+"%")
		argIndex++
	}

	// Sorting
	sortColumn := "created_at"
	sortOrder := "DESC"
	if opts.SortBy != "" {
		switch opts.SortBy {
		case "upvotes":
			sortColumn = "upvotes"
		case "reply_count":
			sortColumn = "reply_count"
		case "updated_at":
			sortColumn = "updated_at"
		case "hot":
			sortColumn = "created_at" // Hot score is calculated separately
		}
	}
	if opts.SortOrder != "" && (opts.SortOrder == "asc" || opts.SortOrder == "ASC") {
		sortOrder = "ASC"
	}

	query += ` ORDER BY is_pinned DESC, ` + sortColumn + ` ` + sortOrder
	query += ` LIMIT $` + string(rune(argIndex+'0'-1)) + ` OFFSET $` + string(rune(argIndex+'0'))
	args = append(args, opts.Limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	discussions := make([]*models.Discussion, 0)
	for rows.Next() {
		discussion := &models.Discussion{}
		err := rows.Scan(
			&discussion.ID,
			&discussion.CourseID,
			&discussion.LessonID,
			&discussion.UserID,
			&discussion.ParentID,
			&discussion.Title,
			&discussion.Content,
			&discussion.ContentHTML,
			&discussion.IsResolved,
			&discussion.IsLocked,
			&discussion.IsPinned,
			&discussion.IsAnonymous,
			&discussion.Upvotes,
			&discussion.Downvotes,
			&discussion.ReplyCount,
			&discussion.ViewCount,
			&discussion.Depth,
			&discussion.CreatedAt,
			&discussion.UpdatedAt,
			&discussion.DeletedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		discussions = append(discussions, discussion)
	}

	return discussions, total, rows.Err()
}

// Update updates an existing discussion
func (r *DiscussionRepository) Update(ctx context.Context, discussion *models.Discussion) error {
	query := `
		UPDATE discussions SET 
			title = $2, content = $3, content_html = $4, is_resolved = $5, 
			is_locked = $6, is_pinned = $7, updated_at = $8
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query,
		discussion.ID,
		discussion.Title,
		discussion.Content,
		discussion.ContentHTML,
		discussion.IsResolved,
		discussion.IsLocked,
		discussion.IsPinned,
		time.Now(),
	)
	return err
}

// Delete soft-deletes a discussion
func (r *DiscussionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE discussions SET deleted_at = $2 WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id, time.Now())
	return err
}

// HardDelete permanently deletes a discussion (use with caution)
func (r *DiscussionRepository) HardDelete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM discussions WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

// IncrementViewCount increments the view count of a discussion
func (r *DiscussionRepository) IncrementViewCount(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE discussions SET view_count = view_count + 1 WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

// IncrementReplyCount increments the reply count of a parent discussion
func (r *DiscussionRepository) IncrementReplyCount(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE discussions SET reply_count = reply_count + 1 WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

// DecrementReplyCount decrements the reply count of a parent discussion
func (r *DiscussionRepository) DecrementReplyCount(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE discussions SET reply_count = GREATEST(reply_count - 1, 0) WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

// GetReplies retrieves replies for a discussion
func (r *DiscussionRepository) GetReplies(ctx context.Context, parentID uuid.UUID, maxDepth int) ([]*models.Discussion, error) {
	query := `
		SELECT id, course_id, lesson_id, user_id, parent_id, title, content, content_html,
			is_resolved, is_locked, is_pinned, is_anonymous, upvotes, downvotes, reply_count,
			view_count, depth, created_at, updated_at, deleted_at
		FROM discussions 
		WHERE parent_id = $1 AND deleted_at IS NULL AND depth <= $2
		ORDER BY created_at ASC
	`
	rows, err := r.db.Query(ctx, query, parentID, maxDepth)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	replies := make([]*models.Discussion, 0)
	for rows.Next() {
		discussion := &models.Discussion{}
		err := rows.Scan(
			&discussion.ID,
			&discussion.CourseID,
			&discussion.LessonID,
			&discussion.UserID,
			&discussion.ParentID,
			&discussion.Title,
			&discussion.Content,
			&discussion.ContentHTML,
			&discussion.IsResolved,
			&discussion.IsLocked,
			&discussion.IsPinned,
			&discussion.IsAnonymous,
			&discussion.Upvotes,
			&discussion.Downvotes,
			&discussion.ReplyCount,
			&discussion.ViewCount,
			&discussion.Depth,
			&discussion.CreatedAt,
			&discussion.UpdatedAt,
			&discussion.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		replies = append(replies, discussion)
	}

	return replies, rows.Err()
}

// GetByCourse retrieves discussions for a specific course
func (r *DiscussionRepository) GetByCourse(ctx context.Context, courseID uuid.UUID, page, limit int) ([]*models.Discussion, int, error) {
	opts := &models.DiscussionListOptions{
		CourseID:  courseID,
		Page:      page,
		Limit:     limit,
		SortBy:    "created_at",
		SortOrder: "DESC",
	}
	return r.List(ctx, opts)
}

// GetHotDiscussions retrieves hot/trending discussions for a course
func (r *DiscussionRepository) GetHotDiscussions(ctx context.Context, courseID uuid.UUID, limit int) ([]*models.HotDiscussion, error) {
	query := `
		SELECT 
			d.id, d.course_id, d.lesson_id, d.user_id, d.parent_id, d.title, d.content, d.content_html,
			d.is_resolved, d.is_locked, d.is_pinned, d.is_anonymous, d.upvotes, d.downvotes, 
			d.reply_count, d.view_count, d.depth, d.created_at, d.updated_at, d.deleted_at,
			((d.upvotes - d.downvotes) * 2 + d.reply_count * 3 + 
			 EXTRACT(EPOCH FROM (NOW() - d.created_at)) / 3600) as hot_score
		FROM discussions d
		WHERE d.course_id = $1 AND d.deleted_at IS NULL AND d.parent_id IS NULL
		ORDER BY hot_score DESC
		LIMIT $2
	`
	rows, err := r.db.Query(ctx, query, courseID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	hotDiscussions := make([]*models.HotDiscussion, 0)
	for rows.Next() {
		discussion := &models.Discussion{}
		hotScore := 0.0
		err := rows.Scan(
			&discussion.ID,
			&discussion.CourseID,
			&discussion.LessonID,
			&discussion.UserID,
			&discussion.ParentID,
			&discussion.Title,
			&discussion.Content,
			&discussion.ContentHTML,
			&discussion.IsResolved,
			&discussion.IsLocked,
			&discussion.IsPinned,
			&discussion.IsAnonymous,
			&discussion.Upvotes,
			&discussion.Downvotes,
			&discussion.ReplyCount,
			&discussion.ViewCount,
			&discussion.Depth,
			&discussion.CreatedAt,
			&discussion.UpdatedAt,
			&discussion.DeletedAt,
			&hotScore,
		)
		if err != nil {
			return nil, err
		}
		hotDiscussions = append(hotDiscussions, &models.HotDiscussion{
			Discussion: discussion,
			HotScore:   hotScore,
		})
	}

	return hotDiscussions, rows.Err()
}

// AddLike adds a like/upvote/downvote to a discussion
func (r *DiscussionRepository) AddLike(ctx context.Context, like *models.DiscussionLike) error {
	query := `
		INSERT INTO discussion_likes (id, discussion_id, user_id, like_type, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (discussion_id, user_id, like_type) DO NOTHING
	`
	_, err := r.db.Exec(ctx, query,
		like.ID,
		like.DiscussionID,
		like.UserID,
		like.LikeType,
		like.CreatedAt,
	)
	return err
}

// RemoveLike removes a like from a discussion
func (r *DiscussionRepository) RemoveLike(ctx context.Context, discussionID, userID uuid.UUID, likeType string) error {
	query := `DELETE FROM discussion_likes WHERE discussion_id = $1 AND user_id = $2 AND like_type = $3`
	_, err := r.db.Exec(ctx, query, discussionID, userID, likeType)
	return err
}

// GetUserLike retrieves a user's like status for a discussion
func (r *DiscussionRepository) GetUserLike(ctx context.Context, discussionID, userID uuid.UUID) (string, error) {
	query := `SELECT like_type FROM discussion_likes WHERE discussion_id = $1 AND user_id = $2`
	var likeType string
	err := r.db.QueryRow(ctx, query, discussionID, userID).Scan(&likeType)
	if err != nil {
		return "", err
	}
	return likeType, nil
}

// UpdateVoteCounts updates the upvote/downvote counts for a discussion
func (r *DiscussionRepository) UpdateVoteCounts(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE discussions SET 
			upvotes = (SELECT COUNT(*) FROM discussion_likes WHERE discussion_id = $1 AND like_type = 'upvote'),
			downvotes = (SELECT COUNT(*) FROM discussion_likes WHERE discussion_id = $1 AND like_type = 'downvote')
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

// AddFavorite adds a discussion to user's favorites
func (r *DiscussionRepository) AddFavorite(ctx context.Context, favorite *models.DiscussionFavorite) error {
	query := `
		INSERT INTO discussion_favorites (id, discussion_id, user_id, created_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (discussion_id, user_id) DO NOTHING
	`
	_, err := r.db.Exec(ctx, query,
		favorite.ID,
		favorite.DiscussionID,
		favorite.UserID,
		favorite.CreatedAt,
	)
	return err
}

// RemoveFavorite removes a discussion from user's favorites
func (r *DiscussionRepository) RemoveFavorite(ctx context.Context, discussionID, userID uuid.UUID) error {
	query := `DELETE FROM discussion_favorites WHERE discussion_id = $1 AND user_id = $2`
	_, err := r.db.Exec(ctx, query, discussionID, userID)
	return err
}

// IsFavorited checks if a discussion is favorited by a user
func (r *DiscussionRepository) IsFavorited(ctx context.Context, discussionID, userID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM discussion_favorites WHERE discussion_id = $1 AND user_id = $2)`
	var exists bool
	err := r.db.QueryRow(ctx, query, discussionID, userID).Scan(&exists)
	return exists, err
}

// GetUserFavorites retrieves a user's favorited discussions
func (r *DiscussionRepository) GetUserFavorites(ctx context.Context, userID uuid.UUID, page, limit int) ([]*models.Discussion, int, error) {
	offset := (page - 1) * limit
	
	countQuery := `SELECT COUNT(*) FROM discussion_favorites df JOIN discussions d ON df.discussion_id = d.id WHERE df.user_id = $1 AND d.deleted_at IS NULL`
	var total int
	err := r.db.QueryRow(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT d.id, d.course_id, d.lesson_id, d.user_id, d.parent_id, d.title, d.content, d.content_html,
			d.is_resolved, d.is_locked, d.is_pinned, d.is_anonymous, d.upvotes, d.downvotes, 
			d.reply_count, d.view_count, d.depth, d.created_at, d.updated_at, d.deleted_at
		FROM discussion_favorites df
		JOIN discussions d ON df.discussion_id = d.id
		WHERE df.user_id = $1 AND d.deleted_at IS NULL
		ORDER BY df.created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	discussions := make([]*models.Discussion, 0)
	for rows.Next() {
		discussion := &models.Discussion{}
		err := rows.Scan(
			&discussion.ID,
			&discussion.CourseID,
			&discussion.LessonID,
			&discussion.UserID,
			&discussion.ParentID,
			&discussion.Title,
			&discussion.Content,
			&discussion.ContentHTML,
			&discussion.IsResolved,
			&discussion.IsLocked,
			&discussion.IsPinned,
			&discussion.IsAnonymous,
			&discussion.Upvotes,
			&discussion.Downvotes,
			&discussion.ReplyCount,
			&discussion.ViewCount,
			&discussion.Depth,
			&discussion.CreatedAt,
			&discussion.UpdatedAt,
			&discussion.DeletedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		discussions = append(discussions, discussion)
	}

	return discussions, total, rows.Err()
}

// AddMention adds a user mention to a discussion
func (r *DiscussionRepository) AddMention(ctx context.Context, mention *models.DiscussionMention) error {
	query := `
		INSERT INTO discussion_mentions (id, discussion_id, mentioned_user_id, mentioned_by, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT DO NOTHING
	`
	_, err := r.db.Exec(ctx, query,
		mention.ID,
		mention.DiscussionID,
		mention.MentionedUserID,
		mention.MentionedBy,
		mention.CreatedAt,
	)
	return err
}

// GetUnreadMentions retrieves unread mentions for a user
func (r *DiscussionRepository) GetUnreadMentions(ctx context.Context, userID uuid.UUID) ([]*models.DiscussionMention, error) {
	query := `
		SELECT id, discussion_id, mentioned_user_id, mentioned_by, is_read, created_at
		FROM discussion_mentions
		WHERE mentioned_user_id = $1 AND is_read = FALSE
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mentions := make([]*models.DiscussionMention, 0)
	for rows.Next() {
		mention := &models.DiscussionMention{}
		err := rows.Scan(
			&mention.ID,
			&mention.DiscussionID,
			&mention.MentionedUserID,
			&mention.MentionedBy,
			&mention.IsRead,
			&mention.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		mentions = append(mentions, mention)
	}

	return mentions, rows.Err()
}

// MarkMentionAsRead marks a mention as read
func (r *DiscussionRepository) MarkMentionAsRead(ctx context.Context, mentionID uuid.UUID) error {
	query := `UPDATE discussion_mentions SET is_read = TRUE WHERE id = $1`
	_, err := r.db.Exec(ctx, query, mentionID)
	return err
}

// GetOrCreateTag gets an existing tag or creates a new one
func (r *DiscussionRepository) GetOrCreateTag(ctx context.Context, name string) (*models.DiscussionTag, error) {
	// Try to get existing tag
	query := `SELECT id, name, color, description, usage_count, created_at FROM discussion_tags WHERE name = $1`
	tag := &models.DiscussionTag{}
	err := r.db.QueryRow(ctx, query, name).Scan(
		&tag.ID,
		&tag.Name,
		&tag.Color,
		&tag.Description,
		&tag.UsageCount,
		&tag.CreatedAt,
	)
	if err == nil {
		return tag, nil
	}

	// Create new tag
	tag.ID = uuid.New()
	tag.Name = name
	tag.Color = "#3B82F6" // Default blue
	tag.CreatedAt = time.Now()
	
	insertQuery := `INSERT INTO discussion_tags (id, name, color, created_at) VALUES ($1, $2, $3, $4)`
	_, err = r.db.Exec(ctx, insertQuery, tag.ID, tag.Name, tag.Color, tag.CreatedAt)
	if err != nil {
		return nil, err
	}

	return tag, nil
}

// GetTags retrieves tags by IDs
func (r *DiscussionRepository) GetTags(ctx context.Context, tagIDs []uuid.UUID) ([]*models.DiscussionTag, error) {
	if len(tagIDs) == 0 {
		return []*models.DiscussionTag{}, nil
	}

	query := `SELECT id, name, color, description, usage_count, created_at FROM discussion_tags WHERE id = ANY($1)`
	rows, err := r.db.Query(ctx, query, tagIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := make([]*models.DiscussionTag, 0)
	for rows.Next() {
		tag := &models.DiscussionTag{}
		err := rows.Scan(
			&tag.ID,
			&tag.Name,
			&tag.Color,
			&tag.Description,
			&tag.UsageCount,
			&tag.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, rows.Err()
}

// GetAllTags retrieves all tags
func (r *DiscussionRepository) GetAllTags(ctx context.Context) ([]*models.DiscussionTag, error) {
	query := `SELECT id, name, color, description, usage_count, created_at FROM discussion_tags ORDER BY usage_count DESC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := make([]*models.DiscussionTag, 0)
	for rows.Next() {
		tag := &models.DiscussionTag{}
		err := rows.Scan(
			&tag.ID,
			&tag.Name,
			&tag.Color,
			&tag.Description,
			&tag.UsageCount,
			&tag.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, rows.Err()
}

// SetDiscussionTags sets tags for a discussion
func (r *DiscussionRepository) SetDiscussionTags(ctx context.Context, discussionID uuid.UUID, tagIDs []uuid.UUID) error {
	// Delete existing tags
	deleteQuery := `DELETE FROM discussion_tag_mapping WHERE discussion_id = $1`
	_, err := r.db.Exec(ctx, deleteQuery, discussionID)
	if err != nil {
		return err
	}

	if len(tagIDs) == 0 {
		return nil
	}

	// Insert new tags
	for _, tagID := range tagIDs {
		mapping := &models.DiscussionTagMapping{
			ID:           uuid.New(),
			DiscussionID: discussionID,
			TagID:        tagID,
			CreatedAt:    time.Now(),
		}
		insertQuery := `INSERT INTO discussion_tag_mapping (id, discussion_id, tag_id, created_at) VALUES ($1, $2, $3, $4)`
		_, err := r.db.Exec(ctx, insertQuery, mapping.ID, mapping.DiscussionID, mapping.TagID, mapping.CreatedAt)
		if err != nil {
			return err
		}

		// Increment tag usage count
		updateQuery := `UPDATE discussion_tags SET usage_count = usage_count + 1 WHERE id = $1`
		r.db.Exec(ctx, updateQuery, tagID)
	}

	return nil
}

// GetDiscussionTags retrieves tags for a discussion
func (r *DiscussionRepository) GetDiscussionTags(ctx context.Context, discussionID uuid.UUID) ([]*models.DiscussionTag, error) {
	query := `
		SELECT t.id, t.name, t.color, t.description, t.usage_count, t.created_at
		FROM discussion_tags t
		JOIN discussion_tag_mapping m ON t.id = m.tag_id
		WHERE m.discussion_id = $1
	`
	rows, err := r.db.Query(ctx, query, discussionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := make([]*models.DiscussionTag, 0)
	for rows.Next() {
		tag := &models.DiscussionTag{}
		err := rows.Scan(
			&tag.ID,
			&tag.Name,
			&tag.Color,
			&tag.Description,
			&tag.UsageCount,
			&tag.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, rows.Err()
}

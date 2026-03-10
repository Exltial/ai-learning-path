// Discussion types for the AI Learning Platform

export interface User {
  id: string
  username: string
  email?: string
  avatar_url?: string
  role?: string
}

export interface DiscussionTag {
  id: string
  name: string
  color: string
  description?: string
  usage_count?: number
  created_at?: string
}

export interface Discussion {
  id: string
  course_id: string
  lesson_id?: string
  user_id: string
  parent_id?: string
  title?: string
  content: string
  content_html?: string
  is_resolved: boolean
  is_locked: boolean
  is_pinned: boolean
  is_anonymous: boolean
  upvotes: number
  downvotes: number
  reply_count: number
  view_count: number
  depth: number
  created_at: string
  updated_at: string
  deleted_at?: string
  
  // Computed fields
  user?: User
  replies?: Discussion[]
  is_liked?: boolean
  is_favorited?: boolean
  tags?: DiscussionTag[]
}

export interface HotDiscussion extends Discussion {
  hot_score: number
}

export interface DiscussionLike {
  id: string
  discussion_id: string
  user_id: string
  like_type: 'upvote' | 'downvote'
  created_at: string
}

export interface DiscussionFavorite {
  id: string
  discussion_id: string
  user_id: string
  created_at: string
}

export interface DiscussionMention {
  id: string
  discussion_id: string
  mentioned_user_id: string
  mentioned_by: string
  is_read: boolean
  created_at: string
}

export interface CreateDiscussionRequest {
  course_id: string
  lesson_id?: string
  title?: string
  content: string
  parent_id?: string
  is_anonymous?: boolean
  tag_ids?: string[]
}

export interface UpdateDiscussionRequest {
  title?: string
  content?: string
  is_resolved?: boolean
  is_locked?: boolean
  is_pinned?: boolean
  tag_ids?: string[]
}

export interface DiscussionListOptions {
  course_id: string
  lesson_id?: string
  parent_id?: string
  user_id?: string
  tag_ids?: string[]
  search?: string
  sort_by?: string
  sort_order?: string
  page?: number
  limit?: number
  with_replies?: boolean
  max_depth?: number
}

export interface DiscussionListResponse {
  success: boolean
  data: Discussion[]
  pagination?: {
    page: number
    limit: number
    total: number
    total_pages: number
  }
}

export interface DiscussionResponse {
  success: boolean
  data?: Discussion
  message?: string
}

import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import {
  MessageSquare,
  Plus,
  TrendingUp,
  Clock,
  ThumbsUp,
  MessageCircle,
  Eye,
  Filter,
  Search,
  Bookmark,
  Tag,
  X,
  ChevronDown,
  Send,
  Bold,
  Italic,
  List,
  Link as LinkIcon,
  Code2,
  AtSign,
} from 'lucide-react'
import DiscussionThread from '@/components/DiscussionThread'
import { Button } from '@/components/ui/Button'
import { Input } from '@/components/ui/Input'
import { Badge } from '@/components/ui/Badge'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/Select'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/Dialog'
import { Textarea } from '@/components/ui/Textarea'
import { useAuth } from '@/contexts/AuthContext'
import { api } from '@/services/api'
import { type Discussion, type DiscussionTag } from '@/types/discussion'

interface CreateDiscussionRequest {
  course_id: string
  lesson_id?: string
  title?: string
  content: string
  parent_id?: string
  is_anonymous?: boolean
  tag_ids?: string[]
}

export default function DiscussionPage() {
  const { courseId } = useParams<{ courseId: string }>()
  const navigate = useNavigate()
  const { isAuthenticated, user } = useAuth()
  
  const [discussions, setDiscussions] = useState<Discussion[]>([])
  const [tags, setTags] = useState<DiscussionTag[]>([])
  const [loading, setLoading] = useState(true)
  const [sortBy, setSortBy] = useState('created_at')
  const [searchQuery, setSearchQuery] = useState('')
  const [selectedTags, setSelectedTags] = useState<string[]>([])
  const [page, setPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const [total, setTotal] = useState(0)
  
  // Create discussion dialog
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false)
  const [newTitle, setNewTitle] = useState('')
  const [newContent, setNewContent] = useState('')
  const [isAnonymous, setIsAnonymous] = useState(false)
  const [selectedTagIds, setSelectedTagIds] = useState<string[]>([])
  const [showMarkdownPreview, setShowMarkdownPreview] = useState(false)
  const [creating, setCreating] = useState(false)

  useEffect(() => {
    if (courseId) {
      loadDiscussions()
      loadTags()
    }
  }, [courseId, sortBy, page, selectedTags])

  const loadDiscussions = async () => {
    if (!courseId) return
    
    setLoading(true)
    try {
      const params = new URLSearchParams({
        course_id: courseId,
        page: page.toString(),
        limit: '20',
        sort_by: sortBy,
      })
      
      if (searchQuery) {
        params.append('search', searchQuery)
      }
      
      selectedTags.forEach(tagId => {
        params.append('tag_ids', tagId)
      })
      
      const response = await api.get(`/api/v1/discussions?${params.toString()}`)
      
      if (response.success) {
        setDiscussions(response.data)
        setTotalPages(response.pagination.total_pages)
        setTotal(response.pagination.total)
      }
    } catch (error) {
      console.error('Failed to load discussions:', error)
    }
    setLoading(false)
  }

  const loadTags = async () => {
    try {
      const response = await api.get('/api/v1/discussions/tags')
      if (response.success) {
        setTags(response.data)
      }
    } catch (error) {
      console.error('Failed to load tags:', error)
    }
  }

  const handleCreateDiscussion = async () => {
    if (!courseId || !newContent.trim() || !newTitle.trim()) return
    
    setCreating(true)
    try {
      const requestData: CreateDiscussionRequest = {
        course_id: courseId,
        title: newTitle,
        content: newContent,
        is_anonymous: isAnonymous,
        tag_ids: selectedTagIds,
      }
      
      const response = await api.post('/api/v1/discussions', requestData)
      
      if (response.success) {
        setIsCreateDialogOpen(false)
        resetForm()
        loadDiscussions()
      }
    } catch (error) {
      console.error('Failed to create discussion:', error)
    }
    setCreating(false)
  }

  const handleReply = async (parentId: string, content: string) => {
    if (!courseId) return
    
    try {
      const requestData: CreateDiscussionRequest = {
        course_id: courseId,
        content: content,
        parent_id: parentId,
      }
      
      const response = await api.post('/api/v1/discussions', requestData)
      
      if (response.success) {
        // Refresh discussions to show new reply
        loadDiscussions()
      }
    } catch (error) {
      console.error('Failed to post reply:', error)
      throw error
    }
  }

  const handleLike = async (discussionId: string, type: 'upvote' | 'downvote') => {
    try {
      const response = await api.post(`/api/v1/discussions/${discussionId}/like?like_type=${type}`)
      
      if (response.success) {
        // Update local state
        setDiscussions(prev => prev.map(d => {
          if (d.id === discussionId) {
            return {
              ...d,
              is_liked: d.is_liked === (type === 'upvote') ? undefined : (type === 'upvote'),
              upvotes: type === 'upvote' ? d.upvotes + 1 : d.upvotes,
            }
          }
          return d
        }))
      }
    } catch (error) {
      console.error('Failed to toggle like:', error)
    }
  }

  const handleFavorite = async (discussionId: string) => {
    try {
      const response = await api.post(`/api/v1/discussions/${discussionId}/favorite`)
      
      if (response.success) {
        // Update local state
        setDiscussions(prev => prev.map(d => {
          if (d.id === discussionId) {
            return {
              ...d,
              is_favorited: !d.is_favorited,
            }
          }
          return d
        }))
      }
    } catch (error) {
      console.error('Failed to toggle favorite:', error)
    }
  }

  const handleEdit = async (discussionId: string, content: string, title?: string) => {
    try {
      const requestData = {
        content,
        title,
      }
      
      const response = await api.put(`/api/v1/discussions/${discussionId}`, requestData)
      
      if (response.success) {
        loadDiscussions()
      }
    } catch (error) {
      console.error('Failed to update discussion:', error)
      throw error
    }
  }

  const handleDelete = async (discussionId: string) => {
    try {
      await api.delete(`/api/v1/discussions/${discussionId}`)
      loadDiscussions()
    } catch (error) {
      console.error('Failed to delete discussion:', error)
      throw error
    }
  }

  const handleResolve = async (discussionId: string) => {
    try {
      const response = await api.post(`/api/v1/discussions/${discussionId}/resolve`)
      
      if (response.success) {
        loadDiscussions()
      }
    } catch (error) {
      console.error('Failed to resolve discussion:', error)
      throw error
    }
  }

  const resetForm = () => {
    setNewTitle('')
    setNewContent('')
    setIsAnonymous(false)
    setSelectedTagIds([])
    setShowMarkdownPreview(false)
  }

  const insertMarkdown = (before: string, after: string = '') => {
    const textarea = document.getElementById('new-discussion-content') as HTMLTextAreaElement
    if (!textarea) return

    const start = textarea.selectionStart
    const end = textarea.selectionEnd
    const text = textarea.value
    const selectedText = text.substring(start, end)
    
    const newText = text.substring(0, start) + before + selectedText + after + text.substring(end)
    setNewContent(newText)
    
    setTimeout(() => {
      textarea.focus()
      textarea.setSelectionRange(start + before.length, end + before.length)
    }, 0)
  }

  const toggleTag = (tagId: string) => {
    setSelectedTags(prev => 
      prev.includes(tagId) 
        ? prev.filter(id => id !== tagId)
        : [...prev, tagId]
    )
  }

  const toggleTagSelection = (tagId: string) => {
    setSelectedTagIds(prev =>
      prev.includes(tagId)
        ? prev.filter(id => id !== tagId)
        : [...prev, tagId]
    )
  }

  const getSortLabel = (sort: string) => {
    switch (sort) {
      case 'created_at':
        return '最新发布'
      case 'updated_at':
        return '最新活动'
      case 'upvotes':
        return '最多点赞'
      case 'reply_count':
        return '最多回复'
      case 'hot':
        return '热门讨论'
      default:
        return '最新发布'
    }
  }

  if (!courseId) {
    return (
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-secondary-800 dark:text-white mb-4">
            课程不存在
          </h1>
          <Button onClick={() => navigate('/courses')}>
            返回课程列表
          </Button>
        </div>
      </div>
    )
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      {/* Header */}
      <div className="mb-8">
        <div className="flex items-center justify-between mb-4">
          <div>
            <h1 className="text-3xl font-bold text-secondary-800 dark:text-white mb-2">
              课程讨论区
            </h1>
            <p className="text-secondary-600 dark:text-secondary-400">
              提问、分享学习心得、参与讨论
            </p>
          </div>
          
          <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
            <DialogTrigger asChild>
              <Button className="flex items-center gap-2">
                <Plus className="h-5 w-5" />
                发起讨论
              </Button>
            </DialogTrigger>
            
            <DialogContent className="max-w-3xl max-h-[90vh] overflow-y-auto">
              <DialogHeader>
                <DialogTitle>发起新讨论</DialogTitle>
              </DialogHeader>
              
              <div className="mt-4 space-y-4">
                <div>
                  <label className="block text-sm font-medium text-secondary-700 dark:text-secondary-300 mb-2">
                    标题 *
                  </label>
                  <Input
                    value={newTitle}
                    onChange={(e) => setNewTitle(e.target.value)}
                    placeholder="简明扼要地描述你的问题或话题"
                    maxLength={200}
                  />
                </div>
                
                <div>
                  <label className="block text-sm font-medium text-secondary-700 dark:text-secondary-300 mb-2">
                    内容 *
                  </label>
                  
                  {/* Markdown Toolbar */}
                  <div className="flex items-center gap-1 p-2 bg-secondary-50 dark:bg-secondary-800 border-b border-secondary-200 dark:border-secondary-700 rounded-t-lg">
                    <Button
                      type="button"
                      variant="ghost"
                      size="sm"
                      onClick={() => insertMarkdown('**', '**')}
                      title="加粗"
                    >
                      <Bold className="h-4 w-4" />
                    </Button>
                    <Button
                      type="button"
                      variant="ghost"
                      size="sm"
                      onClick={() => insertMarkdown('*', '*')}
                      title="斜体"
                    >
                      <Italic className="h-4 w-4" />
                    </Button>
                    <Button
                      type="button"
                      variant="ghost"
                      size="sm"
                      onClick={() => insertMarkdown('- ')}
                      title="列表"
                    >
                      <List className="h-4 w-4" />
                    </Button>
                    <Button
                      type="button"
                      variant="ghost"
                      size="sm"
                      onClick={() => insertMarkdown('[', '](url)')}
                      title="链接"
                    >
                      <LinkIcon className="h-4 w-4" />
                    </Button>
                    <Button
                      type="button"
                      variant="ghost"
                      size="sm"
                      onClick={() => insertMarkdown('@')}
                      title="@提及用户"
                    >
                      <AtSign className="h-4 w-4" />
                    </Button>
                    <Button
                      type="button"
                      variant="ghost"
                      size="sm"
                      onClick={() => insertMarkdown('```\n', '\n```')}
                      title="代码块"
                    >
                      <Code2 className="h-4 w-4" />
                    </Button>
                    <div className="flex-1" />
                    <Button
                      type="button"
                      variant="ghost"
                      size="sm"
                      onClick={() => setShowMarkdownPreview(!showMarkdownPreview)}
                    >
                      {showMarkdownPreview ? '编辑' : '预览'}
                    </Button>
                  </div>
                  
                  {showMarkdownPreview ? (
                    <div className="prose dark:prose-invert max-w-none p-4 border border-secondary-200 dark:border-secondary-700 rounded-b-lg min-h-[300px] bg-white dark:bg-secondary-700">
                      <div dangerouslySetInnerHTML={{ __html: newContent }} />
                    </div>
                  ) : (
                    <Textarea
                      id="new-discussion-content"
                      value={newContent}
                      onChange={(e) => setNewContent(e.target.value)}
                      placeholder="详细描述你的问题或分享... 支持 Markdown 格式"
                      className="min-h-[300px] font-mono"
                    />
                  )}
                </div>
                
                <div>
                  <label className="block text-sm font-medium text-secondary-700 dark:text-secondary-300 mb-2">
                    标签
                  </label>
                  <div className="flex flex-wrap gap-2">
                    {tags.map((tag) => (
                      <Badge
                        key={tag.id}
                        variant={selectedTagIds.includes(tag.id) ? 'primary' : 'secondary'}
                        className="cursor-pointer"
                        onClick={() => toggleTagSelection(tag.id)}
                        style={selectedTagIds.includes(tag.id) ? { backgroundColor: tag.color } : {}}
                      >
                        {tag.name}
                      </Badge>
                    ))}
                  </div>
                </div>
                
                <div className="flex items-center gap-2">
                  <input
                    type="checkbox"
                    id="is-anonymous"
                    checked={isAnonymous}
                    onChange={(e) => setIsAnonymous(e.target.checked)}
                    className="rounded border-secondary-300 dark:border-secondary-600"
                  />
                  <label htmlFor="is-anonymous" className="text-sm text-secondary-600 dark:text-secondary-400">
                    匿名发布
                  </label>
                </div>
                
                <div className="flex justify-end gap-3 pt-4">
                  <Button variant="secondary" onClick={() => setIsCreateDialogOpen(false)}>
                    取消
                  </Button>
                  <Button 
                    onClick={handleCreateDiscussion}
                    disabled={!newTitle.trim() || !newContent.trim() || creating}
                  >
                    {creating ? '发布中...' : '发布讨论'}
                  </Button>
                </div>
              </div>
            </DialogContent>
          </Dialog>
        </div>
        
        {/* Filters and Search */}
        <div className="flex flex-col md:flex-row gap-4 items-start md:items-center justify-between">
          <div className="flex flex-wrap gap-2">
            <Button
              variant={sortBy === 'created_at' ? 'primary' : 'secondary'}
              size="sm"
              onClick={() => setSortBy('created_at')}
            >
              <Clock className="h-4 w-4 mr-1" />
              最新
            </Button>
            <Button
              variant={sortBy === 'hot' ? 'primary' : 'secondary'}
              size="sm"
              onClick={() => setSortBy('hot')}
            >
              <TrendingUp className="h-4 w-4 mr-1" />
              热门
            </Button>
            <Button
              variant={sortBy === 'upvotes' ? 'primary' : 'secondary'}
              size="sm"
              onClick={() => setSortBy('upvotes')}
            >
              <ThumbsUp className="h-4 w-4 mr-1" />
              点赞
            </Button>
            <Button
              variant={sortBy === 'reply_count' ? 'primary' : 'secondary'}
              size="sm"
              onClick={() => setSortBy('reply_count')}
            >
              <MessageCircle className="h-4 w-4 mr-1" />
              回复
            </Button>
          </div>
          
          <div className="flex items-center gap-3 w-full md:w-auto">
            <div className="relative flex-1 md:w-64">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-secondary-400" />
              <Input
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                placeholder="搜索讨论..."
                className="pl-10"
              />
              {searchQuery && (
                <button
                  onClick={() => setSearchQuery('')}
                  className="absolute right-3 top-1/2 transform -translate-y-1/2 text-secondary-400 hover:text-secondary-600"
                >
                  <X className="h-4 w-4" />
                </button>
              )}
            </div>
            
            <Select value={sortBy} onValueChange={setSortBy}>
              <SelectTrigger className="w-full md:w-auto">
                <SelectValue placeholder="排序方式" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="created_at">最新发布</SelectItem>
                <SelectItem value="updated_at">最新活动</SelectItem>
                <SelectItem value="upvotes">最多点赞</SelectItem>
                <SelectItem value="reply_count">最多回复</SelectItem>
                <SelectItem value="hot">热门讨论</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </div>
        
        {/* Tags Filter */}
        {tags.length > 0 && (
          <div className="mt-4 flex flex-wrap gap-2">
            <div className="flex items-center gap-2 text-sm text-secondary-600 dark:text-secondary-400">
              <Filter className="h-4 w-4" />
              <span>标签筛选:</span>
            </div>
            {tags.map((tag) => (
              <Badge
                key={tag.id}
                variant={selectedTags.includes(tag.id) ? 'primary' : 'secondary'}
                className="cursor-pointer"
                onClick={() => toggleTag(tag.id)}
                style={selectedTags.includes(tag.id) ? { backgroundColor: tag.color } : {}}
              >
                {tag.name}
              </Badge>
            ))}
            {selectedTags.length > 0 && (
              <Button
                variant="ghost"
                size="sm"
                onClick={() => setSelectedTags([])}
                className="h-5 px-2 text-xs"
              >
                清除
              </Button>
            )}
          </div>
        )}
      </div>
      
      {/* Stats */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
        <div className="card p-4">
          <div className="flex items-center gap-3">
            <div className="p-2 bg-primary-100 dark:bg-primary-900 rounded-lg">
              <MessageSquare className="h-5 w-5 text-primary-600 dark:text-primary-400" />
            </div>
            <div>
              <p className="text-2xl font-bold text-secondary-800 dark:text-white">{total}</p>
              <p className="text-sm text-secondary-600 dark:text-secondary-400">总讨论数</p>
            </div>
          </div>
        </div>
        
        <div className="card p-4">
          <div className="flex items-center gap-3">
            <div className="p-2 bg-green-100 dark:bg-green-900 rounded-lg">
              <ThumbsUp className="h-5 w-5 text-green-600 dark:text-green-400" />
            </div>
            <div>
              <p className="text-2xl font-bold text-secondary-800 dark:text-white">
                {discussions.reduce((sum, d) => sum + (d.upvotes || 0), 0)}
              </p>
              <p className="text-sm text-secondary-600 dark:text-secondary-400">总点赞数</p>
            </div>
          </div>
        </div>
        
        <div className="card p-4">
          <div className="flex items-center gap-3">
            <div className="p-2 bg-blue-100 dark:bg-blue-900 rounded-lg">
              <MessageCircle className="h-5 w-5 text-blue-600 dark:text-blue-400" />
            </div>
            <div>
              <p className="text-2xl font-bold text-secondary-800 dark:text-white">
                {discussions.reduce((sum, d) => sum + (d.reply_count || 0), 0)}
              </p>
              <p className="text-sm text-secondary-600 dark:text-secondary-400">总回复数</p>
            </div>
          </div>
        </div>
        
        <div className="card p-4">
          <div className="flex items-center gap-3">
            <div className="p-2 bg-purple-100 dark:bg-purple-900 rounded-lg">
              <Eye className="h-5 w-5 text-purple-600 dark:text-purple-400" />
            </div>
            <div>
              <p className="text-2xl font-bold text-secondary-800 dark:text-white">
                {discussions.reduce((sum, d) => sum + (d.view_count || 0), 0)}
              </p>
              <p className="text-sm text-secondary-600 dark:text-secondary-400">总浏览数</p>
            </div>
          </div>
        </div>
      </div>
      
      {/* Discussions List */}
      {loading ? (
        <div className="space-y-4">
          {[...Array(5)].map((_, i) => (
            <div key={i} className="card p-6 animate-pulse">
              <div className="h-6 bg-secondary-200 dark:bg-secondary-700 rounded w-3/4 mb-4" />
              <div className="space-y-2">
                <div className="h-4 bg-secondary-200 dark:bg-secondary-700 rounded" />
                <div className="h-4 bg-secondary-200 dark:bg-secondary-700 rounded w-5/6" />
              </div>
            </div>
          ))}
        </div>
      ) : discussions.length === 0 ? (
        <div className="card p-12 text-center">
          <MessageSquare className="h-16 w-16 text-secondary-300 dark:text-secondary-600 mx-auto mb-4" />
          <h3 className="text-xl font-semibold text-secondary-800 dark:text-white mb-2">
            还没有讨论
          </h3>
          <p className="text-secondary-600 dark:text-secondary-400 mb-6">
            快来发起第一个讨论吧！
          </p>
          <Button onClick={() => setIsCreateDialogOpen(true)}>
            发起讨论
          </Button>
        </div>
      ) : (
        <>
          <div className="space-y-4">
            {discussions.map((discussion) => (
              <DiscussionThread
                key={discussion.id}
                discussion={discussion}
                onReply={handleReply}
                onLike={handleLike}
                onFavorite={handleFavorite}
                onEdit={handleEdit}
                onDelete={handleDelete}
                onResolve={handleResolve}
                currentUserId={user?.id}
              />
            ))}
          </div>
          
          {/* Pagination */}
          {totalPages > 1 && (
            <div className="flex justify-center items-center gap-2 mt-8">
              <Button
                variant="secondary"
                onClick={() => setPage(p => Math.max(1, p - 1))}
                disabled={page === 1}
              >
                上一页
              </Button>
              
              <span className="text-sm text-secondary-600 dark:text-secondary-400">
                第 {page} 页，共 {totalPages} 页
              </span>
              
              <Button
                variant="secondary"
                onClick={() => setPage(p => Math.min(totalPages, p + 1))}
                disabled={page === totalPages}
              >
                下一页
              </Button>
            </div>
          )}
        </>
      )}
    </div>
  )
}

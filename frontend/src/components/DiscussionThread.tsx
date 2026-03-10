import { useState } from 'react'
import { 
  MessageSquare, 
  ThumbsUp, 
  ThumbsDown, 
  Star, 
  Share2, 
  Flag, 
  MoreVertical,
  ChevronDown,
  ChevronUp,
  CheckCircle,
  Lock,
  Pin,
  Edit2,
  Trash2,
  Reply,
  Send,
  AtSign,
  Code2,
  Bold,
  Italic,
  List,
  Link as LinkIcon,
} from 'lucide-react'
import { formatDistanceToNow } from 'date-fns'
import { zhCN } from 'date-fns/locale'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import rehypeHighlight from 'rehype-highlight'
import 'highlight.js/styles/github.css'
import { Avatar } from '@/components/ui/Avatar'
import { Button } from '@/components/ui/Button'
import { Badge } from '@/components/ui/Badge'
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from '@/components/ui/DropdownMenu'
import { type Discussion } from '@/types/discussion'

interface DiscussionThreadProps {
  discussion: Discussion
  depth?: number
  onReply?: (parentId: string, content: string) => Promise<void>
  onLike?: (discussionId: string, type: 'upvote' | 'downvote') => Promise<void>
  onFavorite?: (discussionId: string) => Promise<void>
  onEdit?: (discussionId: string, content: string, title?: string) => Promise<void>
  onDelete?: (discussionId: string) => Promise<void>
  onResolve?: (discussionId: string) => Promise<void>
  currentUserId?: string
  maxDepth?: number
}

export default function DiscussionThread({
  discussion,
  depth = 0,
  onReply,
  onLike,
  onFavorite,
  onEdit,
  onDelete,
  onResolve,
  currentUserId,
  maxDepth = 10,
}: DiscussionThreadProps) {
  const [isReplying, setIsReplying] = useState(false)
  const [replyContent, setReplyContent] = useState('')
  const [isExpanded, setIsExpanded] = useState(true)
  const [isEditing, setIsEditing] = useState(false)
  const [editContent, setEditContent] = useState(discussion.content)
  const [editTitle, setEditTitle] = useState(discussion.title || '')
  const [showMarkdownPreview, setShowMarkdownPreview] = useState(false)

  const isAuthor = currentUserId === discussion.user.id
  const canReply = depth < maxDepth && !discussion.is_locked
  const canEdit = isAuthor && !isEditing
  const canDelete = isAuthor
  const canResolve = isAuthor && discussion.parent_id === null

  const handleReply = async () => {
    if (!replyContent.trim() || !onReply) return
    
    try {
      await onReply(discussion.id, replyContent)
      setReplyContent('')
      setIsReplying(false)
    } catch (error) {
      console.error('Failed to post reply:', error)
    }
  }

  const handleLike = async (type: 'upvote' | 'downvote') => {
    if (!onLike) return
    try {
      await onLike(discussion.id, type)
    } catch (error) {
      console.error('Failed to toggle like:', error)
    }
  }

  const handleFavorite = async () => {
    if (!onFavorite) return
    try {
      await onFavorite(discussion.id)
    } catch (error) {
      console.error('Failed to toggle favorite:', error)
    }
  }

  const handleEdit = async () => {
    if (!onEdit || !editContent.trim()) return
    try {
      await onEdit(discussion.id, editContent, discussion.parent_id ? undefined : editTitle)
      setIsEditing(false)
    } catch (error) {
      console.error('Failed to update discussion:', error)
    }
  }

  const handleDelete = async () => {
    if (!onDelete || !confirm('确定要删除这条讨论吗？此操作不可恢复。')) return
    try {
      await onDelete(discussion.id)
    } catch (error) {
      console.error('Failed to delete discussion:', error)
    }
  }

  const handleResolve = async () => {
    if (!onResolve || !confirm('标记此讨论为已解决？')) return
    try {
      await onResolve(discussion.id)
    } catch (error) {
      console.error('Failed to resolve discussion:', error)
    }
  }

  const insertMarkdown = (before: string, after: string = '') => {
    const textarea = document.getElementById(`discussion-edit-${discussion.id}`) as HTMLTextAreaElement
    if (!textarea) return

    const start = textarea.selectionStart
    const end = textarea.selectionEnd
    const text = textarea.value
    const selectedText = text.substring(start, end)
    
    const newText = text.substring(0, start) + before + selectedText + after + text.substring(end)
    setEditContent(newText)
    
    setTimeout(() => {
      textarea.focus()
      textarea.setSelectionRange(start + before.length, end + before.length)
    }, 0)
  }

  const renderMarkdownToolbar = () => (
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
  )

  const renderActions = () => (
    <div className="flex items-center gap-2 mt-3">
      <Button
        variant="ghost"
        size="sm"
        onClick={() => handleLike('upvote')}
        className={`flex items-center gap-1 ${
          discussion.is_liked ? 'text-green-600 dark:text-green-400' : ''
        }`}
      >
        <ThumbsUp className="h-4 w-4" />
        <span className="text-sm">{discussion.upvotes || 0}</span>
      </Button>
      
      <Button
        variant="ghost"
        size="sm"
        onClick={() => handleLike('downvote')}
        className={`flex items-center gap-1 ${
          discussion.is_liked === false ? 'text-red-600 dark:text-red-400' : ''
        }`}
      >
        <ThumbsDown className="h-4 w-4" />
        <span className="text-sm">{discussion.downvotes || 0}</span>
      </Button>

      <Button
        variant="ghost"
        size="sm"
        onClick={handleFavorite}
        className={`flex items-center gap-1 ${
          discussion.is_favorited ? 'text-yellow-600 dark:text-yellow-400' : ''
        }`}
      >
        <Star className="h-4 w-4" />
        <span className="text-sm">收藏</span>
      </Button>

      {canReply && (
        <Button
          variant="ghost"
          size="sm"
          onClick={() => setIsReplying(!isReplying)}
          className="flex items-center gap-1"
        >
          <Reply className="h-4 w-4" />
          <span className="text-sm">回复</span>
        </Button>
      )}

      {discussion.replies && discussion.replies.length > 0 && (
        <Button
          variant="ghost"
          size="sm"
          onClick={() => setIsExpanded(!isExpanded)}
          className="flex items-center gap-1"
        >
          {isExpanded ? (
            <ChevronUp className="h-4 w-4" />
          ) : (
            <ChevronDown className="h-4 w-4" />
          )}
          <span className="text-sm">{discussion.replies.length} 条回复</span>
        </Button>
      )}

      <div className="flex-1" />

      {(canEdit || canDelete || canResolve) && (
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" size="sm">
              <MoreVertical className="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            {canResolve && (
              <DropdownMenuItem onClick={handleResolve}>
                <CheckCircle className="h-4 w-4 mr-2" />
                标记为已解决
              </DropdownMenuItem>
            )}
            {canEdit && (
              <DropdownMenuItem onClick={() => setIsEditing(true)}>
                <Edit2 className="h-4 w-4 mr-2" />
                编辑
              </DropdownMenuItem>
            )}
            {canDelete && (
              <DropdownMenuItem onClick={handleDelete} className="text-red-600">
                <Trash2 className="h-4 w-4 mr-2" />
                删除
              </DropdownMenuItem>
            )}
            <DropdownMenuItem>
              <Flag className="h-4 w-4 mr-2" />
              举报
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      )}
    </div>
  )

  return (
    <div className={`${depth > 0 ? 'ml-6 md:ml-12 border-l-2 border-secondary-200 dark:border-secondary-700 pl-4 md:pl-6' : ''}`}>
      <div className="bg-white dark:bg-secondary-800 rounded-lg p-4 md:p-6 my-3 shadow-sm">
        {/* Header */}
        <div className="flex items-start gap-3 mb-3">
          <Avatar 
            src={discussion.user?.avatar_url} 
            name={discussion.user?.username || '匿名用户'}
            size="md"
          />
          
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2 flex-wrap">
              <span className="font-medium text-secondary-800 dark:text-white">
                {discussion.is_anonymous ? '匿名用户' : discussion.user?.username || '未知用户'}
              </span>
              
              {discussion.is_pinned && (
                <Badge variant="primary" size="sm">
                  <Pin className="h-3 w-3 mr-1" />
                  置顶
                </Badge>
              )}
              
              {discussion.is_resolved && (
                <Badge variant="success" size="sm">
                  <CheckCircle className="h-3 w-3 mr-1" />
                  已解决
                </Badge>
              )}
              
              {discussion.is_locked && (
                <Badge variant="secondary" size="sm">
                  <Lock className="h-3 w-3 mr-1" />
                  已锁定
                </Badge>
              )}
              
              {discussion.tags?.map((tag) => (
                <Badge 
                  key={tag.id} 
                  size="sm"
                  style={{ backgroundColor: tag.color, color: 'white' }}
                >
                  {tag.name}
                </Badge>
              ))}
              
              <span className="text-xs text-secondary-500 dark:text-secondary-400">
                {formatDistanceToNow(new Date(discussion.created_at), { 
                  addSuffix: true, 
                  locale: zhCN 
                })}
              </span>
            </div>
            
            {discussion.title && (
              <h3 className="text-lg font-semibold text-secondary-800 dark:text-white mt-2">
                {discussion.title}
              </h3>
            )}
          </div>
        </div>

        {/* Content */}
        {isEditing ? (
          <div className="mt-3">
            {discussion.parent_id === null && (
              <input
                type="text"
                value={editTitle}
                onChange={(e) => setEditTitle(e.target.value)}
                className="w-full px-3 py-2 mb-2 border border-secondary-300 dark:border-secondary-600 rounded-lg bg-white dark:bg-secondary-700 text-secondary-800 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                placeholder="标题"
              />
            )}
            
            {renderMarkdownToolbar()}
            
            {showMarkdownPreview ? (
              <div className="prose dark:prose-invert max-w-none p-4 border border-secondary-200 dark:border-secondary-700 rounded-b-lg min-h-[200px]">
                <ReactMarkdown 
                  remarkPlugins={[remarkGfm]}
                  rehypePlugins={[rehypeHighlight]}
                  components={{
                    code({ node, className, children, ...props }: any) {
                      return (
                        <code className={`${className} bg-secondary-100 dark:bg-secondary-800 rounded px-1 py-0.5`} {...props}>
                          {children}
                        </code>
                      )
                    },
                  }}
                >
                  {editContent}
                </ReactMarkdown>
              </div>
            ) : (
              <textarea
                id={`discussion-edit-${discussion.id}`}
                value={editContent}
                onChange={(e) => setEditContent(e.target.value)}
                className="w-full px-3 py-2 border border-secondary-300 dark:border-secondary-600 rounded-b-lg bg-white dark:bg-secondary-700 text-secondary-800 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent min-h-[200px] font-mono text-sm"
                placeholder="编写你的内容... 支持 Markdown 格式"
              />
            )}
            
            <div className="flex justify-end gap-2 mt-3">
              <Button variant="secondary" onClick={() => setIsEditing(false)}>
                取消
              </Button>
              <Button onClick={handleEdit}>
                保存
              </Button>
            </div>
          </div>
        ) : (
          <div className="prose dark:prose-invert max-w-none mt-3">
            <ReactMarkdown 
              remarkPlugins={[remarkGfm]}
              rehypePlugins={[rehypeHighlight]}
              components={{
                code({ node, className, children, ...props }: any) {
                  const match = /language-(\w+)/.exec(className || '')
                  return match ? (
                    <pre className="bg-secondary-100 dark:bg-secondary-800 rounded-lg p-4 overflow-x-auto">
                      <code className={className} {...props}>
                        {children}
                      </code>
                    </pre>
                  ) : (
                    <code className={`${className} bg-secondary-100 dark:bg-secondary-800 rounded px-1 py-0.5`} {...props}>
                      {children}
                    </code>
                  )
                },
              }}
            >
              {discussion.content}
            </ReactMarkdown>
          </div>
        )}

        {/* Actions */}
        {!isEditing && renderActions()}

        {/* Reply Form */}
        {isReplying && canReply && (
          <div className="mt-4 p-4 bg-secondary-50 dark:bg-secondary-900 rounded-lg">
            <div className="flex items-center gap-2 mb-2">
              <AtSign className="h-4 w-4 text-secondary-500" />
              <span className="text-sm text-secondary-600 dark:text-secondary-400">
                回复 @{discussion.user?.username || '匿名用户'}
              </span>
            </div>
            
            {renderMarkdownToolbar()}
            
            <textarea
              value={replyContent}
              onChange={(e) => setReplyContent(e.target.value)}
              className="w-full px-3 py-2 border border-secondary-300 dark:border-secondary-600 rounded-b-lg bg-white dark:bg-secondary-700 text-secondary-800 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent min-h-[100px] font-mono text-sm"
              placeholder="编写你的回复... 支持 Markdown 格式"
            />
            
            <div className="flex justify-end gap-2 mt-3">
              <Button variant="secondary" onClick={() => setIsReplying(false)}>
                取消
              </Button>
              <Button onClick={handleReply} disabled={!replyContent.trim()}>
                <Send className="h-4 w-4 mr-2" />
                发布回复
              </Button>
            </div>
          </div>
        )}
      </div>

      {/* Nested Replies */}
      {isExpanded && discussion.replies && discussion.replies.length > 0 && (
        <div className="mt-2">
          {discussion.replies.map((reply) => (
            <DiscussionThread
              key={reply.id}
              discussion={reply}
              depth={depth + 1}
              onReply={onReply}
              onLike={onLike}
              onFavorite={onFavorite}
              onEdit={onEdit}
              onDelete={onDelete}
              onResolve={onResolve}
              currentUserId={currentUserId}
              maxDepth={maxDepth}
            />
          ))}
        </div>
      )}
    </div>
  )
}

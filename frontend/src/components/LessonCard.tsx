import { CheckCircle, Circle, Lock, PlayCircle, Code, HelpCircle, Folder } from 'lucide-react'

export interface Lesson {
  id: string
  title: string
  duration: string
  completed: boolean
  locked: boolean
  type: 'video' | 'exercise' | 'quiz' | 'project'
}

interface LessonCardProps {
  lesson: Lesson
  onClick?: () => void
}

export default function LessonCard({ lesson, onClick }: LessonCardProps) {
  const typeIcons = {
    video: PlayCircle,
    exercise: Code,
    quiz: HelpCircle,
    project: Folder,
  }

  const Icon = lesson.completed
    ? CheckCircle
    : lesson.locked
    ? Lock
    : typeIcons[lesson.type] || PlayCircle

  return (
    <button
      onClick={onClick}
      disabled={lesson.locked}
      className={`w-full flex items-center p-4 rounded-lg border transition-all duration-200 ${
        lesson.completed
          ? 'bg-green-50 border-green-200 dark:bg-green-900/20 dark:border-green-800'
          : lesson.locked
          ? 'bg-secondary-50 border-secondary-200 dark:bg-secondary-800 dark:border-secondary-700 cursor-not-allowed'
          : 'bg-white border-secondary-200 dark:bg-secondary-800 dark:border-secondary-700 hover:border-primary-300 dark:hover:border-primary-700 cursor-pointer'
      }`}
    >
      <Icon
        className={`h-6 w-6 mr-4 flex-shrink-0 ${
          lesson.completed
            ? 'text-green-600'
            : lesson.locked
            ? 'text-secondary-400'
            : 'text-primary-600'
        }`}
      />
      
      <div className="flex-1 text-left">
        <h4 className={`font-medium ${
          lesson.locked
            ? 'text-secondary-400'
            : 'text-secondary-800 dark:text-white'
        }`}>
          {lesson.title}
        </h4>
        <p className="text-sm text-secondary-500 dark:text-secondary-400 mt-1">
          {lesson.duration}
        </p>
      </div>

      {!lesson.locked && !lesson.completed && (
        <span className="text-primary-600 text-sm font-medium">
          开始学习
        </span>
      )}
    </button>
  )
}

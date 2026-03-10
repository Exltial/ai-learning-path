import { Link } from 'react-router-dom'
import { BookOpen, Clock, Star } from 'lucide-react'

export interface Course {
  id: string
  title: string
  description: string
  thumbnail: string
  level: 'beginner' | 'intermediate' | 'advanced'
  duration: string
  students: number
  rating: number
  lessons: number
  tags: string[]
}

interface CourseCardProps {
  course: Course
}

export default function CourseCard({ course }: CourseCardProps) {
  const levelLabels = {
    beginner: '入门',
    intermediate: '进阶',
    advanced: '高级',
  }

  const levelColors = {
    beginner: 'bg-green-100 text-green-800',
    intermediate: 'bg-yellow-100 text-yellow-800',
    advanced: 'bg-red-100 text-red-800',
  }

  return (
    <Link to={`/courses/${course.id}`} className="block group">
      <div className="card h-full flex flex-col">
        {/* Thumbnail */}
        <div className="relative h-48 mb-4 rounded-lg overflow-hidden bg-secondary-200">
          <img
            src={course.thumbnail}
            alt={course.title}
            className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
          />
          <div className="absolute top-2 right-2">
            <span className={`px-2 py-1 rounded-full text-xs font-medium ${levelColors[course.level]}`}>
              {levelLabels[course.level]}
            </span>
          </div>
        </div>

        {/* Content */}
        <h3 className="text-lg font-semibold text-secondary-800 dark:text-white mb-2 group-hover:text-primary-600 transition-colors">
          {course.title}
        </h3>
        
        <p className="text-secondary-600 dark:text-secondary-400 text-sm mb-4 flex-1">
          {course.description}
        </p>

        {/* Tags */}
        <div className="flex flex-wrap gap-2 mb-4">
          {course.tags.slice(0, 3).map((tag) => (
            <span
              key={tag}
              className="px-2 py-1 bg-primary-50 dark:bg-primary-900 text-primary-700 dark:text-primary-300 text-xs rounded-full"
            >
              {tag}
            </span>
          ))}
        </div>

        {/* Meta Info */}
        <div className="flex items-center justify-between text-sm text-secondary-500 dark:text-secondary-400 pt-4 border-t border-secondary-200 dark:border-secondary-700">
          <div className="flex items-center space-x-4">
            <span className="flex items-center">
              <Clock className="h-4 w-4 mr-1" />
              {course.duration}
            </span>
            <span className="flex items-center">
              <BookOpen className="h-4 w-4 mr-1" />
              {course.lessons} 课
            </span>
          </div>
          <div className="flex items-center space-x-1">
            <Star className="h-4 w-4 fill-yellow-400 text-yellow-400" />
            <span className="font-medium">{course.rating}</span>
            <span className="text-secondary-400">({course.students})</span>
          </div>
        </div>
      </div>
    </Link>
  )
}

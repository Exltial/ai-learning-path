import { Outlet } from 'react-router-dom'
import Header from './Header'
import Footer from './Footer'
import MobileNav from './MobileNav'
import PWAInstallPrompt from './PWAInstallPrompt'

export default function Layout() {
  return (
    <div className="min-h-screen flex flex-col bg-secondary-50 dark:bg-secondary-900">
      <PWAInstallPrompt />
      <Header />
      <main className="flex-1 pb-20 md:pb-0">
        <Outlet />
      </main>
      <Footer className="hidden md:block" />
      <MobileNav />
    </div>
  )
}

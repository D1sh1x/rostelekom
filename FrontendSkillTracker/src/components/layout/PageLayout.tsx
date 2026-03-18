import { Outlet } from 'react-router-dom'
import AppSidebar from './AppSidebar'

export default function PageLayout() {
  return (
    <div className="flex h-screen overflow-hidden bg-background">
      <AppSidebar />
      <main className="flex-1 overflow-y-auto scrollbar-thin">
        <div className="min-h-full p-6 lg:p-8">
          <Outlet />
        </div>
      </main>
    </div>
  )
}

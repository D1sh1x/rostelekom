import React from 'react';
import { NavLink, useNavigate } from 'react-router-dom';
import { 
  LayoutDashboard, 
  CheckSquare, 
  Users, 
  LogOut, 
  Clock,
  PlusCircle
} from 'lucide-react';
import { useAuth } from '../context/AuthContext';

const Sidebar: React.FC = () => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = async () => {
    await logout();
    navigate('/login');
  };

  return (
    <aside className="w-64 h-screen bg-surface border-r border-border flex flex-col fixed left-0 top-0">
      <div className="p-6 border-b border-border">
        <h1 className="text-xl font-bold flex items-center gap-2">
          <div className="w-8 h-8 bg-primary rounded-lg flex items-center justify-center">
            <CheckSquare size={20} className="text-white" />
          </div>
          SkillTracker
        </h1>
      </div>

      <nav className="flex-1 p-4 flex flex-col gap-2">
        <NavLink 
          to="/" 
          className={({ isActive }) => `flex items-center gap-3 px-4 py-3 rounded-lg transition-colors ${isActive ? 'bg-primary-bg text-primary' : 'hover:bg-bg-accent text-text-secondary'}`}
        >
          <LayoutDashboard size={20} />
          Dashboard
        </NavLink>
        
        <NavLink 
          to="/tasks" 
          className={({ isActive }) => `flex items-center gap-3 px-4 py-3 rounded-lg transition-colors ${isActive ? 'bg-primary-bg text-primary' : 'hover:bg-bg-accent text-text-secondary'}`}
        >
          <Clock size={20} />
          All Tasks
        </NavLink>

        {user?.role === 'manager' && (
          <>
            <NavLink 
              to="/users" 
              className={({ isActive }) => `flex items-center gap-3 px-4 py-3 rounded-lg transition-colors ${isActive ? 'bg-primary-bg text-primary' : 'hover:bg-bg-accent text-text-secondary'}`}
            >
              <Users size={20} />
              Employees
            </NavLink>
            <NavLink 
              to="/tasks/new" 
              className={({ isActive }) => `flex items-center gap-3 px-4 py-3 rounded-lg transition-colors ${isActive ? 'bg-primary-bg text-primary' : 'hover:bg-bg-accent text-text-secondary'}`}
            >
              <PlusCircle size={20} />
              Create Task
            </NavLink>
          </>
        )}
      </nav>

      <div className="p-4 border-t border-border mt-auto">
        <div className="flex items-center gap-3 px-4 py-3 mb-4">
          <div className="w-10 h-10 rounded-full bg-accent flex items-center justify-center font-bold text-primary">
            {user?.username.substring(0, 2).toUpperCase()}
          </div>
          <div className="flex flex-col overflow-hidden">
            <span className="font-medium truncate">{user?.username}</span>
            <span className="text-xs text-text-muted capitalize">{user?.role}</span>
          </div>
        </div>
        
        <button 
          onClick={handleLogout}
          className="w-full flex items-center gap-3 px-4 py-3 text-text-secondary hover:text-error hover:bg-error/10 rounded-lg transition-colors"
        >
          <LogOut size={20} />
          Logout
        </button>
      </div>
    </aside>
  );
};

export default Sidebar;

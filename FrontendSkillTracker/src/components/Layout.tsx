import React from 'react';
import { useAuth } from '../contexts/AuthContext';
import { 
  LayoutDashboard, 
  ClipboardList, 
  Users, 
  LogOut, 
  Bell, 
  Search,
  ChevronRight,
  Shield,
} from 'lucide-react';
import { Link, useLocation } from 'react-router-dom';
import { motion } from 'framer-motion';

const SidebarItem: React.FC<{ 
  to: string, 
  icon: React.ReactNode, 
  label: string, 
  active?: boolean 
}> = ({ to, icon, label, active }) => (
  <Link to={to}>
    <motion.div 
      whileHover={{ x: 4 }}
      className={`flex items-center gap-3 px-4 py-3 rounded-xl transition-all duration-200 group ${
        active 
          ? 'bg-blue-600/20 text-blue-400 border border-blue-500/20 shadow-[0_0_15px_rgba(59,130,246,0.15)]' 
          : 'text-slate-400 hover:text-white hover:bg-white/5'
      }`}
    >
      <div className={`${active ? 'text-blue-400' : 'text-slate-500 group-hover:text-slate-300'}`}>
        {icon}
      </div>
      <span className="font-medium tracking-wide text-sm">{label}</span>
      {active && <ChevronRight className="ml-auto w-4 h-4" />}
    </motion.div>
  </Link>
);

const Layout: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { user, logout } = useAuth();
  const location = useLocation();

  const navItems = [
    { to: '/', icon: <LayoutDashboard className="w-5 h-5" />, label: 'DASHBOARD' },
    { to: '/quests', icon: <ClipboardList className="w-5 h-5" />, label: 'QUEST LOG' },
  ];

  if (user?.role === 'manager') {
    navItems.push({ to: '/personnel', icon: <Users className="w-5 h-5" />, label: 'PERSONNEL' });
  }

  return (
    <div className="flex h-screen bg-background overflow-hidden">
      {/* Sidebar */}
      <aside className="w-72 glass-panel m-4 rounded-3xl flex flex-col border-white/5">
        <div className="p-8 flex items-center gap-3">
          <div className="w-10 h-10 bg-blue-600/20 rounded-lg flex items-center justify-center border border-blue-500/30">
            <Shield className="w-6 h-6 text-blue-400" />
          </div>
          <div>
            <h1 className="text-xl font-bold tracking-tighter text-white">SKILL<span className="text-blue-500">NEXUS</span></h1>
            <p className="text-[10px] text-slate-500 tracking-[0.2em] font-bold">OPERATIONS CENTER</p>
          </div>
        </div>

        <nav className="flex-1 px-4 space-y-2">
          <p className="px-4 text-[10px] font-bold text-slate-600 tracking-widest uppercase mb-4">Navigation</p>
          {navItems.map((item) => (
            <SidebarItem 
              key={item.to} 
              {...item} 
              active={location.pathname === item.to} 
            />
          ))}
        </nav>

        <div className="p-4 mt-auto">
          <div className="p-4 bg-slate-900/40 rounded-2xl border border-white/5 mb-4">
            <div className="flex items-center gap-3 mb-3">
              <div className="w-8 h-8 rounded-full bg-blue-500/20 flex items-center justify-center border border-blue-500/30 text-xs font-bold text-blue-400">
                {user?.name?.[0].toUpperCase()}
              </div>
              <div className="overflow-hidden">
                <p className="text-sm font-bold text-white truncate">{user?.name}</p>
                <p className="text-[10px] text-blue-500 font-bold uppercase tracking-tighter">{user?.role}</p>
              </div>
            </div>
            <button 
              onClick={logout}
              className="w-full flex items-center justify-center gap-2 py-2 rounded-lg text-xs font-bold text-red-400 hover:bg-red-500/10 transition-all border border-red-500/20"
            >
              <LogOut className="w-4 h-4" /> TERMINATE SESSION
            </button>
          </div>
        </div>
      </aside>

      {/* Main Content Area */}
      <div className="flex-1 flex flex-col overflow-hidden relative">
        {/* Top Navbar */}
        <header className="h-20 flex items-center justify-between px-8 z-10">
          <div className="flex items-center gap-4 bg-slate-900/50 px-4 py-2 rounded-xl border border-white/5 w-96">
            <Search className="w-4 h-4 text-slate-500" />
            <input 
              type="text" 
              placeholder="Search across Nexus intelligence..." 
              className="bg-transparent border-none outline-none text-sm text-slate-300 w-full placeholder:text-slate-600"
            />
          </div>

          <div className="flex items-center gap-6">
            <div className="flex items-center gap-2 px-3 py-1.5 bg-green-500/5 border border-green-500/20 rounded-full">
              <div className="w-2 h-2 bg-green-500 rounded-full animate-pulse" />
              <span className="text-[10px] font-bold text-green-500 tracking-widest">UPLINK STABLE</span>
            </div>
            
            <button className="relative p-2 text-slate-400 hover:text-white transition-colors group">
              <Bell className="w-5 h-5" />
              <span className="absolute top-2 right-2 w-2 h-2 bg-blue-500 rounded-full border-2 border-background" />
              <div className="absolute -bottom-8 left-1/2 -translate-x-1/2 bg-slate-800 text-[10px] text-white px-2 py-1 rounded opacity-0 group-hover:opacity-100 transition-opacity whitespace-nowrap border border-white/10 pointer-events-none uppercase tracking-tighter font-bold">
                Intelligence Alerts
              </div>
            </button>
          </div>
        </header>

        {/* Content */}
        <main className="flex-1 overflow-y-auto p-8 relative custom-scrollbar">
          {children}
        </main>
      </div>
    </div>
  );
};

export default Layout;

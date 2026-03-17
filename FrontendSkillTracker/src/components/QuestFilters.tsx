import React from 'react';
import { Search, Filter } from 'lucide-react';
import type { TaskStatus } from '../types';

interface QuestFiltersProps {
  onSearch: (query: string) => void;
  onStatusChange: (status: TaskStatus | 'all') => void;
  activeStatus: TaskStatus | 'all';
}

const QuestFilters: React.FC<QuestFiltersProps> = ({ onSearch, onStatusChange, activeStatus }) => {
  const statuses: (TaskStatus | 'all')[] = ['all', 'pending', 'in_progress', 'completed'];

  return (
    <div className="flex flex-col md:flex-row gap-4 mb-8">
      <div className="flex-1 glass-panel flex items-center gap-3 px-4 py-2 rounded-xl border-white/5 focus-within:border-blue-500/30 transition-all">
        <Search className="w-4 h-4 text-slate-500" />
        <input 
          type="text" 
          placeholder="Filter Intelligence..." 
          className="bg-transparent border-none outline-none text-sm text-slate-300 w-full placeholder:text-slate-600"
          onChange={(e) => onSearch(e.target.value)}
        />
      </div>

      <div className="flex items-center gap-2 bg-slate-900/50 p-1 rounded-xl border border-white/5">
        {statuses.map((status) => (
          <button
            key={status}
            onClick={() => onStatusChange(status)}
            className={`px-4 py-2 rounded-lg text-[10px] font-bold uppercase tracking-widest transition-all ${
              activeStatus === status 
                ? 'bg-blue-600 text-white shadow-[0_0_15px_rgba(59,130,246,0.5)]' 
                : 'text-slate-500 hover:text-slate-300 hover:bg-white/5'
            }`}
          >
            {status.replace('_', ' ')}
          </button>
        ))}
      </div>
      
      <button className="p-2 glass-panel rounded-xl border-white/5 text-slate-500 hover:text-white transition-all">
        <Filter className="w-5 h-5" />
      </button>
    </div>
  );
};

export default QuestFilters;

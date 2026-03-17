import React from 'react';
import type { Task } from '../types';
import { motion } from 'framer-motion';
import { Calendar, User, ArrowUpRight } from 'lucide-react';
import StatusBadge from './StatusBadge';

const QuestCard: React.FC<{ task: Task; onClick?: () => void }> = ({ task, onClick }) => {
  return (
    <motion.div
      whileHover={{ y: -4 }}
      onClick={onClick}
      className="glass-panel p-5 rounded-2xl border-white/5 group cursor-pointer transition-all duration-300 hover:border-white/10 hover:shadow-[0_20px_40px_rgba(0,0,0,0.4)]"
    >
      <div className="flex justify-between items-start mb-4">
        <StatusBadge status={task.status} />
        <button className="p-2 rounded-lg bg-white/5 text-slate-500 group-hover:text-blue-400 group-hover:bg-blue-500/10 transition-all">
          <ArrowUpRight className="w-4 h-4" />
        </button>
      </div>

      <h3 className="text-lg font-bold text-white mb-2 group-hover:text-glow transition-all">{task.title}</h3>
      <p className="text-xs text-slate-400 line-clamp-2 mb-6 leading-relaxed uppercase font-medium tracking-tight">
        {task.description}
      </p>

      {/* Progress Bar */}
      <div className="mb-6 space-y-2">
        <div className="flex justify-between text-[10px] font-bold tracking-widest text-slate-500">
          <span>SYNC PROGRESS</span>
          <span className="text-blue-400">{task.progress}%</span>
        </div>
        <div className="h-1 bg-slate-800 rounded-full overflow-hidden">
          <motion.div 
            initial={{ width: 0 }}
            animate={{ width: `${task.progress}%` }}
            className="h-full bg-gradient-to-r from-blue-600 to-blue-400"
          />
        </div>
      </div>

      <div className="flex items-center justify-between pt-4 border-t border-white/5">
        <div className="flex items-center gap-2">
          <div className="w-6 h-6 rounded-full bg-slate-800 border border-white/5 flex items-center justify-center">
            <User className="w-3 h-3 text-slate-500" />
          </div>
          <span className="text-[10px] font-bold text-slate-400 uppercase tracking-tighter">
            {task.employee?.name || `Agent #${task.employee_id}`}
          </span>
        </div>
        <div className="flex items-center gap-1.5 text-slate-500">
          <Calendar className="w-3 h-3" />
          <span className="text-[10px] font-mono">
            {new Date(task.deadline).toLocaleDateString('en-GB')}
          </span>
        </div>
      </div>
    </motion.div>
  );
};

export default QuestCard;

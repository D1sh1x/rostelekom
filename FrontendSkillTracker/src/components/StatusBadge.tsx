import React from 'react';
import type { TaskStatus } from '../types';
import { Zap, Clock, CheckCircle2 } from 'lucide-react';

const StatusBadge: React.FC<{ status: TaskStatus }> = ({ status }) => {
  const configs = {
    pending: {
      color: 'bg-yellow-500/10 text-yellow-500 border-yellow-500/20',
      icon: <Clock className="w-3 h-3" />,
      label: 'PENDING'
    },
    in_progress: {
      color: 'bg-blue-500/10 text-blue-400 border-blue-500/20 shadow-[0_0_10px_rgba(59,130,246,0.1)]',
      icon: <Zap className="w-3 h-3 animate-pulse" />,
      label: 'NEURAL SYNC'
    },
    completed: {
      color: 'bg-green-500/10 text-green-500 border-green-500/20',
      icon: <CheckCircle2 className="w-3 h-3" />,
      label: 'FINALIZED'
    }
  };

  const { color, icon, label } = configs[status];

  return (
    <div className={`px-2 py-1 rounded-md border ${color} flex items-center gap-1.5 text-[10px] font-bold tracking-widest`}>
      {icon}
      {label}
    </div>
  );
};

export default StatusBadge;

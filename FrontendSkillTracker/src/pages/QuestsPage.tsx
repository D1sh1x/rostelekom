import React, { useState, useEffect } from 'react';
import { taskApi } from '../api/tasks';
import type { Task, TaskStatus } from '../types';
import QuestCard from '../components/QuestCard';
import QuestFilters from '../components/QuestFilters';
import { motion, AnimatePresence } from 'framer-motion';
import { ClipboardList, Plus, Loader2 } from 'lucide-react';
import { useAuth } from '../contexts/AuthContext';

const QuestsPage: React.FC = () => {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState('');
  const [statusFilter, setStatusFilter] = useState<TaskStatus | 'all'>('all');
  const { user } = useAuth();

  useEffect(() => {
    fetchQuests();
  }, [statusFilter]);

  const fetchQuests = async () => {
    setLoading(true);
    try {
      const filters = statusFilter === 'all' ? {} : { status: statusFilter };
      const data = await taskApi.getTasks(filters);
      setTasks(data || []);
    } catch (error) {
      console.error('Failed to fetch intelligence:', error);
    } finally {
      setLoading(false);
    }
  };

  const filteredTasks = tasks.filter(task => 
    task.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
    task.description.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <div className="space-y-8 pb-10">
      <div className="flex items-end justify-between">
        <div>
          <h2 className="text-sm font-bold text-blue-500 uppercase tracking-[0.3em] mb-1">Intelligence Database</h2>
          <h1 className="text-4xl font-bold text-white tracking-tight">Active <span className="text-glow">Quest Log</span></h1>
        </div>
        
        {user?.role === 'manager' && (
          <motion.button 
            whileHover={{ scale: 1.02 }}
            whileTap={{ scale: 0.98 }}
            className="flex items-center gap-2 px-6 py-3 bg-blue-600 hover:bg-blue-500 text-white rounded-xl font-bold text-xs tracking-widest shadow-[0_0_20px_rgba(59,130,246,0.3)] transition-all"
          >
            <Plus className="w-4 h-4" /> INITIALIZE QUEST
          </motion.button>
        )}
      </div>

      <QuestFilters 
        onSearch={setSearchQuery} 
        onStatusChange={setStatusFilter} 
        activeStatus={statusFilter} 
      />

      {loading ? (
        <div className="flex flex-col items-center justify-center py-20 text-slate-500 gap-4">
          <Loader2 className="w-10 h-10 animate-spin text-blue-500" />
          <p className="text-[10px] font-bold tracking-[0.3em] uppercase">Synchronizing intelligence...</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
          <AnimatePresence mode="popLayout">
            {filteredTasks.length > 0 ? (
              filteredTasks.map((task) => (
                <motion.div
                  key={task.id}
                  layout
                  initial={{ opacity: 0, scale: 0.9 }}
                  animate={{ opacity: 1, scale: 1 }}
                  exit={{ opacity: 0, scale: 0.9 }}
                  transition={{ duration: 0.3 }}
                >
                  <QuestCard task={task} />
                </motion.div>
              ))
            ) : (
              <motion.div 
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                className="col-span-full py-20 text-center glass-panel rounded-3xl border-dashed border-white/10"
              >
                <ClipboardList className="w-12 h-12 text-slate-700 mx-auto mb-4" />
                <p className="text-slate-500 font-bold uppercase tracking-widest">No matching intelligence records found in Nexus.</p>
              </motion.div>
            )}
          </AnimatePresence>
        </div>
      )}
    </div>
  );
};

export default QuestsPage;

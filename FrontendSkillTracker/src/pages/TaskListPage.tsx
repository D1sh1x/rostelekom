import React, { useState, useEffect } from 'react';
import { Search, Filter, Plus, X } from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';
import api from '../services/api';
import type { Task, TaskStatus } from '../types';
import TaskCard from '../components/TaskCard';
import { useAuth } from '../context/AuthContext';
import { useNavigate } from 'react-router-dom';

const TaskListPage: React.FC = () => {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState('');
  const [status, setStatus] = useState<TaskStatus | ''>('');
  const [showFilters, setShowFilters] = useState(false);
  
  const { user } = useAuth();
  const navigate = useNavigate();

  const fetchTasks = async () => {
    setLoading(true);
    try {
      const params: any = {};
      if (search) params.search = search;
      if (status) params.status = status;
      
      const response = await api.get('/tasks', { params });
      setTasks(response.data || []);
    } catch (err) {
      console.error('Failed to fetch tasks:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    const timeoutId = setTimeout(fetchTasks, 300);
    return () => clearTimeout(timeoutId);
  }, [search, status]);

  const clearFilters = () => {
    setSearch('');
    setStatus('');
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h1 className="text-3xl font-bold">All Tasks</h1>
          <p className="text-text-secondary mt-1">Manage and track all project activities</p>
        </div>
        {user?.role === 'manager' && (
          <button 
            onClick={() => navigate('/tasks/new')}
            className="primary flex items-center gap-2"
          >
            <Plus size={20} />
            Create Task
          </button>
        )}
      </div>

      <div className="flex flex-col md:flex-row gap-4">
        <div className="relative flex-1">
          <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-text-muted">
            <Search size={18} />
          </div>
          <input
            type="text"
            placeholder="Search tasks by title or description..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="pl-10"
          />
        </div>
        <button 
          onClick={() => setShowFilters(!showFilters)}
          className={`flex items-center gap-2 px-4 ${showFilters ? 'border-primary text-primary' : ''}`}
        >
          <Filter size={18} />
          Filters
          {(status) && (
            <span className="w-2 h-2 bg-primary rounded-full" />
          )}
        </button>
      </div>

      <AnimatePresence>
        {showFilters && (
          <motion.div
            initial={{ height: 0, opacity: 0 }}
            animate={{ height: 'auto', opacity: 1 }}
            exit={{ height: 0, opacity: 0 }}
            className="overflow-hidden"
          >
            <div className="card glass p-4 flex flex-wrap items-center gap-4">
              <div className="flex flex-col gap-1.5 min-w-[150px]">
                <label className="text-xs font-semibold text-text-muted uppercase">Status</label>
                <select 
                  value={status} 
                  onChange={(e) => setStatus(e.target.value as any)}
                  className="h-10 text-sm"
                >
                  <option value="">All Statuses</option>
                  <option value="pending">Pending</option>
                  <option value="in_progress">In Progress</option>
                  <option value="completed">Completed</option>
                </select>
              </div>

              {(status || search) && (
                <button 
                  onClick={clearFilters}
                  className="mt-auto h-10 flex items-center gap-2 text-sm text-text-muted hover:text-white"
                >
                  <X size={16} />
                  Clear Filters
                </button>
              )}
            </div>
          </motion.div>
        )}
      </AnimatePresence>

      {loading ? (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {[1, 2, 3, 4, 5, 6].map(i => (
            <div key={i} className="h-48 bg-surface rounded-xl animate-pulse" />
          ))}
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {tasks.map(task => (
            <TaskCard key={task.id} task={task} />
          ))}
          {tasks.length === 0 && (
            <div className="col-span-full py-24 text-center">
              <div className="inline-flex items-center justify-center w-16 h-16 bg-bg-accent rounded-full mb-4 text-text-muted">
                <Search size={32} />
              </div>
              <h3 className="text-xl font-semibold">No tasks found</h3>
              <p className="text-text-secondary mt-2">Try adjusting your filters or search terms.</p>
            </div>
          )}
        </div>
      )}
    </div>
  );
};

export default TaskListPage;

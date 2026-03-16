import React, { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { 
  TrendingUp, 
  CheckCircle2, 
  Clock, 
  AlertCircle,
  Plus
} from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import api from '../services/api';
import type { Task } from '../types';
import TaskCard from '../components/TaskCard';
import { useAuth } from '../context/AuthContext';

const Dashboard: React.FC = () => {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState(true);
  const { user } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    const fetchTasks = async () => {
      try {
        const response = await api.get('/tasks');
        setTasks(response.data || []);
      } catch (err) {
        console.error('Failed to fetch tasks:', err);
      } finally {
        setLoading(false);
      }
    };
    fetchTasks();
  }, []);

  const stats = {
    total: tasks.length,
    pending: tasks.filter(t => t.status === 'pending').length,
    inProgress: tasks.filter(t => t.status === 'in_progress').length,
    completed: tasks.filter(t => t.status === 'completed').length,
  };

  const container = {
    hidden: { opacity: 0 },
    show: {
      opacity: 1,
      transition: {
        staggerChildren: 0.1
      }
    }
  };

  const item = {
    hidden: { y: 20, opacity: 0 },
    show: { y: 0, opacity: 1 }
  };

  if (loading) {
    return <div className="animate-pulse space-y-8">
      <div className="h-32 bg-surface rounded-xl"></div>
      <div className="grid grid-cols-3 gap-6">
        <div className="h-40 bg-surface rounded-xl"></div>
        <div className="h-40 bg-surface rounded-xl"></div>
        <div className="h-40 bg-surface rounded-xl"></div>
      </div>
    </div>;
  }

  return (
    <div className="space-y-8">
      <header className="flex justify-between items-end">
        <div>
          <h1 className="text-3xl font-bold">Welcome back, {user?.username}</h1>
          <p className="text-text-secondary mt-1">Here's what's happening with your projects today.</p>
        </div>
        {user?.role === 'manager' && (
          <button 
            onClick={() => navigate('/tasks/new')}
            className="primary flex items-center gap-2"
          >
            <Plus size={20} />
            Create New Task
          </button>
        )}
      </header>

      {/* Stats Grid */}
      <motion.div 
        variants={container}
        initial="hidden"
        animate="show"
        className="grid grid-cols-1 md:grid-cols-4 gap-6"
      >
        <motion.div variants={item} className="card border-l-4 border-primary">
          <div className="flex items-center justify-between mb-2">
            <span className="text-text-secondary text-sm font-medium">Total Tasks</span>
            <TrendingUp size={20} className="text-primary" />
          </div>
          <p className="text-3xl font-bold">{stats.total}</p>
        </motion.div>

        <motion.div variants={item} className="card border-l-4 border-warning">
          <div className="flex items-center justify-between mb-2">
            <span className="text-text-secondary text-sm font-medium">Pending</span>
            <Clock size={20} className="text-warning" />
          </div>
          <p className="text-3xl font-bold">{stats.pending}</p>
        </motion.div>

        <motion.div variants={item} className="card border-l-4 border-primary">
          <div className="flex items-center justify-between mb-2">
            <span className="text-text-secondary text-sm font-medium">In Progress</span>
            <AlertCircle size={20} className="text-primary" />
          </div>
          <p className="text-3xl font-bold">{stats.inProgress}</p>
        </motion.div>

        <motion.div variants={item} className="card border-l-4 border-success">
          <div className="flex items-center justify-between mb-2">
            <span className="text-text-secondary text-sm font-medium">Completed</span>
            <CheckCircle2 size={20} className="text-success" />
          </div>
          <p className="text-3xl font-bold">{stats.completed}</p>
        </motion.div>
      </motion.div>

      {/* Recent Tasks */}
      <section>
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-xl font-bold">Recent Tasks</h2>
          <button 
            onClick={() => navigate('/tasks')}
            className="text-sm text-primary hover:underline"
          >
            View all
          </button>
        </div>
        
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {tasks.slice(0, 6).map(task => (
            <TaskCard key={task.id} task={task} />
          ))}
          {tasks.length === 0 && (
            <div className="col-span-3 card text-center py-12">
              <p className="text-text-secondary">No tasks found. Get started by creating one!</p>
            </div>
          )}
        </div>
      </section>
    </div>
  );
};

export default Dashboard;

import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Save, X, Calendar, User as UserIcon, Type, FileText } from 'lucide-react';
import api from '../services/api';
import type { User } from '../types';

const CreateTaskPage: React.FC = () => {
  const navigate = useNavigate();
  const [employees, setEmployees] = useState<User[]>([]);
  const [loading, setLoading] = useState(false);
  
  const [formData, setFormData] = useState({
    title: '',
    description: '',
    employee_id: '',
    deadline: '',
    progress: 0
  });

  useEffect(() => {
    const fetchEmployees = async () => {
      try {
        const response = await api.get('/users');
        setEmployees(response.data.filter((u: User) => u.role === 'employee'));
      } catch (err) {
        console.error('Failed to fetch employees:', err);
      }
    };
    fetchEmployees();
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    
    try {
      await api.post('/tasks', {
        ...formData,
        employee_id: parseInt(formData.employee_id),
        deadline: new Date(formData.deadline).toISOString(),
      });
      navigate('/tasks');
    } catch (err) {
      console.error('Failed to create task:', err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-3xl mx-auto space-y-8">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">New Task</h1>
          <p className="text-text-secondary mt-1">Assign a new professional objective</p>
        </div>
        <button 
          onClick={() => navigate(-1)}
          className="p-2 text-text-muted hover:text-white transition-colors"
        >
          <X size={24} />
        </button>
      </div>

      <form onSubmit={handleSubmit} className="space-y-6">
        <div className="card glass p-8 space-y-6">
          <div className="space-y-2">
            <label className="text-sm font-semibold text-text-secondary uppercase">Task Title</label>
            <div className="relative">
              <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-text-muted">
                <Type size={18} />
              </div>
              <input
                required
                placeholder="e.g., Q1 Security Audit"
                value={formData.title}
                onChange={(e) => setFormData({ ...formData, title: e.target.value })}
                className="pl-10"
              />
            </div>
          </div>

          <div className="space-y-2">
            <label className="text-sm font-semibold text-text-secondary uppercase">Description</label>
            <div className="relative">
              <div className="absolute top-3 left-3 pointer-events-none text-text-muted">
                <FileText size={18} />
              </div>
              <textarea
                required
                rows={4}
                placeholder="Detail the scope and requirements..."
                value={formData.description}
                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                className="pl-10 resize-none pt-3"
              />
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="space-y-2">
              <label className="text-sm font-semibold text-text-secondary uppercase">Assign To</label>
              <div className="relative">
                <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-text-muted">
                  <UserIcon size={18} />
                </div>
                <select
                  required
                  value={formData.employee_id}
                  onChange={(e) => setFormData({ ...formData, employee_id: e.target.value })}
                  className="pl-10 h-11"
                >
                  <option value="">Select Employee</option>
                  {employees.map(emp => (
                    <option key={emp.id} value={emp.id}>{emp.username}</option>
                  ))}
                </select>
              </div>
            </div>

            <div className="space-y-2">
              <label className="text-sm font-semibold text-text-secondary uppercase">Deadline</label>
              <div className="relative">
                <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-text-muted">
                  <Calendar size={18} />
                </div>
                <input
                  type="date"
                  required
                  value={formData.deadline}
                  onChange={(e) => setFormData({ ...formData, deadline: e.target.value })}
                  className="pl-10 h-11"
                />
              </div>
            </div>
          </div>

          <div className="space-y-2">
            <label className="text-sm font-semibold text-text-secondary uppercase">Initial Progress ({formData.progress}%)</label>
            <input 
              type="range"
              min="0"
              max="100"
              value={formData.progress}
              onChange={(e) => setFormData({ ...formData, progress: parseInt(e.target.value) })}
              className="accent-primary h-2 bg-bg-accent rounded-lg appearance-none cursor-pointer"
            />
          </div>
        </div>

        <div className="flex gap-4">
          <button 
            type="submit" 
            disabled={loading}
            className="primary flex-1 h-12 flex items-center justify-center gap-2"
          >
            {loading ? (
              <div className="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
            ) : (
              <>
                <Save size={18} />
                Create Task
              </>
            )}
          </button>
          <button 
            type="button"
            onClick={() => navigate(-1)}
            className="flex-1 border border-border h-12 hover:bg-bg-accent"
          >
            Cancel
          </button>
        </div>
      </form>
    </div>
  );
};

export default CreateTaskPage;

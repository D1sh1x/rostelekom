import React, { useState, useEffect, useRef } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { motion } from 'framer-motion';
import { 
  ArrowLeft, 
  Calendar, 
  User as UserIcon, 
  Paperclip,
  Send,
  Upload,
  Clock,
  CheckCircle,
  AlertCircle
} from 'lucide-react';
import api from '../services/api';
import type { Task, TaskHistory, Comment, TaskStatus } from '../types';
import Badge from '../components/Badge';

const TaskDetailPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  
  const [task, setTask] = useState<Task | null>(null);
  const [history, setHistory] = useState<TaskHistory[]>([]);
  const [comments, setComments] = useState<Comment[]>([]);
  const [loading, setLoading] = useState(true);
  const [newComment, setNewComment] = useState('');
  const [activeTab, setActiveTab] = useState<'comments' | 'history' | 'attachments'>('comments');
  const [uploading, setUploading] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const fetchData = async () => {
    try {
      const [taskRes, historyRes, commentsRes] = await Promise.all([
        api.get(`/tasks/${id}`),
        api.get(`/tasks/${id}/history`),
        api.get(`/tasks/${id}/comments`)
      ]);
      setTask(taskRes.data);
      setHistory(historyRes.data || []);
      setComments(commentsRes.data || []);
    } catch (err) {
      console.error('Failed to fetch task details:', err);
      navigate('/tasks');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, [id]);

  const handleStatusChange = async (newStatus: TaskStatus) => {
    if (!task) return;
    try {
      await api.put(`/tasks/${task.id}`, {
        ...task,
        status: newStatus
      });
      fetchData();
    } catch (err) {
      console.error('Failed to update status:', err);
    }
  };

  const handleCommentSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newComment.trim() || !task) return;
    try {
      await api.post('/comments', {
        task_id: task.id,
        content: newComment
      });
      setNewComment('');
      fetchData();
    } catch (err) {
      console.error('Failed to post comment:', err);
    }
  };

  const handleFileUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file || !task) return;
    
    setUploading(true);
    const formData = new FormData();
    formData.append('file', file);
    
    try {
      await api.post(`/tasks/${task.id}/attachments`, formData, {
        headers: { 'Content-Type': 'multipart/form-data' }
      });
      fetchData();
    } catch (err) {
      console.error('Upload failed:', err);
    } finally {
      setUploading(false);
    }
  };

  if (loading || !task) return <div className="p-8 text-center animate-pulse">Loading task...</div>;

  return (
    <div className="space-y-8">
      <button 
        onClick={() => navigate(-1)}
        className="flex items-center gap-2 text-text-secondary hover:text-white transition-colors"
      >
        <ArrowLeft size={20} />
        Back to tasks
      </button>

      <div className="flex flex-col lg:flex-row gap-8">
        {/* Main Content */}
        <div className="flex-1 space-y-8">
          <section className="card glass p-8">
            <div className="flex flex-wrap justify-between items-start gap-4 mb-6">
              <div className="space-y-2">
                <Badge variant={task.status}>{task.status.replace('_', ' ')}</Badge>
                <h1 className="text-4xl font-bold">{task.title}</h1>
              </div>
              
              <div className="flex items-center gap-3">
                <select 
                  value={task.status}
                  onChange={(e) => handleStatusChange(e.target.value as TaskStatus)}
                  className="w-40 h-10 text-sm"
                >
                  <option value="pending">Pending</option>
                  <option value="in_progress">In Progress</option>
                  <option value="completed">Completed</option>
                </select>
              </div>
            </div>

            <p className="text-text-secondary text-lg leading-relaxed mb-8 whitespace-pre-wrap">
              {task.description}
            </p>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-6 border-t border-border pt-8">
              <div className="space-y-1">
                <span className="text-xs font-bold text-text-muted uppercase">Deadline</span>
                <div className="flex items-center gap-2 font-medium">
                  <Calendar size={18} className="text-primary" />
                  {new Date(task.deadline).toLocaleDateString(undefined, { dateStyle: 'long' })}
                </div>
              </div>
              <div className="space-y-1">
                <span className="text-xs font-bold text-text-muted uppercase">Assignee</span>
                <div className="flex items-center gap-2 font-medium">
                  <UserIcon size={18} className="text-primary" />
                  Employee #{task.employee_id}
                </div>
              </div>
              <div className="space-y-1">
                <span className="text-xs font-bold text-text-muted uppercase">Creator</span>
                <div className="flex items-center gap-2 font-medium">
                  <UserIcon size={18} className="text-primary" />
                  Manager #{task.creator_id}
                </div>
              </div>
            </div>

            <div className="mt-8 space-y-2">
              <div className="flex justify-between text-sm">
                <span className="font-medium">Completion Progress</span>
                <span className="text-primary font-bold">{task.progress}%</span>
              </div>
              <div className="w-full bg-bg-accent h-3 rounded-full overflow-hidden">
                <motion.div 
                  initial={{ width: 0 }}
                  animate={{ width: `${task.progress}%` }}
                  className="h-full bg-primary"
                />
              </div>
            </div>
          </section>

          {/* Interactions Section */}
          <div className="card glass overflow-hidden">
            <div className="flex border-b border-border">
              {(['comments', 'history', 'attachments'] as const).map(tab => (
                <button
                  key={tab}
                  onClick={() => setActiveTab(tab)}
                  className={`px-8 py-4 text-sm font-bold uppercase tracking-wider transition-all relative ${
                    activeTab === tab ? 'text-primary' : 'text-text-muted hover:text-text-secondary'
                  }`}
                >
                  {tab}
                  {activeTab === tab && (
                    <motion.div layoutId="activeTab" className="absolute bottom-0 left-0 right-0 h-1 bg-primary" />
                  )}
                </button>
              ))}
            </div>

            <div className="p-8">
              {activeTab === 'comments' && (
                <div className="space-y-8">
                  <form onSubmit={handleCommentSubmit} className="flex gap-4">
                    <input 
                      value={newComment}
                      onChange={(e) => setNewComment(e.target.value)}
                      placeholder="Add a thought or update..."
                      className="flex-1 bg-bg-accent"
                    />
                    <button type="submit" className="primary px-6 flex items-center gap-2">
                      <Send size={18} />
                      Post
                    </button>
                  </form>

                  <div className="space-y-6">
                    {comments.map(comment => (
                      <div key={comment.id} className="flex gap-4">
                        <div className="w-10 h-10 rounded-full bg-primary/20 flex items-center justify-center font-bold text-primary shrink-0">
                          {comment.user?.username.substring(0, 1).toUpperCase() || '?'}
                        </div>
                        <div className="space-y-1">
                          <div className="flex items-center gap-2">
                            <span className="font-bold">{comment.user?.username || `User ${comment.user_id}`}</span>
                            <span className="text-xs text-text-muted">
                              {new Date(comment.created_at).toLocaleString()}
                            </span>
                          </div>
                          <p className="text-text-secondary">{comment.content}</p>
                        </div>
                      </div>
                    ))}
                    {comments.length === 0 && <p className="text-center text-text-muted py-8">No comments yet.</p>}
                  </div>
                </div>
              )}

              {activeTab === 'history' && (
                <div className="space-y-6">
                  {history.map(h => (
                    <div key={h.id} className="flex items-start gap-4 border-l-2 border-primary/20 pl-6 relative">
                      <div className="absolute -left-[9px] top-0 w-4 h-4 rounded-full bg-bg-main border-2 border-primary" />
                      <div className="space-y-1">
                        <div className="flex items-center gap-2 text-sm">
                          <span className="text-text-muted font-medium">
                            {new Date(h.created_at).toLocaleString()}
                          </span>
                        </div>
                        <p className="font-medium">
                          Status changed from <Badge variant={h.old_status as any}>{h.old_status}</Badge> to <Badge variant={h.new_status as any}>{h.new_status}</Badge>
                        </p>
                        <p className="text-xs text-text-muted">Changed by User #{h.changed_by}</p>
                      </div>
                    </div>
                  ))}
                  {history.length === 0 && <p className="text-center text-text-muted py-8">No history available.</p>}
                </div>
              )}

              {activeTab === 'attachments' && (
                <div className="space-y-6">
                  <div className="flex justify-between items-center">
                    <h3 className="text-lg font-bold">Files</h3>
                    <button 
                      onClick={() => fileInputRef.current?.click()}
                      disabled={uploading}
                      className="text-sm bg-primary/10 text-primary hover:bg-primary/20 flex items-center gap-2"
                    >
                      {uploading ? <Clock size={16} className="animate-spin" /> : <Upload size={16} />}
                      Upload File
                    </button>
                    <input 
                      type="file" 
                      className="hidden" 
                      ref={fileInputRef}
                      onChange={handleFileUpload}
                    />
                  </div>

                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    {task.attachments?.map(file => (
                      <div key={file.id} className="card p-4 flex items-center gap-4 hover:bg-bg-accent transition-colors">
                        <div className="w-12 h-12 bg-surface rounded-lg flex items-center justify-center text-primary">
                          <Paperclip size={24} />
                        </div>
                        <div className="flex-1 overflow-hidden">
                          <p className="font-medium truncate">{file.file_name}</p>
                          <p className="text-xs text-text-muted">{(file.file_size / 1024).toFixed(1)} KB</p>
                        </div>
                        <a 
                          href={`http://localhost:8080/${file.file_path}`} 
                          target="_blank" 
                          rel="noreferrer"
                          className="p-2 hover:text-primary transition-colors"
                        >
                          Download
                        </a>
                      </div>
                    ))}
                    {(!task.attachments || task.attachments.length === 0) && (
                      <div className="col-span-2 py-12 text-center text-text-muted">
                        No files attached.
                      </div>
                    )}
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>

        {/* Sidebar Info */}
        <div className="w-full lg:w-80 space-y-6">
          <div className="card glass p-6">
            <h3 className="font-bold mb-4 flex items-center gap-2">
              <CheckCircle size={18} className="text-success" />
              Quick Actions
            </h3>
            <div className="space-y-3">
              <button 
                onClick={() => handleStatusChange('completed')}
                className="w-full text-left px-4 py-3 rounded-lg border border-border hover:border-success/50 hover:bg-success/5 transition-all text-sm font-medium flex items-center justify-between group"
              >
                Mark as Completed
                <CheckCircle size={16} className="text-text-muted group-hover:text-success" />
              </button>
              <button 
                className="w-full text-left px-4 py-3 rounded-lg border border-border hover:border-primary/50 hover:bg-primary/5 transition-all text-sm font-medium flex items-center justify-between group"
              >
                Set Reminder
                <Clock size={16} className="text-text-muted group-hover:text-primary" />
              </button>
            </div>
          </div>

          <div className="card glass p-6 bg-primary/5 border-primary/20">
            <h3 className="font-bold mb-2 flex items-center gap-2 text-primary">
              <AlertCircle size={18} />
              Tips
            </h3>
            <p className="text-xs text-text-secondary leading-relaxed">
              Don't forget to upload your work progress in the attachments tab for your manager's review.
            </p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default TaskDetailPage;

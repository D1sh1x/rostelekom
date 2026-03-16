import React from 'react';
import { useNavigate } from 'react-router-dom';
import { Calendar, User as UserIcon } from 'lucide-react';
import type { Task } from '../types';
import Badge from './Badge';

interface TaskCardProps {
  task: Task;
}

const TaskCard: React.FC<TaskCardProps> = ({ task }) => {
  const navigate = useNavigate();

  return (
    <div
      onClick={() => navigate(`/tasks/${task.id}`)}
      className="card hover:border-primary/50 cursor-pointer group"
    >
      <div className="flex justify-between items-start mb-4">
        <Badge variant={task.status}>
          {task.status.replace('_', ' ')}
        </Badge>
        <span className="text-xs text-text-muted flex items-center gap-1">
          <Calendar size={14} />
          {new Date(task.deadline).toLocaleDateString()}
        </span>
      </div>

      <h3 className="text-lg font-semibold mb-2 group-hover:text-primary transition-colors">
        {task.title}
      </h3>

      <p className="text-text-secondary text-sm line-clamp-2 mb-4">
        {task.description}
      </p>

      <div className="flex items-center justify-between border-t border-border pt-4 mt-auto">
        <div className="flex items-center gap-2 text-xs text-text-muted">
          <UserIcon size={14} />
          <span>ID: {task.employee_id}</span>
        </div>
        <div className="w-24 bg-bg-accent h-1.5 rounded-full overflow-hidden">
          <div
            className="h-full bg-primary transition-all duration-500"
            style={{ width: `${task.progress}%` }}
          />
        </div>
      </div>
    </div>
  );
};

export default TaskCard;

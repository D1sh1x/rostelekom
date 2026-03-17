export type Role = 'manager' | 'employee';
export type TaskStatus = 'pending' | 'in_progress' | 'completed';

export interface User {
  id: number;
  username: string;
  role: Role;
  name: string;
  created_at: string;
}

export interface FileAttachment {
  id: number;
  task_id: number;
  file_name: string;
  file_path: string;
  file_size: number;
  uploaded_at: string;
}

export interface TaskStatusHistory {
  id: number;
  task_id: number;
  old_status: TaskStatus;
  new_status: TaskStatus;
  changed_by: number;
  created_at: string;
  user?: User;
}

export interface Task {
  id: number;
  employee_id: number;
  creator_id: number;
  title: string;
  description: string;
  deadline: string;
  status: TaskStatus;
  progress: number;
  created_at: string;
  updated_at: string;
  employee?: User;
  creator?: User;
  attachments?: FileAttachment[];
  history?: TaskStatusHistory[];
}

export interface Comment {
  id: number;
  taskId: number;
  userId: number;
  text: string;
  createdAt: string;
  user?: User;
}

export interface TaskFilters {
  status?: TaskStatus;
  employeeId?: number;
  creatorId?: number;
  search?: string;
  startDate?: string;
  endDate?: string;
}

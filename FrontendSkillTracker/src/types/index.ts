export type UserRole = 'manager' | 'employee';

export interface User {
  id: number;
  username: string;
  role: UserRole;
  created_at: string;
}

export type TaskStatus = 'pending' | 'in_progress' | 'completed';

export interface FileAttachment {
  id: number;
  task_id: number;
  file_name: string;
  file_path: string;
  file_size: number;
  uploaded_at: string;
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
  attachments?: FileAttachment[];
}

export interface TaskHistory {
  id: number;
  task_id: number;
  old_status: string;
  new_status: string;
  changed_by: number;
  created_at: string;
}

export interface Comment {
  id: number;
  task_id: number;
  user_id: number;
  content: string;
  created_at: string;
  updated_at: string;
  user?: User;
}

export interface AuthResponse {
  access_token: string;
  refresh_token: string;
  user: User;
}

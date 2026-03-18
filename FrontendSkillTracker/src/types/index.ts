export type Role = 'manager' | 'employee'
export type TaskStatus = 'pending' | 'in_progress' | 'completed'

export interface User {
  id: number
  username: string
  role: Role
  name: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  access_token: string
  refresh_token: string
  user: User
}

export interface Skill {
  id: number
  name: string
  description: string
  created_at: string
}

export interface Task {
  id: number
  employee_id: number
  creator_id: number
  title: string
  description: string
  deadline: string
  status: TaskStatus
  progress: number
  required_skills: Skill[]
  created_at: string
  updated_at: string
}

export interface TaskRequest {
  employee_id: number
  title: string
  description?: string
  deadline: string
  progress?: number
  status?: TaskStatus
}

export interface TaskFilter {
  status?: TaskStatus | ''
  employee_id?: number
  creator_id?: number
  search?: string
  from_date?: string
  to_date?: string
}

export interface TaskHistory {
  id: number
  task_id: number
  old_status: TaskStatus
  new_status: TaskStatus
  changed_by: number
  created_at: string
}

export interface Comment {
  id: number
  task_id: number
  user_id: number
  text: string
  created_at: string
}

export interface Attachment {
  id: number
  task_id: number
  file_name: string
  file_size: number
  uploaded_at: string
}

export interface RecommendedEmployee {
  id: number
  username: string
  name: string
  skills: Skill[]
  match_score: number
  matched_skills: string[]
  missing_skills: string[]
}

export interface UserRequest {
  username: string
  password: string
  role: Role
  name: string
}

export interface UpdateUserRequest {
  username: string
  password?: string
  role: Role
  name: string
}

export interface SkillRequest {
  name: string
  description?: string
}

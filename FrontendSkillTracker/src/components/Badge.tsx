import React from 'react';

interface BadgeProps {
  children: React.ReactNode;
  variant?: 'pending' | 'in_progress' | 'completed' | 'primary' | 'role';
  className?: string;
}

const Badge: React.FC<BadgeProps> = ({ children, variant = 'primary', className = '' }) => {
  const variants = {
    pending: 'bg-warning/10 text-warning border-warning/20',
    in_progress: 'bg-primary/10 text-primary border-primary/20',
    completed: 'bg-success/10 text-success border-success/20',
    primary: 'bg-primary/20 text-primary border-primary/30',
    role: 'bg-bg-accent text-text-secondary border-border'
  };

  return (
    <span className={`px-2.5 py-0.5 rounded-full text-xs font-semibold border ${variants[variant]} ${className}`}>
      {children}
    </span>
  );
};

export default Badge;

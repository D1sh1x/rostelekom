import React, { useState, useEffect } from 'react';
import { UserPlus, Shield, ShieldCheck } from 'lucide-react';
import { motion } from 'framer-motion';
import api from '../services/api';
import type { User } from '../types';

const UsersPage: React.FC = () => {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchUsers = async () => {
      try {
        const response = await api.get('/users');
        setUsers(response.data || []);
      } catch (err) {
        console.error('Failed to fetch users:', err);
      } finally {
        setLoading(false);
      }
    };
    fetchUsers();
  }, []);

  const container = {
    hidden: { opacity: 0 },
    show: {
      opacity: 1,
      transition: {
        staggerChildren: 0.05
      }
    }
  };

  const item = {
    hidden: { x: -20, opacity: 0 },
    show: { x: 0, opacity: 1 }
  };

  return (
    <div className="space-y-8">
      <div className="flex justify-between items-end">
        <div>
          <h1 className="text-3xl font-bold">Team Members</h1>
          <p className="text-text-secondary mt-1">Manage employees and their roles</p>
        </div>
        <button className="primary flex items-center gap-2">
          <UserPlus size={20} />
          Add Member
        </button>
      </div>

      <div className="card glass overflow-hidden">
        <table className="w-full text-left border-collapse">
          <thead>
            <tr className="bg-bg-accent/50 text-text-muted text-xs uppercase tracking-wider">
              <th className="px-6 py-4 font-semibold">User</th>
              <th className="px-6 py-4 font-semibold">Role</th>
              <th className="px-6 py-4 font-semibold">Member Since</th>
              <th className="px-6 py-4 font-semibold text-right">Actions</th>
            </tr>
          </thead>
          <motion.tbody 
            variants={container}
            initial="hidden"
            animate="show"
            className="divide-y divide-border"
          >
            {users.map(user => (
              <motion.tr key={user.id} variants={item} className="hover:bg-bg-accent/30 transition-colors group">
                <td className="px-6 py-4">
                  <div className="flex items-center gap-3">
                    <div className="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center font-bold text-primary">
                      {user.username.substring(0, 2).toUpperCase()}
                    </div>
                    <div>
                      <p className="font-semibold">{user.username}</p>
                      <p className="text-xs text-text-muted">ID: #{user.id}</p>
                    </div>
                  </div>
                </td>
                <td className="px-6 py-4">
                  <div className="flex items-center gap-2">
                    {user.role === 'manager' ? (
                      <ShieldCheck size={16} className="text-primary" />
                    ) : (
                      <Shield size={16} className="text-success" />
                    )}
                    <span className={`text-sm capitalize ${user.role === 'manager' ? 'text-primary' : 'text-success'}`}>
                      {user.role}
                    </span>
                  </div>
                </td>
                <td className="px-6 py-4 text-sm text-text-secondary">
                  {new Date(user.created_at).toLocaleDateString()}
                </td>
                <td className="px-6 py-4 text-right">
                  <button className="text-sm text-text-muted hover:text-white transition-colors">
                    Edit Settings
                  </button>
                </td>
              </motion.tr>
            ))}
          </motion.tbody>
        </table>
        {loading && (
          <div className="py-20 flex flex-col items-center justify-center gap-4">
            <div className="w-10 h-10 border-4 border-primary border-t-transparent rounded-full animate-spin"></div>
            <p className="text-text-muted animate-pulse">Retrieving team data...</p>
          </div>
        )}
        {!loading && users.length === 0 && (
          <div className="py-20 text-center">
            <p className="text-text-muted">No team members found.</p>
          </div>
        )}
      </div>
    </div>
  );
};

export default UsersPage;

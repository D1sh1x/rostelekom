import React from 'react';
import { motion } from 'framer-motion';
import { Zap, Target, TrendingUp, History, Clock } from 'lucide-react';

const DashboardWidget: React.FC<{ 
  title: string; 
  value: string | number; 
  icon: React.ReactNode; 
  trend?: string;
  color: string;
}> = ({ title, value, icon, trend, color }) => (
  <motion.div 
    initial={{ opacity: 0, y: 20 }}
    animate={{ opacity: 1, y: 0 }}
    className="glass-panel p-6 rounded-2xl border-white/5 flex items-start justify-between group hover:border-white/10 transition-colors"
  >
    <div>
      <p className="text-[10px] font-bold text-slate-500 uppercase tracking-[0.2em] mb-2">{title}</p>
      <h3 className="text-3xl font-bold text-white mb-2">{value}</h3>
      {trend && (
        <div className="flex items-center gap-1">
          <TrendingUp className="w-3 h-3 text-green-500" />
          <span className="text-[10px] text-green-500 font-bold">{trend} SYNC RATE</span>
        </div>
      )}
    </div>
    <div className={`p-3 rounded-xl bg-${color}-500/10 border border-${color}-500/20 text-${color}-400 group-hover:scale-110 transition-transform`}>
      {icon}
    </div>
  </motion.div>
);

const DashboardPage: React.FC = () => {
  return (
    <div className="space-y-8 pb-10">
      <div className="flex items-end justify-between">
        <div>
          <h2 className="text-sm font-bold text-blue-500 uppercase tracking-[0.3em] mb-1">Commander's Interface</h2>
          <h1 className="text-4xl font-bold text-white tracking-tight">Operations <span className="text-glow">Overview</span></h1>
        </div>
        <div className="text-right">
          <p className="text-xs text-slate-500 font-bold uppercase tracking-widest mb-1">System Time</p>
          <div className="text-xl font-mono text-white flex items-center gap-2 justify-end">
            <Clock className="w-4 h-4 text-slate-500" />
            {new Date().toLocaleTimeString('en-US', { hour12: false })}
          </div>
        </div>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <DashboardWidget 
          title="Active Quests" 
          value="12" 
          icon={<Zap className="w-6 h-6" />}
          trend="+5.2%"
          color="blue"
        />
        <DashboardWidget 
          title="Sync Level" 
          value="84%" 
          icon={<TrendingUp className="w-6 h-6" />}
          color="purple"
        />
        <DashboardWidget 
          title="Agents Online" 
          value="8" 
          icon={<Target className="w-6 h-6" />}
          color="green"
        />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Recent Intelligence (History) */}
        <div className="lg:col-span-2 space-y-6">
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-bold text-white flex items-center gap-2">
              <History className="w-5 h-5 text-blue-500" /> RECENT INTELLIGENCE
            </h3>
            <button className="text-[10px] font-bold text-slate-500 hover:text-white uppercase tracking-widest">View Archives</button>
          </div>
          
          <div className="space-y-4">
            {[1, 2, 3].map((i) => (
              <motion.div 
                key={i}
                initial={{ opacity: 0, x: -20 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ delay: i * 0.1 }}
                className="glass-panel p-4 rounded-xl border-white/5 flex items-center gap-4 hover:bg-white/5 transition-colors cursor-pointer group"
              >
                <div className="w-10 h-10 rounded-lg bg-slate-800 flex items-center justify-center border border-white/5 text-slate-400 group-hover:text-blue-400">
                  <Clock className="w-5 h-5" />
                </div>
                <div className="flex-1">
                  <div className="flex items-center justify-between mb-1">
                    <h4 className="text-sm font-bold text-white">Quest Status Modification</h4>
                    <span className="text-[10px] font-mono text-slate-600 uppercase">T-minus {i*2}h</span>
                  </div>
                  <p className="text-xs text-slate-400">Agent {312 + i} updated Quest Alpha to <span className="text-blue-500 font-bold uppercase underline decoration-2 underline-offset-4">In Progress</span></p>
                </div>
              </motion.div>
            ))}
          </div>
        </div>

        {/* System Status */}
        <div className="space-y-6">
          <h3 className="text-lg font-bold text-white flex items-center gap-2">
            <Zap className="w-5 h-5 text-purple-500" /> SYSTEM ARCHITECTURE
          </h3>
          <div className="glass-panel p-6 rounded-2xl border-white/5 space-y-6">
            <div className="space-y-2">
              <div className="flex justify-between items-end">
                <span className="text-[10px] font-bold text-slate-500 uppercase tracking-widest">Network Stability</span>
                <span className="text-xs font-mono text-blue-400">99.8%</span>
              </div>
              <div className="h-1 bg-slate-800 rounded-full overflow-hidden">
                <div className="h-full bg-blue-500 w-[99.8%]" />
              </div>
            </div>
            <div className="space-y-2">
              <div className="flex justify-between items-end">
                <span className="text-[10px] font-bold text-slate-500 uppercase tracking-widest">Resource Load</span>
                <span className="text-xs font-mono text-purple-400">42%</span>
              </div>
              <div className="h-1 bg-slate-800 rounded-full overflow-hidden">
                <div className="h-full bg-purple-500 w-[42%]" />
              </div>
            </div>
            <button className="w-full py-2 bg-white/5 hover:bg-white/10 border border-white/5 rounded-lg text-[10px] font-bold text-slate-400 transition-all uppercase tracking-widest">
              Diagnostic Report
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default DashboardPage;

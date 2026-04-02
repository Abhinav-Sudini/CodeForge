import Link from 'next/link';
import { Terminal, Code2, UserCircle } from 'lucide-react';

export default function Navbar() {
  return (
    <nav className="w-full border-b border-white/10 bg-black/50 backdrop-blur-md sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          <Link href="/" className="flex items-center gap-2 group">
            <div className="p-2 bg-indigo-500/20 rounded-lg group-hover:bg-indigo-500/30 transition-colors">
              <Terminal className="w-5 h-5 text-indigo-400" />
            </div>
            <span className="font-bold text-xl tracking-tight text-white">
              Code<span className="text-indigo-400">Forge</span>
            </span>
          </Link>
          
          <div className="flex items-center gap-6">
            <Link href="/" className="text-sm font-medium text-gray-300 hover:text-white transition-colors flex items-center gap-1.5">
              <Code2 className="w-4 h-4" />
              Problems
            </Link>
            {/* Login placeholder - as requested */}
            <button className="flex items-center gap-2 px-4 py-2 text-sm font-medium text-white bg-white/5 hover:bg-white/10 rounded-full border border-white/10 transition-all cursor-not-allowed opacity-70">
              <UserCircle className="w-4 h-4" />
              Login
            </button>
          </div>
        </div>
      </div>
    </nav>
  );
}

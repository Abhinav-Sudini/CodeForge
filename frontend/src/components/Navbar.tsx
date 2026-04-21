"use client";

import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import {
  Terminal,
  Code2,
  UserCircle,
  Trophy,
  MessageSquare,
  Settings,
  LogOut,
} from "lucide-react";
import { useState, useEffect } from "react";

export default function Navbar() {
  const pathname = usePathname();
  const router = useRouter();
  const [isLoggedIn, setIsLoggedIn] = useState<boolean | null>(null);

  useEffect(() => {
    const checkAuth = () => {
      setIsLoggedIn(localStorage.getItem("cf_authed") === "1");
    };
    checkAuth();
    window.addEventListener("auth-change", checkAuth);
    window.addEventListener("storage", checkAuth);
    return () => {
      window.removeEventListener("auth-change", checkAuth);
      window.removeEventListener("storage", checkAuth);
    };
  }, []);

  const handleLogout = () => {
    localStorage.removeItem("cf_authed");
    setIsLoggedIn(false);
    router.push("/login");
  };

  if (pathname.startsWith("/problems/")) {
    return null;
  }

  return (
    <nav className="w-full border-b border-white/10 bg-black/50 backdrop-blur-md sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          <div className="flex items-center gap-10">
            <Link href="/" className="flex items-center gap-2 group">
              <div className="p-2 bg-indigo-500/20 rounded-lg group-hover:bg-indigo-500/30 transition-colors">
                <Terminal className="w-5 h-5 text-indigo-400" />
              </div>
              <span className="font-bold text-xl tracking-tight text-white">
                Code<span className="text-indigo-400">Forge</span>
              </span>
            </Link>

            <div className="hidden md:flex items-center gap-6">
              <Link
                href="/"
                className="text-sm font-bold text-white transition-colors flex items-center gap-1.5 border-b-2 border-indigo-500 pb-1 pt-1"
              >
                <Code2 className="w-4 h-4" />
                Problems
              </Link>
              <Link
                href="#"
                className="text-sm font-medium text-neutral-400 hover:text-white transition-colors flex items-center gap-1.5 pb-1 pt-1 border-b-2 border-transparent"
              >
                <Trophy className="w-4 h-4" />
                Contests
              </Link>
              <Link
                href="#"
                className="text-sm font-medium text-neutral-400 hover:text-white transition-colors flex items-center gap-1.5 pb-1 pt-1 border-b-2 border-transparent"
              >
                <MessageSquare className="w-4 h-4" />
                Discuss
              </Link>
              <Link
                href="#"
                className="text-sm font-medium text-neutral-400 hover:text-white transition-colors flex items-center gap-1.5 pb-1 pt-1 border-b-2 border-transparent"
              >
                <Settings className="w-4 h-4" />
                Profile
              </Link>
            </div>
          </div>

          <div className="flex items-center gap-4">
            {isLoggedIn ? (
              <button
                onClick={handleLogout}
                className="flex items-center gap-2 px-4 py-2 text-sm font-medium text-white bg-white/5 hover:bg-white/10 rounded-full border border-white/10 transition-all hover:shadow-[0_0_15px_rgba(99,102,241,0.3)]"
              >
                <LogOut className="w-4 h-4" />
                Logout
              </button>
            ) : (
              <Link
                href="/login"
                className="flex items-center gap-2 px-4 py-2 text-sm font-medium text-white bg-white/5 hover:bg-white/10 rounded-full border border-white/10 transition-all hover:shadow-[0_0_15px_rgba(99,102,241,0.3)]"
              >
                <UserCircle className="w-4 h-4" />
                Login
              </Link>
            )}
          </div>
        </div>
      </div>
    </nav>
  );
}

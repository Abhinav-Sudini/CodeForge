"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { UserCircle, Clock, LogOut, ArrowLeft } from "lucide-react";

export default function ProfilePage() {
  const router = useRouter();
  const [email, setEmail] = useState<string | null>(null);
  const [isLoggedIn, setIsLoggedIn] = useState<boolean | null>(null);

  useEffect(() => {
    const authed = localStorage.getItem("cf_authed") === "1";
    const storedEmail = localStorage.getItem("cf_email") || null;
    setIsLoggedIn(authed);
    setEmail(storedEmail);
  }, []);

  const handleLogout = () => {
    localStorage.removeItem("cf_authed");
    localStorage.removeItem("cf_email");
    window.dispatchEvent(new Event("auth-change"));
    router.push("/login");
  };

  if (isLoggedIn === null) return null;

  if (!isLoggedIn) {
    return (
      <div className="min-h-[60vh] flex flex-col items-center justify-center gap-6 text-center">
        <div className="w-20 h-20 rounded-full bg-white/5 border border-white/10 flex items-center justify-center">
          <UserCircle className="w-10 h-10 text-neutral-600" />
        </div>
        <div>
          <p className="text-white font-bold text-xl mb-2">You&apos;re not logged in</p>
          <p className="text-neutral-500 text-sm">Sign in to view your profile and submission history.</p>
        </div>
        <Link
          href="/login"
          className="px-6 py-2.5 bg-indigo-600 hover:bg-indigo-500 text-white text-sm font-bold rounded-xl transition-colors"
        >
          Login
        </Link>
      </div>
    );
  }

  const displayName = email ? email.split("@")[0] : "User";
  const avatarLetter = displayName[0]?.toUpperCase() ?? "U";

  return (
    <div className="max-w-3xl mx-auto py-12 px-4 sm:px-6 animate-fade-in-up">

      <div className="relative overflow-hidden rounded-3xl bg-gradient-to-br from-neutral-900 via-neutral-900 to-indigo-950/40 border border-white/5 p-8 shadow-2xl mb-8">
        <div className="absolute top-0 right-0 p-12 opacity-20 pointer-events-none">
          <div className="w-48 h-48 bg-indigo-500 rounded-full blur-[80px]" />
        </div>
        <div className="relative z-10 flex items-center gap-6">
          <div className="w-20 h-20 rounded-2xl bg-indigo-500/20 border border-indigo-500/30 flex items-center justify-center text-3xl font-extrabold text-indigo-300 shrink-0 shadow-lg shadow-indigo-500/10">
            {avatarLetter}
          </div>
          <div className="flex-1 min-w-0">
            <h1 className="text-2xl font-extrabold text-white tracking-tight truncate">{displayName}</h1>
            {email && (
              <p className="text-sm text-neutral-500 mt-1 truncate">{email}</p>
            )}
            <span className="inline-flex items-center gap-1.5 mt-2 px-2.5 py-1 bg-indigo-500/10 border border-indigo-500/20 rounded-full text-[11px] font-bold text-indigo-400 uppercase tracking-widest">
              <span className="w-1.5 h-1.5 rounded-full bg-indigo-400 animate-pulse" />
              Active
            </span>
          </div>
          <button
            onClick={handleLogout}
            className="flex items-center gap-2 px-4 py-2 text-sm font-medium text-neutral-400 hover:text-white bg-white/5 hover:bg-white/10 border border-white/10 rounded-xl transition-all shrink-0"
          >
            <LogOut className="w-4 h-4" />
            Logout
          </button>
        </div>
      </div>


      <div className="bg-neutral-900/60 border border-white/5 rounded-2xl p-6">
        <div className="flex items-start gap-3">
          <Clock className="w-5 h-5 text-amber-400 shrink-0 mt-0.5" />
          <div>
            <p className="text-sm font-semibold text-neutral-300 mb-1">Submission history</p>
            <p className="text-sm text-neutral-500">
              Per-problem submission history is available on each problem&apos;s page under the{" "}
              <strong className="text-neutral-400">Submissions</strong> tab. Detailed analytics will be added in a future update.
            </p>
          </div>
        </div>
      </div>

      <div className="mt-6 flex justify-center">
        <Link
          href="/"
          className="flex items-center gap-2 text-sm text-neutral-500 hover:text-neutral-300 transition-colors"
        >
          <ArrowLeft className="w-4 h-4" />
          Back to Problems
        </Link>
      </div>
    </div>
  );
}

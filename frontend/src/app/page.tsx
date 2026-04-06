import { fetchQuestions } from "@/lib/api";
import Link from "next/link";
import { ChevronRight, Target, Activity } from "lucide-react";

export const dynamic = 'force-dynamic'; // We want the questions to be fetched fresh!

export default async function Home() {
  const questions = await fetchQuestions();

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      {/* Hero Section */}
      <div className="relative mb-16 overflow-hidden rounded-3xl bg-gradient-to-br from-neutral-900 via-neutral-900 to-indigo-950/40 border border-white/5 p-8 sm:p-12">
        <div className="absolute top-0 right-0 p-12 opacity-20 pointer-events-none">
          <div className="w-64 h-64 bg-indigo-500 rounded-full blur-3xl" />
        </div>
        
        <div className="relative z-10">
          <h1 className="text-4xl sm:text-5xl font-extrabold tracking-tight text-white mb-4">
            Master your Craft.<br/>
            <span className="text-transparent bg-clip-text bg-gradient-to-r from-indigo-400 to-cyan-400">
              Forge your Code.
            </span>
          </h1>
          <p className="text-lg text-neutral-400 max-w-2xl mb-8">
            Dive into our curated list of competitive programming challenges. Pick a runtime, write optimal algorithms, and instantly climb the leaderboard.
          </p>
          <div className="flex flex-wrap gap-4">
            <div className="flex items-center gap-2 px-4 py-2 rounded-full bg-white/5 border border-white/10 text-sm font-medium text-neutral-300">
              <Target className="w-4 h-4 text-emerald-400" />
              {questions.length} Challenges Available
            </div>
            <div className="flex items-center gap-2 px-4 py-2 rounded-full bg-white/5 border border-white/10 text-sm font-medium text-neutral-300">
              <Activity className="w-4 h-4 text-blue-400" />
              6 Backend Runtimes Supported
            </div>
          </div>
        </div>
      </div>

      {/* Problem List Section */}
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-2xl font-bold text-white tracking-tight">Available Problem Set</h2>
      </div>

      <div className="bg-neutral-900/50 backdrop-blur-sm border border-white/5 rounded-2xl overflow-hidden shadow-2xl">
        <div className="overflow-x-auto">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="border-b border-white/10 bg-white/[0.02]">
                <th className="px-6 py-4 text-xs font-semibold text-neutral-400 uppercase tracking-wider">Status</th>
                <th className="px-6 py-4 text-xs font-semibold text-neutral-400 uppercase tracking-wider">Title</th>
                <th className="px-6 py-4 text-xs font-semibold text-neutral-400 uppercase tracking-wider">Category</th>
                <th className="px-6 py-4 text-xs font-semibold text-neutral-400 uppercase tracking-wider text-right">Action</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-white/5">
              {questions.length === 0 ? (
                <tr>
                  <td colSpan={4} className="px-6 py-12 text-center text-neutral-500">
                    No questions found. Make sure the backend scraper has finished!
                  </td>
                </tr>
              ) : (
                questions.map((q) => (
                  <tr key={q.QuestionId} className="group hover:bg-white/[0.02] transition-colors flex flex-col sm:table-row">
                    <td className="px-6 py-4 whitespace-nowrap hidden sm:table-cell">
                      <div className="w-2 h-2 rounded-full bg-neutral-600 group-hover:bg-indigo-400 transition-colors" />
                    </td>
                    <td className="px-6 py-4">
                      <Link href={`/problems/${q.QuestionId}`} className="text-base sm:text-sm font-medium text-neutral-200 group-hover:text-indigo-300 transition-colors">
                        {q.QuestionName}
                      </Link>
                    </td>
                    <td className="px-6 py-4">
                      <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-indigo-500/10 text-indigo-300 border border-indigo-500/20">
                        {q.QuestionCategory || "Algorithmic"}
                      </span>
                    </td>
                    <td className="px-6 py-4 sm:text-right">
                      <Link 
                        href={`/problems/${q.QuestionId}`}
                        className="inline-flex items-center px-3 py-1.5 text-xs font-semibold bg-white/5 hover:bg-white/10 border border-white/10 rounded-lg transition-all text-neutral-300 hover:text-white"
                      >
                        Solve
                        <ChevronRight className="w-3 h-3 ml-1" />
                      </Link>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}

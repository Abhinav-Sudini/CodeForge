"use client";

import { useState, useMemo } from "react";
import Link from "next/link";
import { Search, ChevronRight, ChevronLeft, ArrowUpDown, Target, Activity } from "lucide-react";
import { QuestionMinimal } from "@/lib/api";

interface ProblemsDashboardProps {
  questions: QuestionMinimal[];
}

const PAGE_SIZE = 20;

const CATEGORY_COLORS = [
  { hoverBg: "hover:bg-indigo-400/10", borderHover: "hover:border-indigo-400/40", dot: "bg-indigo-400", activeBg: "bg-indigo-500/20", activeBorder: "border-indigo-500" },
  { hoverBg: "hover:bg-emerald-400/10", borderHover: "hover:border-emerald-400/40", dot: "bg-emerald-400", activeBg: "bg-emerald-500/20", activeBorder: "border-emerald-500" },
  { hoverBg: "hover:bg-amber-400/10", borderHover: "hover:border-amber-400/40", dot: "bg-amber-400", activeBg: "bg-amber-500/20", activeBorder: "border-amber-500" },
  { hoverBg: "hover:bg-rose-400/10", borderHover: "hover:border-rose-400/40", dot: "bg-rose-400", activeBg: "bg-rose-500/20", activeBorder: "border-rose-500" },
  { hoverBg: "hover:bg-sky-400/10", borderHover: "hover:border-sky-400/40", dot: "bg-sky-400", activeBg: "bg-sky-500/20", activeBorder: "border-sky-500" },
  { hoverBg: "hover:bg-orange-400/10", borderHover: "hover:border-orange-400/40", dot: "bg-orange-400", activeBg: "bg-orange-500/20", activeBorder: "border-orange-500" },
  { hoverBg: "hover:bg-fuchsia-400/10", borderHover: "hover:border-fuchsia-400/40", dot: "bg-fuchsia-400", activeBg: "bg-fuchsia-500/20", activeBorder: "border-fuchsia-500" },
  { hoverBg: "hover:bg-lime-400/10", borderHover: "hover:border-lime-400/40", dot: "bg-lime-400", activeBg: "bg-lime-500/20", activeBorder: "border-lime-500" },
  { hoverBg: "hover:bg-pink-400/10", borderHover: "hover:border-pink-400/40", dot: "bg-pink-400", activeBg: "bg-pink-500/20", activeBorder: "border-pink-500" },
  { hoverBg: "hover:bg-teal-400/10", borderHover: "hover:border-teal-400/40", dot: "bg-teal-400", activeBg: "bg-teal-500/20", activeBorder: "border-teal-500" },
];

export default function ProblemsDashboard({ questions }: ProblemsDashboardProps) {
  const [searchQuery, setSearchQuery] = useState("");
  const [currentCategory, setCurrentCategory] = useState("All");
  const [sortMode, setSortMode] = useState<"id" | "name">("id");
  const [currentPage, setCurrentPage] = useState(1);

  
  const categories = useMemo(() => {
    const counts: Record<string, number> = { All: questions.length };
    questions.forEach((q) => {
      const cat = q.QuestionCategory || "Uncategorized";
      counts[cat] = (counts[cat] || 0) + 1;
    });
    return Object.entries(counts).sort((a, b) => b[1] - a[1]); // sort by count descending
  }, [questions]);

  // map colors
  const categoryColorMap = useMemo(() => {
    const map: Record<string, typeof CATEGORY_COLORS[0]> = {
      "All": { hoverBg: "hover:bg-white/10", borderHover: "hover:border-white/20", dot: "bg-white", activeBg: "bg-white/10", activeBorder: "border-white/30" }
    };
    let idx = 0;
    categories.forEach(([cat]) => {
      if (cat !== "All") {
        map[cat] = CATEGORY_COLORS[idx % CATEGORY_COLORS.length];
        idx++;
      }
    });
    return map;
  }, [categories]);

  
  const filteredQuestions = useMemo(() => {
    let result = questions;

    
    if (currentCategory !== "All") {
      result = result.filter(
        (q) => (q.QuestionCategory || "Uncategorized") === currentCategory
      );
    }

    
    if (searchQuery.trim() !== "") {
      const lowerQ = searchQuery.toLowerCase();
      result = result.filter(
        (q) =>
          q.QuestionName.toLowerCase().includes(lowerQ) ||
          q.QuestionId.toString().includes(lowerQ) ||
          (q.QuestionCategory || "Uncategorized").toLowerCase().includes(lowerQ)
      );
    }

    
    if (sortMode === "id") {
      result = [...result].sort((a, b) => a.QuestionId - b.QuestionId);
    } else {
      result = [...result].sort((a, b) =>
        a.QuestionName.localeCompare(b.QuestionName)
      );
    }

    return result;
  }, [questions, currentCategory, searchQuery, sortMode]);

  const totalPages = Math.max(1, Math.ceil(filteredQuestions.length / PAGE_SIZE));
  
  
  if (currentPage > totalPages) {
    setCurrentPage(totalPages);
  }

  const paginatedQuestions = useMemo(() => {
    const start = (currentPage - 1) * PAGE_SIZE;
    return filteredQuestions.slice(start, start + PAGE_SIZE);
  }, [filteredQuestions, currentPage]);

  return (
    <div className="w-full relative z-10 animate-fade-in-up">
      
      <div className="relative mb-12 overflow-hidden rounded-3xl bg-gradient-to-br from-neutral-900 via-neutral-900 to-indigo-950/40 border border-white/5 p-8 sm:p-12 shadow-2xl">
        <div className="absolute top-0 right-0 p-12 opacity-20 pointer-events-none">
          <div className="w-64 h-64 bg-indigo-500 rounded-full blur-[100px]" />
        </div>
        <div className="relative z-10">
          <p className="font-mono text-[10px] tracking-widest uppercase text-indigo-400 mb-4 opacity-80">practice arena</p>
          <h1 className="text-4xl sm:text-5xl font-extrabold tracking-tight text-white mb-4">
            Solve. <span className="text-transparent bg-clip-text bg-gradient-to-r from-indigo-400 to-cyan-400">Level up.</span>
          </h1>
          <p className="text-lg text-neutral-400 max-w-2xl mb-8 leading-relaxed">
            Pick a problem, write your solution, and sharpen your algorithmic thinking — one challenge at a time.
          </p>
          <div className="flex flex-wrap gap-4">
            <div className="flex flex-col gap-1 bg-black/40 border border-white/10 px-6 py-3 rounded-2xl backdrop-blur-md">
              <span className="font-mono text-2xl font-bold text-white">{questions.length}</span>
              <span className="text-[10px] font-bold uppercase tracking-widest text-neutral-500">Total</span>
            </div>
            <div className="flex flex-col gap-1 bg-black/40 border border-white/10 px-6 py-3 rounded-2xl backdrop-blur-md">
              <span className="font-mono text-2xl font-bold text-emerald-400">{categories.length - 1}</span>
              <span className="text-[10px] font-bold uppercase tracking-widest text-neutral-500">Categories</span>
            </div>
            <div className="flex flex-col gap-1 bg-black/40 border border-white/10 px-6 py-3 rounded-2xl backdrop-blur-md">
              <span className="font-mono text-2xl font-bold text-amber-400">{filteredQuestions.length}</span>
              <span className="text-[10px] font-bold uppercase tracking-widest text-neutral-500">Showing</span>
            </div>
          </div>
        </div>
      </div>

      
      <div className="flex flex-col md:flex-row md:items-center gap-4 mb-6">
        <div className="relative flex-1">
          <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <Search className="h-5 w-5 text-neutral-500" />
          </div>
          <input
            type="text"
            className="w-full bg-neutral-900/50 backdrop-blur-sm border border-white/10 rounded-xl py-2.5 pl-10 pr-4 text-white placeholder-neutral-500 focus:outline-none focus:ring-2 focus:ring-indigo-500/50 transition-all text-sm"
            placeholder="Search by name, ID, or category..."
            value={searchQuery}
            onChange={(e) => {
              setSearchQuery(e.target.value);
              setCurrentPage(1);
            }}
          />
        </div>

        <div className="flex items-center gap-2">
          <button
            onClick={() => setSortMode("id")}
            className={`flex items-center gap-2 px-4 py-2.5 text-sm font-semibold rounded-xl border transition-all ${
              sortMode === "id"
                ? "bg-indigo-500/10 border-indigo-500/30 text-indigo-400"
                : "bg-neutral-900/50 border-white/10 text-neutral-400 hover:text-white"
            }`}
          >
            Sort by ID
          </button>
          <button
            onClick={() => setSortMode("name")}
            className={`flex items-center gap-2 px-4 py-2.5 text-sm font-semibold rounded-xl border transition-all ${
              sortMode === "name"
                ? "bg-indigo-500/10 border-indigo-500/30 text-indigo-400"
                : "bg-neutral-900/50 border-white/10 text-neutral-400 hover:text-white"
            }`}
          >
            Sort by Name
            <ArrowUpDown className="w-4 h-4 ml-1 opacity-50" />
          </button>
        </div>
      </div>

      
      <div className="flex flex-wrap gap-2 mb-8">
        {categories.map(([cat, count]) => {
          const color = categoryColorMap[cat] || categoryColorMap["All"];
          const isActive = currentCategory === cat;
          return (
            <button
              key={cat}
              onClick={() => {
                setCurrentCategory(cat);
                setCurrentPage(1);
              }}
              className={`flex items-center gap-2 px-4 py-1.5 rounded-full text-xs font-semibold border transition-all ${
                isActive
                  ? `${color.activeBg} ${color.activeBorder} text-white shadow-[0_0_15px_rgba(79,70,229,0.15)]`
                  : `bg-transparent border-white/10 text-neutral-400 hover:text-white ${color.borderHover} ${color.hoverBg}`
              }`}
            >
              <span className={`w-1.5 h-1.5 rounded-full ${color.dot} shadow-[0_0_8px_currentColor] opacity-80`} />
              {cat}
              <span
                className={`px-1.5 py-0.5 rounded-full text-[10px] tabular-nums font-mono ${
                  isActive
                    ? "bg-black/30 text-white"
                    : "bg-black/20 text-current opacity-80"
                }`}
              >
                {count}
              </span>
            </button>
          );
        })}
      </div>

      
      <div className="bg-neutral-900/50 backdrop-blur-md border border-white/5 rounded-2xl overflow-hidden shadow-2xl transition-all duration-300">
        <div className="overflow-x-auto min-h-[400px]">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="border-b border-white/10 bg-white/[0.02]">
                <th className="px-6 py-4 text-xs font-semibold text-neutral-500 uppercase tracking-wider w-20">ID</th>
                <th className="px-6 py-4 text-xs font-semibold text-neutral-500 uppercase tracking-wider">Problem Title</th>
                <th className="px-6 py-4 text-xs font-semibold text-neutral-500 uppercase tracking-wider">Category</th>
                <th className="px-6 py-4 text-xs font-semibold text-neutral-500 uppercase tracking-wider text-right">Action</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-white/5">
              {paginatedQuestions.length === 0 ? (
                <tr>
                  <td colSpan={4} className="px-6 py-20 text-center">
                    <div className="flex flex-col items-center justify-center">
                      <div className="w-16 h-16 bg-white/5 rounded-full flex items-center justify-center mb-4">
                        <Search className="w-8 h-8 text-neutral-600" />
                      </div>
                      <p className="text-neutral-400 font-medium">No challenges found matching your criteria.</p>
                      <button 
                        onClick={() => { setSearchQuery(""); setCurrentCategory("All"); }}
                        className="mt-4 text-indigo-400 text-sm hover:text-indigo-300 transition-colors"
                      >
                        Clear filters
                      </button>
                    </div>
                  </td>
                </tr>
              ) : (
                paginatedQuestions.map((q) => (
                  <tr key={q.QuestionId} className="group hover:bg-indigo-500/5 transition-colors flex flex-col sm:table-row">
                    <td className="px-6 py-4 whitespace-nowrap hidden sm:table-cell">
                      <span className="font-mono text-xs text-neutral-500 bg-black/40 px-2 py-1 rounded inline-block">
                        #{q.QuestionId}
                      </span>
                    </td>
                    <td className="px-6 py-4">
                      <Link href={`/problems/${q.QuestionId}`} className="text-base sm:text-sm font-semibold text-neutral-200 group-hover:text-indigo-400 transition-colors mr-2">
                        {q.QuestionName}
                      </Link>
                      <span className="sm:hidden font-mono text-xs text-neutral-500 bg-black/40 px-2 py-1 rounded ml-2">
                        #{q.QuestionId}
                      </span>
                    </td>
                    <td className="px-6 py-4">
                      {(() => {
                        const cat = q.QuestionCategory || "Uncategorized";
                        const color = categoryColorMap[cat] || categoryColorMap["All"];
                        return (
                          <span className={`inline-flex items-center gap-2 px-2.5 py-1.5 rounded-md text-[11px] font-bold tracking-wide uppercase text-neutral-400 bg-white/5 border border-white/5 group-hover:border-white/10 transition-colors`}>
                            <span className={`w-1.5 h-1.5 rounded-full ${color.dot} opacity-80`} />
                            {cat}
                          </span>
                        );
                      })()}
                    </td>
                    <td className="px-6 py-4 sm:text-right">
                      <Link 
                        href={`/problems/${q.QuestionId}`}
                        className="inline-flex items-center px-4 py-2 text-xs font-bold bg-white/5 hover:bg-indigo-600 border border-white/10 hover:border-indigo-500 rounded-lg transition-all text-neutral-300 hover:text-white"
                      >
                        Solve
                        <ChevronRight className="w-3 h-3 ml-1.5" />
                      </Link>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>

      
      {filteredQuestions.length > 0 && (
        <div className="flex items-center justify-between mt-6">
          <div className="text-sm text-neutral-500 font-mono hidden sm:block">
            Showing {((currentPage - 1) * PAGE_SIZE) + 1} to {Math.min(currentPage * PAGE_SIZE, filteredQuestions.length)} of {filteredQuestions.length}
          </div>
          <div className="flex items-center gap-2">
            <button
              onClick={() => setCurrentPage(p => Math.max(1, p - 1))}
              disabled={currentPage === 1}
              className="flex items-center justify-center w-10 h-10 rounded-xl bg-neutral-900 border border-white/10 text-neutral-400 hover:text-white hover:border-white/30 disabled:opacity-30 disabled:cursor-not-allowed transition-all"
            >
              <ChevronLeft className="w-5 h-5" />
            </button>
            <div className="flex gap-1">
              
              {Array.from({ length: Math.min(5, totalPages) }, (_, i) => {
                let pageNum = currentPage;
                if (currentPage <= 3) pageNum = i + 1;
                else if (currentPage >= totalPages - 2) pageNum = totalPages - 4 + i;
                else pageNum = currentPage - 2 + i;

                if (pageNum < 1 || pageNum > totalPages) return null;

                return (
                  <button
                    key={pageNum}
                    onClick={() => setCurrentPage(pageNum)}
                    className={`w-10 h-10 rounded-xl font-mono text-sm transition-all border ${
                      currentPage === pageNum
                        ? "bg-indigo-600 border-indigo-500 text-white shadow-lg shadow-indigo-500/20"
                        : "bg-surface border-transparent text-neutral-400 hover:bg-white/5 hover:text-white"
                    }`}
                  >
                    {pageNum}
                  </button>
                );
              })}
            </div>
            <button
              onClick={() => setCurrentPage(p => Math.min(totalPages, p + 1))}
              disabled={currentPage === totalPages}
              className="flex items-center justify-center w-10 h-10 rounded-xl bg-neutral-900 border border-white/10 text-neutral-400 hover:text-white hover:border-white/30 disabled:opacity-30 disabled:cursor-not-allowed transition-all"
            >
              <ChevronRight className="w-5 h-5" />
            </button>
          </div>
        </div>
      )}
    </div>
  );
}

"use client";

import { useState, useRef, useEffect } from "react";
import Link from "next/link";
import CodeEditor from "./CodeEditor";
import { Clock, HardDrive, CheckCircle, XCircle, ArrowLeft } from "lucide-react";
import parse, { HTMLReactParserOptions } from "html-react-parser";
import katex from "katex";
import "katex/dist/katex.min.css";
import { QuestionDetail, SubmissionVerdict, fetchQuestionSubmissions } from "@/lib/api";

const parseOptions: HTMLReactParserOptions = {
  replace: (domNode: any) => {
    if (domNode.attribs && domNode.attribs.class && domNode.attribs.class.includes("math")) {
      const isBlock = domNode.attribs.class.includes("math-display");
      let mathText = "";
      if (domNode.children && domNode.children.length > 0) {
        mathText = domNode.children[0].data || "";
      }
      try {
        const html = katex.renderToString(mathText, {
          displayMode: isBlock,
          throwOnError: false,
        });
        return <span dangerouslySetInnerHTML={{ __html: html }} />;
      } catch (e) {
        return <span>{mathText}</span>;
      }
    }
  },
};

interface ProblemWorkspaceProps {
  question: QuestionDetail;
  submissions: SubmissionVerdict[];
}

export default function ProblemWorkspace({ question, submissions }: ProblemWorkspaceProps) {
  const [leftWidth, setLeftWidth] = useState(45);
  const [activeTab, setActiveTab] = useState<"description" | "submissions">("description");
  const [selectedHistoricalVerdict, setSelectedHistoricalVerdict] = useState<SubmissionVerdict | null>(null);
  
  const [localSubmissions, setLocalSubmissions] = useState<SubmissionVerdict[]>(submissions);
  
  useEffect(() => {
    fetchQuestionSubmissions(question.QuestionID).then((data) => {
      if (data && Array.isArray(data)) {
        setLocalSubmissions(data);
      }
    });
  }, [question.QuestionID]);
  
  const containerRef = useRef<HTMLDivElement>(null);
  const isDraggingRef = useRef(false);

  useEffect(() => {
    const handleMouseMove = (e: MouseEvent) => {
      if (!isDraggingRef.current || !containerRef.current) return;
      const containerRect = containerRef.current.getBoundingClientRect();
      let newWidthPercentage = ((e.clientX - containerRect.left) / containerRect.width) * 100;
      
      if (newWidthPercentage < 20) newWidthPercentage = 20;
      if (newWidthPercentage > 80) newWidthPercentage = 80;
      
      setLeftWidth(newWidthPercentage);
    };

    const handleMouseUp = () => {
      if (isDraggingRef.current) {
        isDraggingRef.current = false;
        document.body.style.cursor = "default";
        document.body.style.userSelect = "auto";
      }
    };

    document.addEventListener("mousemove", handleMouseMove);
    document.addEventListener("mouseup", handleMouseUp);
    return () => {
      document.removeEventListener("mousemove", handleMouseMove);
      document.removeEventListener("mouseup", handleMouseUp);
    };
  }, []);

  return (
    <div className="flex flex-col h-screen w-full bg-neutral-950 font-sans">
      
      <nav className="flex items-center justify-between px-6 h-14 bg-[#111118] border-b border-white/5 shrink-0 select-none">
         <Link href="/" className="flex items-center gap-2 px-3 py-1.5 rounded-lg text-xs font-bold uppercase tracking-widest text-neutral-400 hover:text-white hover:bg-white/5 transition-colors">
            <ArrowLeft className="w-4 h-4" />
            Problems
         </Link>
         <div className="text-white font-bold text-[13px] flex items-center gap-3">
            <span className="px-2 py-0.5 bg-white/5 rounded text-neutral-500 font-mono tracking-wider">#{question.QuestionID}</span> 
            {question.QuestionName}
         </div>
         <div className="w-[100px]" /> 
      </nav>

      
      <div ref={containerRef} className="flex flex-1 w-full overflow-hidden">
        
        <div 
          style={{ width: `${leftWidth}%` }} 
          className="flex flex-col h-full bg-[#0a0a0f] z-10"
        >
        
        <div className="flex bg-neutral-950 border-b border-white/10 shrink-0">
          <button
            onClick={() => setActiveTab("description")}
            className={`px-6 py-3 text-xs font-bold tracking-widest uppercase transition-all border-b-2 ${
              activeTab === "description" 
                ? "text-indigo-400 border-indigo-500 bg-indigo-500/5" 
                : "text-neutral-500 border-transparent hover:text-neutral-300 hover:bg-white/5"
            }`}
          >
            Description
          </button>
          <button
            onClick={() => setActiveTab("submissions")}
            className={`px-6 py-3 text-xs font-bold tracking-widest uppercase transition-all border-b-2 ${
              activeTab === "submissions" 
                ? "text-indigo-400 border-indigo-500 bg-indigo-500/5" 
                : "text-neutral-500 border-transparent hover:text-neutral-300 hover:bg-white/5"
            }`}
          >
            Submissions
          </button>
        </div>

        
        <div className="flex-1 overflow-y-auto custom-scrollbar p-6 sm:p-8">
          {activeTab === "description" && (
            <div className="animate-fade-in-up">
              <div className="flex items-center gap-3 mb-2">
                <span className="px-3 py-1 bg-indigo-500/10 text-indigo-400 border border-indigo-500/20 text-xs font-bold rounded-full uppercase tracking-wider">
                  {question.QuestionCategory}
                </span>
                <span className="text-neutral-500 text-sm font-medium">Problem #{question.QuestionID}</span>
              </div>
              
              <h1 className="text-3xl font-extrabold text-white mb-6 tracking-tight">
                {question.QuestionName}
              </h1>

              <div className="flex flex-wrap gap-4 mb-8">
                <div className="flex items-center gap-2 text-sm text-neutral-400 bg-white/5 border border-white/5 px-3 py-1.5 rounded-lg">
                  <Clock className="w-4 h-4 text-amber-400" />
                  Time Limit: {question.TimeConstraint}ms
                </div>
                <div className="flex items-center gap-2 text-sm text-neutral-400 bg-white/5 border border-white/5 px-3 py-1.5 rounded-lg">
                  <HardDrive className="w-4 h-4 text-cyan-400" />
                  Memory Limit: {question.MemConstraint}KB
                </div>
              </div>

              <div className="prose prose-invert max-w-none text-neutral-300">
                <h3 className="text-white font-semibold text-sm uppercase tracking-widest border-b border-white/10 pb-2 mb-4 text-neutral-400">Description</h3>
                <div className="math-rendered-content text-sm leading-relaxed">{parse(question.QuestionDescription, parseOptions)}</div>

                <h3 className="text-white font-semibold text-sm uppercase tracking-widest border-b border-white/10 pb-2 mt-8 mb-4 text-neutral-400">Input Specification</h3>
                <div className="math-rendered-content text-sm leading-relaxed">{parse(question.InputDescription, parseOptions)}</div>

                <h3 className="text-white font-semibold text-sm uppercase tracking-widest border-b border-white/10 pb-2 mt-8 mb-4 text-neutral-400">Output Specification</h3>
                <div className="math-rendered-content text-sm leading-relaxed">{parse(question.OutputDescription, parseOptions)}</div>

                {question.ConstraintsDescription && (
                  <>
                    <h3 className="text-white font-semibold text-sm uppercase tracking-widest border-b border-white/10 pb-2 mt-8 mb-4 text-neutral-400">Constraints</h3>
                    <div className="math-rendered-content text-sm leading-relaxed bg-neutral-900/50 p-4 rounded-xl border border-white/5">{parse(question.ConstraintsDescription, parseOptions)}</div>
                  </>
                )}

                {question.ExampleInputs && question.ExampleInputs.length > 0 && (
                  <div className="mt-8">
                    <h3 className="text-white font-semibold text-sm uppercase tracking-widest border-b border-white/10 pb-2 mb-4 text-neutral-400">Examples</h3>
                    {question.ExampleInputs.map((input, idx) => (
                      <div key={idx} className="mb-6 bg-[#16161f] rounded-xl overflow-hidden border border-white/5">
                        <div className="p-4 grid gap-4 grid-cols-1 sm:grid-cols-2">
                          <div>
                            <div className="text-[10px] text-neutral-500 mb-2 font-bold uppercase tracking-wider">Input</div>
                            <pre className="font-mono text-sm text-indigo-300 whitespace-pre-wrap">{input}</pre>
                          </div>
                          <div>
                            <div className="text-[10px] text-neutral-500 mb-2 font-bold uppercase tracking-wider">Output</div>
                            <pre className="font-mono text-sm text-emerald-300 whitespace-pre-wrap">{question.ExampleOutputs[idx]}</pre>
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            </div>
          )}

          {activeTab === "submissions" && (
            <div className="animate-fade-in-up">
              {localSubmissions && localSubmissions.length > 0 ? (
                <div className="w-full">
                  <table className="w-full text-left border-collapse">
                    <thead>
                      <tr className="border-b border-white/10">
                        <th className="px-4 py-3 text-[10px] font-bold text-neutral-500 uppercase tracking-widest">Status</th>
                        <th className="px-4 py-3 text-[10px] font-bold text-neutral-500 uppercase tracking-widest">Memory</th>
                        <th className="px-4 py-3 text-[10px] font-bold text-neutral-500 uppercase tracking-widest">Time</th>
                        <th className="px-4 py-3 text-[10px] font-bold text-neutral-500 uppercase tracking-widest">Date</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-white/5">
                      {[...localSubmissions].reverse().map(sub => (
                        <tr 
                          key={sub.Submission_id} 
                          onClick={() => setSelectedHistoricalVerdict(sub)}
                          className={`hover:bg-white/5 transition-colors group cursor-pointer ${selectedHistoricalVerdict?.Submission_id === sub.Submission_id ? 'bg-white/10' : ''}`}
                        >
                          <td className="px-4 py-3">
                            <span className={`inline-flex items-center gap-1.5 font-bold text-xs ${sub.Verdict === "Accepted" ? "text-emerald-400" : "text-red-400"}`}>
                                {sub.Verdict === "Accepted" ? <CheckCircle className="w-3.5 h-3.5" /> : <XCircle className="w-3.5 h-3.5" />}
                                {sub.Verdict}
                            </span>
                          </td>
                          <td className="px-4 py-3 text-xs font-mono text-neutral-300">{sub.Mem_usage > 0 ? `${sub.Mem_usage} KB` : "0 KB"}</td>
                          <td className="px-4 py-3 text-xs font-mono text-neutral-300">{sub.Time_ms > 0 ? `${sub.Time_ms} ms` : "0 ms"}</td>
                          <td className="px-4 py-3 text-xs text-neutral-500">{new Date(sub.SubmissionTime).toLocaleDateString()}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              ) : (
                <div className="py-20 flex flex-col items-center justify-center text-center">
                  <div className="w-16 h-16 bg-white/5 rounded-full flex items-center justify-center mb-4">
                    <CheckCircle className="w-6 h-6 text-neutral-600" />
                  </div>
                  <p className="text-neutral-400 font-medium">No submissions yet.</p>
                  <p className="text-neutral-500 text-sm mt-1">Submit your first solution using the editor.</p>
                </div>
              )}
            </div>
          )}
        </div>
      </div>

      
      <div 
        onMouseDown={() => {
          isDraggingRef.current = true;
          document.body.style.cursor = "col-resize";
          document.body.style.userSelect = "none";
        }}
        className="w-1.5 bg-neutral-900 border-x border-white/5 hover:bg-indigo-500 cursor-col-resize z-50 transition-colors active:bg-indigo-400 grab-handle"
      />

      
      <div style={{ width: `${100 - leftWidth}%` }} className="flex flex-col h-full z-10">
        <CodeEditor 
          questionId={question.QuestionID} 
          historicalVerdict={selectedHistoricalVerdict} 
          onNewVerdict={(v) => setLocalSubmissions(prev => [...prev, v])}
        />
      </div>
      
      </div> 

    </div>
  );
}

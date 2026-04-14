"use client";

import { useRef, useState, useEffect } from "react";
import Editor from "@monaco-editor/react";
import { submitCode, fetchVerdict, SubmissionVerdict } from "@/lib/api";
import { Play, Loader2, CheckCircle, XCircle, ChevronUp, Terminal } from "lucide-react";

interface CodeEditorProps {
  questionId: number;
  historicalVerdict?: SubmissionVerdict | null;
  onNewVerdict?: (v: SubmissionVerdict) => void;
}

const RUNTIMES = [
  { id: "c++23", name: "C++ 23", language: "cpp", template: `#include <iostream>\nusing namespace std;\n\nint main() {\n    // solve here\n    return 0;\n}` },
  { id: "c++20", name: "C++ 20", language: "cpp", template: `#include <iostream>\nusing namespace std;\n\nint main() {\n    // solve here\n    return 0;\n}` },
  { id: "c++17", name: "C++ 17", language: "cpp", template: `#include <iostream>\nusing namespace std;\n\nint main() {\n    // solve here\n    return 0;\n}` },
  { id: "gcc-c17", name: "C 17", language: "c", template: `#include <stdio.h>\n\nint main() {\n    // solve here\n    return 0;\n}` },
  { id: "python3", name: "Python 3", language: "python", template: `import sys\n\ndef solve():\n    pass\n\nif __name__ == "__main__":\n    solve()` },
  { id: "node-25", name: "JavaScript", language: "javascript", template: `const fs = require('fs');\nconst input = fs.readFileSync(0, 'utf8');\n\nconsole.log(input);` },
];

export default function CodeEditor({ questionId, historicalVerdict, onNewVerdict }: CodeEditorProps) {
  const [runtime, setRuntime] = useState(RUNTIMES[0].id);
  const [language, setLanguage] = useState(RUNTIMES[0].language);
  const [code, setCode] = useState(RUNTIMES[0].template);
  
  const [submitting, setSubmitting] = useState(false);
  const [verdict, setVerdict] = useState<SubmissionVerdict | null>(null);

  
  const [consoleHeight, setConsoleHeight] = useState(44); // 44 is collapsed height
  const [isConsoleExpanded, setIsConsoleExpanded] = useState(false);
  
  const containerRef = useRef<HTMLDivElement>(null);
  const isDraggingRef = useRef(false);

  
  const pollRef = useRef<ReturnType<typeof setInterval> | null>(null);
  const timeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  useEffect(() => {
    if (historicalVerdict) {
      setVerdict(historicalVerdict);
      if (consoleHeight < 250) {
        setConsoleHeight(250);
        setIsConsoleExpanded(true);
      }
    }
  }, [historicalVerdict]);

  useEffect(() => {
    const handleMouseMove = (e: MouseEvent) => {
      if (!isDraggingRef.current || !containerRef.current) return;
      const containerRect = containerRef.current.getBoundingClientRect();
      let newHeight = containerRect.bottom - e.clientY;
      
      if (newHeight >= 44 && newHeight < containerRect.height - 100) {
        setConsoleHeight(newHeight);
        if (newHeight > 60) {
          setIsConsoleExpanded(true);
        } else {
          setIsConsoleExpanded(false);
        }
      }
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

  const stopPolling = () => {
    if (pollRef.current) clearInterval(pollRef.current);
    if (timeoutRef.current) clearTimeout(timeoutRef.current);
  };

  const handleRuntimeChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const selectedId = e.target.value;
    setRuntime(selectedId);
    const linked = RUNTIMES.find((r) => r.id === selectedId);
    if (linked) {
      setLanguage(linked.language);
      // reset template
      const oldRuntime = RUNTIMES.find((r) => r.language === language);
      if (code === oldRuntime?.template || code.trim() === "") {
        setCode(linked.template);
      }
    }
  };

  const toggleConsole = () => {
    if (isConsoleExpanded) {
      setConsoleHeight(44);
      setIsConsoleExpanded(false);
    } else {
      setConsoleHeight(250);
      setIsConsoleExpanded(true);
    }
  };

  const handleSimulateRun = () => {
    setVerdict({ Verdict: "Simulated Run" } as SubmissionVerdict);
    if (!isConsoleExpanded) toggleConsole();
  };

  const handleSubmit = async () => {
    if (!code.trim() || submitting) return;

    setSubmitting(true);
    setVerdict(null);
    stopPolling();
    
    
    if (consoleHeight < 200) {
      setConsoleHeight(250);
      setIsConsoleExpanded(true);
    }

    try {
      const subId = await submitCode({ QuestionId: questionId, runtime, code });
      if (!subId) throw new Error("No submission ID returned from server");

      pollRef.current = setInterval(async () => {
        const res = await fetchVerdict(subId);
        if (res && res.Verdict) {
          if (res.Verdict !== "queued" && res.Verdict !== "Running") {
            stopPolling();
            setVerdict(res);
            setSubmitting(false);
            if (onNewVerdict) onNewVerdict(res);
          } else {
            setVerdict(res); // update UI to show 'queued' or 'Running'
          }
        }
      }, 1500);

      // 30s timeout
      timeoutRef.current = setTimeout(() => {
        stopPolling();
        setSubmitting(false);
        setVerdict({ Verdict: "Timed out waiting for judge" } as SubmissionVerdict);
      }, 30000);

    } catch (err) {
      console.error(err);
      stopPolling();
      setSubmitting(false);
      setVerdict({ Verdict: "Submission error — check console" } as SubmissionVerdict);
    }
  };

  return (
    <div ref={containerRef} className="flex flex-col h-full w-full bg-[#1e1e1e] relative">
      
      <div className="flex items-center justify-between p-[10px] px-4 bg-[#111118] border-b border-white/5">
        <select
          value={runtime}
          onChange={handleRuntimeChange}
          className="bg-black/50 border border-white/10 text-white text-xs rounded-md focus:outline-none focus:ring-1 focus:ring-indigo-500 p-1.5 font-mono"
        >
          {RUNTIMES.map((r) => (
            <option key={r.id} value={r.id}>{r.name}</option>
          ))}
        </select>

        <div className="flex items-center gap-2">
          
          <button
            onClick={handleSimulateRun}
            disabled={submitting}
            className="flex items-center gap-2 bg-white/5 hover:bg-white/10 text-neutral-300 transition-colors text-xs px-4 py-1.5 rounded-md font-bold disabled:opacity-50"
          >
            Run Code
          </button>
          
          
          <button
            onClick={handleSubmit}
            disabled={submitting}
            className="flex items-center gap-2 bg-indigo-600 hover:bg-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed text-white text-xs px-4 py-1.5 rounded-md font-bold transition-all"
          >
            {submitting ? <Loader2 className="w-3.5 h-3.5 animate-spin" /> : <Play className="w-3.5 h-3.5" />}
            {submitting ? "Evaluating" : "Submit"}
          </button>
        </div>
      </div>

      
      <div className="flex-1 w-full relative min-h-0">
        <Editor
          height="100%"
          language={language}
          theme="vs-dark"
          value={code}
          onChange={(val) => setCode(val ?? "")}
          options={{
            minimap: { enabled: false },
            fontSize: 14,
            fontFamily: "'JetBrains Mono', 'Fira Code', monospace",
            padding: { top: 16 },
            scrollBeyondLastLine: false,
            cursorSmoothCaretAnimation: "on",
            roundedSelection: true
          }}
        />
      </div>

      
      <div 
        style={{ height: `${consoleHeight}px` }} 
        className="w-full bg-[#111118] border-t border-white/10 flex flex-col z-50 shrink-0 transition-[height] duration-75 relative"
      >
        
        <div 
          onMouseDown={() => {
            isDraggingRef.current = true;
            document.body.style.cursor = "ns-resize";
            document.body.style.userSelect = "none";
          }}
          className="h-1 w-full absolute top-0 -translate-y-1/2 cursor-ns-resize hover:bg-indigo-500 z-50 transition-colors"
        />

        
        <div 
          onClick={toggleConsole}
          className="h-[44px] min-h-[44px] px-5 flex items-center justify-between cursor-pointer hover:bg-white/5 transition-colors select-none"
        >
          <span className="flex items-center gap-2 text-xs font-bold text-neutral-400 uppercase tracking-widest">
            <ChevronUp className={`w-4 h-4 transition-transform duration-300 ${isConsoleExpanded ? "rotate-180" : ""}`} />
            Console
          </span>
          <span className="font-mono text-[10px] uppercase font-bold text-neutral-500 flex items-center gap-2">
            {submitting && <Loader2 className="w-3 h-3 text-indigo-400 animate-spin" />}
            {verdict?.Verdict === "Accepted" && <span className="text-emerald-400">Accepted</span>}
            {verdict?.Verdict && verdict.Verdict !== "Accepted" && verdict.Verdict !== "queued" && verdict.Verdict !== "Running" && <span className="text-red-400">{verdict.Verdict}</span>}
            {(verdict?.Verdict === "queued" || verdict?.Verdict === "Running") && <span className="text-indigo-400">{verdict.Verdict}</span>}
          </span>
        </div>

        
        {isConsoleExpanded && (
          <div className="flex-1 overflow-y-auto px-6 py-4 custom-scrollbar">
            {submitting || verdict?.Verdict === "queued" || verdict?.Verdict === "Running" ? (
              <div className="flex flex-col items-center justify-center h-full text-neutral-500 text-sm gap-4">
                <Loader2 className="w-8 h-8 animate-spin text-indigo-500/50" />
                {verdict?.Verdict === "Running" ? "Executing on judge workers..." : "In judging queue..."}
              </div>
            ) : verdict ? (
              <div className={`p-5 rounded-xl border ${
                verdict.Verdict === "Accepted"
                  ? "bg-emerald-500/5 border-emerald-500/20"
                  : verdict.Verdict === "Simulated Run" 
                  ? "bg-white/5 border-white/10" 
                  : "bg-red-500/5 border-red-500/20"
              }`}>
                <div className="flex items-center gap-2 font-bold mb-4">
                  {verdict.Verdict === "Accepted" ? (
                    <CheckCircle className="w-5 h-5 text-emerald-400" />
                  ) : verdict.Verdict === "Simulated Run" ? (
                    <Terminal className="w-5 h-5 text-neutral-400" />
                  ) : (
                    <XCircle className="w-5 h-5 text-red-400" />
                  )}
                  <span className={verdict.Verdict === "Accepted" ? "text-emerald-400" : verdict.Verdict === "Simulated Run" ? "text-neutral-300" : "text-red-400"}>
                    {verdict.Verdict}
                  </span>
                </div>
                
                {verdict.Verdict !== "Simulated Run" && (
                  <div className="text-xs font-mono opacity-80 flex gap-6 text-neutral-400 bg-black/40 px-3 py-2 rounded border border-white/5 inline-flex mb-4">
                    {verdict.Time_ms > 0 && <span>Time: <strong className="text-neutral-200">{verdict.Time_ms}ms</strong></span>}
                    {verdict.Mem_usage > 0 && <span>Memory: <strong className="text-neutral-200">{verdict.Mem_usage}KB</strong></span>}
                  </div>
                )}

                {verdict.Stderr && (
                  <div className="mt-2">
                    <span className="text-[10px] font-bold text-red-500/80 uppercase tracking-widest bg-red-500/10 px-2 py-0.5 rounded">Error Stack Trace</span>
                    <pre className="mt-2 p-3 bg-black/60 rounded-lg text-xs overflow-x-auto border border-red-500/20 text-red-300 font-mono">
                      {verdict.Stderr}
                    </pre>
                  </div>
                )}
                
                {verdict.WA_Test_case > 0 && (
                  <div className="mt-4">
                    <span className="text-sm font-semibold text-neutral-300 border-b border-white/10 pb-1 flex w-full">
                      Wrong Answer on Test Case {verdict.WA_Test_case}
                    </span>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-3 mt-3">
                      <div className="p-3 bg-black/50 rounded-lg border border-red-500/10">
                        <div className="text-[10px] font-bold text-neutral-500 uppercase tracking-wider mb-2">Your Actual Output</div>
                        <pre className="text-red-400 overflow-x-auto text-xs font-mono">{verdict.WA_Stdout}</pre>
                      </div>
                      <div className="p-3 bg-black/50 rounded-lg border border-emerald-500/10">
                        <div className="text-[10px] font-bold text-neutral-500 uppercase tracking-wider mb-2">Expected Correct Output</div>
                        <pre className="text-emerald-400 overflow-x-auto text-xs font-mono">{verdict.Required_Stdout}</pre>
                      </div>
                    </div>
                  </div>
                )}
              </div>
            ) : (
              <div className="flex flex-col items-center justify-center h-full text-neutral-600 text-sm italic">
                Hit Run or Submit to see the output here.
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
}

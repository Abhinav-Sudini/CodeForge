"use client";

import { useState, useEffect } from "react";
import Editor from "@monaco-editor/react";
import { submitCode, fetchVerdict, SubmissionVerdict } from "@/lib/api";
import { Play, Loader2, CheckCircle, XCircle, AlertTriangle } from "lucide-react";

interface CodeEditorProps {
  questionId: number;
}

const RUNTIMES = [
  { id: "python3", name: "Python 3", language: "python" },
  { id: "c++17", name: "C++ 17", language: "cpp" },
  { id: "c++20", name: "C++ 20", language: "cpp" },
  { id: "node-25", name: "Node.js", language: "javascript" },
];

export default function CodeEditor({ questionId }: CodeEditorProps) {
  const [runtime, setRuntime] = useState(RUNTIMES[0].id);
  const [language, setLanguage] = useState(RUNTIMES[0].language);
  const [code, setCode] = useState("// Write your optimal algorithm here...\n");
  const [submitting, setSubmitting] = useState(false);
  const [verdict, setVerdict] = useState<SubmissionVerdict | null>(null);

  const handleRuntimeChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const selectedId = e.target.value;
    setRuntime(selectedId);
    const linkedLang = RUNTIMES.find((r) => r.id === selectedId)?.language || "plaintext";
    setLanguage(linkedLang);
  };

  const handleRun = async () => {
    if (!code.trim()) return;
    setSubmitting(true);
    setVerdict(null);

    try {
      const subId = await submitCode({ QuestionId: questionId, runtime, code });
      if (!subId) throw new Error("Failed to submit");

      // Polling mechanism
      const poll = setInterval(async () => {
        const res = await fetchVerdict(subId);
        if (res && res.Verdict && res.Verdict !== "PENDING") {
          clearInterval(poll);
          setVerdict(res);
          setSubmitting(false);
        }
      }, 1500); // Check every 1.5 seconds

      // Timeout safety (30 seconds)
      setTimeout(() => {
        clearInterval(poll);
        if (submitting) {
            setSubmitting(false);
            setVerdict({ Verdict: "TIMEOUT_UNKNOWN" } as SubmissionVerdict);
        }
      }, 30000);

    } catch (err) {
      console.error(err);
      setSubmitting(false);
      setVerdict({ Verdict: "ERROR" } as SubmissionVerdict);
    }
  };

  return (
    <div className="flex flex-col h-full bg-neutral-900 border-l border-white/10 relative">
      <div className="flex items-center justify-between p-3 border-b border-white/5 bg-neutral-950/50">
        <select 
          value={runtime}
          onChange={handleRuntimeChange}
          className="bg-neutral-800 border border-white/10 text-white text-sm rounded-lg focus:ring-indigo-500 focus:border-indigo-500 block p-2"
        >
          {RUNTIMES.map((r) => (
            <option key={r.id} value={r.id}>{r.name}</option>
          ))}
        </select>
        
        <button
          onClick={handleRun}
          disabled={submitting}
          className="flex items-center gap-2 bg-indigo-600 hover:bg-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed text-white text-sm px-4 py-2 rounded-lg font-medium transition-all"
        >
          {submitting ? <Loader2 className="w-4 h-4 animate-spin" /> : <Play className="w-4 h-4" />}
          {submitting ? "Evaluating..." : "Submit Code"}
        </button>
      </div>

      <div className="flex-1 w-full relative">
        <Editor
          height="100%"
          language={language}
          theme="vs-dark"
          value={code}
          onChange={(val) => setCode(val || "")}
          options={{
            minimap: { enabled: false },
            fontSize: 14,
            fontFamily: "'JetBrains Mono', 'Fira Code', monospace",
            padding: { top: 16 },
            scrollBeyondLastLine: false,
          }}
        />
      </div>

      {/* Verdict Panel */}
      {verdict && (
        <div className={`p-4 border-t ${
          verdict.Verdict === "Accepted" 
            ? "bg-emerald-950/30 border-emerald-500/30 text-emerald-400" 
            : "bg-red-950/30 border-red-500/30 text-red-400"
        }`}>
          <div className="flex items-center gap-2 font-bold mb-2">
            {verdict.Verdict === "Accepted" ? <CheckCircle className="w-5 h-5" /> : <XCircle className="w-5 h-5" />}
            {verdict.Verdict}
          </div>
          <div className="text-sm font-mono opacity-80 flex gap-4">
            {verdict.Time_ms > 0 && <span>Time: {verdict.Time_ms}ms</span>}
            {verdict.Mem_usage > 0 && <span>Memory: {verdict.Mem_usage}KB</span>}
          </div>
          {verdict.Stderr && (
            <pre className="mt-4 p-3 bg-red-950/50 rounded-lg text-xs overflow-x-auto border border-red-500/20">
              {verdict.Stderr}
            </pre>
          )}
          {verdict.WA_Test_case > 0 && (
            <div className="mt-4 text-xs">
              <span className="font-semibold text-neutral-300">Failed on Test Case {verdict.WA_Test_case}</span>
              <div className="grid grid-cols-2 gap-2 mt-2">
                <div className="p-2 bg-black/50 rounded">
                  <div className="text-neutral-500 mb-1">Received Output:</div>
                  <pre className="text-red-300 overflow-x-auto">{verdict.WA_Stdout}</pre>
                </div>
                <div className="p-2 bg-black/50 rounded">
                  <div className="text-neutral-500 mb-1">Expected Output:</div>
                  <pre className="text-emerald-300 overflow-x-auto">{verdict.Required_Stdout}</pre>
                </div>
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
}

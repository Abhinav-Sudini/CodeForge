import { fetchQuestionDetail, fetchQuestionSubmissions } from "@/lib/api";
import CodeEditor from "@/components/CodeEditor";
import { notFound } from "next/navigation";
import { Clock, HardDrive, CheckCircle, XCircle } from "lucide-react";
import parse, { HTMLReactParserOptions } from 'html-react-parser';
import katex from 'katex';
import 'katex/dist/katex.min.css';

const parseOptions: HTMLReactParserOptions = {
  replace: (domNode: any) => {
    if (domNode.attribs && domNode.attribs.class && domNode.attribs.class.includes('math')) {
      const isBlock = domNode.attribs.class.includes('math-display');
      
      let mathText = '';
      if (domNode.children && domNode.children.length > 0) {
         mathText = domNode.children[0].data || '';
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
  }
};

export const dynamic = 'force-dynamic';

export default async function ProblemPage({ params }: { params: { id: string } }) {
  const { id } = await params;
  const questionId = parseInt(id, 10);
  
  if (isNaN(questionId)) return notFound();

  const [question, submissions] = await Promise.all([
    fetchQuestionDetail(questionId),
    fetchQuestionSubmissions(questionId)
  ]);

  if (!question) return notFound();

  return (
    <div className="flex flex-col md:flex-row h-[calc(100vh-64px)] w-full overflow-hidden bg-neutral-950">
      
      {/* Left Pane: Problem Description */}
      <div className="w-full md:w-1/2 lg:w-[45%] h-full overflow-y-auto custom-scrollbar bg-neutral-950 p-6 sm:p-8">
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
          <h3 className="text-white font-semibold text-lg border-b border-white/10 pb-2 mb-4">Description</h3>
          <div className="math-rendered-content">{parse(question.QuestionDescription, parseOptions)}</div>

          <h3 className="text-white font-semibold text-lg border-b border-white/10 pb-2 mt-8 mb-4">Input</h3>
          <div className="math-rendered-content">{parse(question.InputDescription, parseOptions)}</div>

          <h3 className="text-white font-semibold text-lg border-b border-white/10 pb-2 mt-8 mb-4">Output</h3>
          <div className="math-rendered-content">{parse(question.OutputDescription, parseOptions)}</div>

          {question.ConstraintsDescription && (
            <>
              <h3 className="text-white font-semibold text-lg border-b border-white/10 pb-2 mt-8 mb-4">Constraints</h3>
              <div className="math-rendered-content bg-neutral-900 p-4 rounded-xl border border-white/5">{parse(question.ConstraintsDescription, parseOptions)}</div>
            </>
          )}

          {question.ExampleInputs && question.ExampleInputs.length > 0 && (
            <div className="mt-8">
              <h3 className="text-white font-semibold text-lg border-b border-white/10 pb-2 mb-4">Examples</h3>
              {question.ExampleInputs.map((input, idx) => (
                <div key={idx} className="mb-6 bg-neutral-900 rounded-xl overflow-hidden border border-white/5">
                  <div className="px-4 py-2 bg-black/40 border-b border-white/5 font-semibold text-sm text-neutral-400">Example {idx + 1}</div>
                  <div className="p-4 grid gap-4 grid-cols-1 sm:grid-cols-2">
                    <div>
                      <div className="text-xs text-neutral-500 mb-1 font-semibold uppercase tracking-wider">Input</div>
                      <pre className="font-mono text-sm text-neutral-200 whitespace-pre-wrap">{input}</pre>
                    </div>
                    <div>
                      <div className="text-xs text-neutral-500 mb-1 font-semibold uppercase tracking-wider">Output</div>
                      <pre className="font-mono text-sm text-neutral-200 whitespace-pre-wrap">{question.ExampleOutputs[idx]}</pre>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Previous Submissions */}
        <div className="mt-12 pt-8 border-t border-white/10">
          <h3 className="text-white font-semibold text-lg mb-4">Your Recent Submissions</h3>
          {submissions && submissions.length > 0 ? (
            <div className="flex flex-col gap-2">
              {submissions.map(sub => (
                <div key={sub.Submission_id} className="flex flex-wrap items-center justify-between p-3 rounded-lg bg-neutral-900/50 border border-white/5">
                  <div className="flex items-center gap-3">
                    {sub.Verdict === "Accepted" ? <CheckCircle className="w-4 h-4 text-emerald-400" /> : <XCircle className="w-4 h-4 text-red-400" />}
                    <span className="text-sm font-medium text-white">{sub.Verdict}</span>
                  </div>
                  <div className="text-xs text-neutral-500 font-mono">
                    {sub.Time_ms > 0 ? `${sub.Time_ms}ms` : "-"} | {sub.Mem_usage > 0 ? `${sub.Mem_usage}KB` : "-"}
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <p className="text-neutral-500 text-sm italic">You haven't submitted any solutions for this problem yet.</p>
          )}
        </div>
      </div>

      {/* Right Pane: Code Editor */}
      <div className="w-full md:w-1/2 lg:w-[55%] h-[50vh] md:h-full">
        <CodeEditor questionId={questionId} />
      </div>

    </div>
  );
}

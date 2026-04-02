export const API_BASE = "http://localhost:7000/api";

export interface QuestionMinimal {
  QuestionId: number;
  QuestionName: string;
  QuestionCategory: string;
}

export interface QuestionDetail extends QuestionMinimal {
  QuestionDescription: string;
  InputDescription: string;
  OutputDescription: string;
  ConstraintsDescription: string;
  TimeConstraint: number;
  MemConstraint: number;
  ExampleInputs: string[];
  ExampleOutputs: string[];
}

export interface SubmissionVerdict {
  Submission_id: number;
  QuestionId: number;
  Verdict: string;
  Mem_usage: number;
  Time_ms: number;
  WA_Test_case: number;
  WA_Stdin: string;
  WA_Stdout: string;
  Required_Stdout: string;
  Stderr: string;
  SubmissionTime: string;
}

export interface CodeSubmissionContext {
  QuestionId: number;
  runtime: string;
  code: string;
}

// -- API Helpers --

export async function fetchQuestions(): Promise<QuestionMinimal[]> {
  try {
    const res = await fetch(`${API_BASE}/question/`, { next: { revalidate: 10 } });
    if (!res.ok) throw new Error("Failed to fetch questions");
    const data = await res.json();
    return data.Questions || [];
  } catch (err) {
    console.error(err);
    return [];
  }
}

export async function fetchQuestionDetail(id: number): Promise<QuestionDetail | null> {
  try {
    const res = await fetch(`${API_BASE}/question/${id}/`, { cache: 'no-store' });
    if (!res.ok) throw new Error("Failed to fetch question detail");
    const data = await res.json();
    return data;
  } catch (err) {
    console.error(err);
    return null;
  }
}

export async function submitCode(payload: CodeSubmissionContext): Promise<number | null> {
  try {
    const res = await fetch(`${API_BASE}/judge/`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });
    if (!res.ok) {
      const errText = await res.text();
      console.error("Backend error response:", errText);
      throw new Error(`Submission rejected: ${errText}`);
    }
    const data = await res.json();
    return data.Submission_id;
  } catch (err) {
    console.error(err);
    return null;
  }
}

export async function fetchVerdict(submissionId: number): Promise<SubmissionVerdict | null> {
  try {
    const res = await fetch(`${API_BASE}/submissions/${submissionId}/`);
    if (!res.ok) throw new Error("Verdict fetch failed");
    return await res.json();
  } catch (err) {
    console.error(err);
    return null;
  }
}

export async function fetchQuestionSubmissions(qId: number): Promise<SubmissionVerdict[]> {
  try {
    const res = await fetch(`${API_BASE}/question/submissions/${qId}/`, { cache: 'no-store' });
    if (!res.ok) return [];
    const data = await res.json();
    return data.AllSubmissions || [];
  } catch (err) {
    console.error(err);
    return [];
  }
}

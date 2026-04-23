export const API_BASE = "https://azure.abhi.dedyn.io/cfbackend/api";

export interface QuestionMinimal {
  QuestionId: number;
  QuestionName: string;
  QuestionCategory: string;
}

export interface QuestionDetail {
  QuestionID: number;
  QuestionName: string;
  QuestionCategory: string;
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
  Submitted_code?: string;
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

// Auth

export async function registerUser(
  email: string,
  password: string,
): Promise<{ ok: boolean; message: string }> {
  try {
    const res = await fetch(`/api/auth/register`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ UserName: email, Password: password }),
    });
    if (res.ok) {
      return { ok: true, message: "user created" };
    }
    return { ok: false, message: "Registration failed. An account with this email might already exist." };
  } catch (err) {
    console.error(err);
    return { ok: false, message: "Network error. Please try again." };
  }
}

export async function loginUser(
  email: string,
  password: string,
): Promise<{ ok: boolean; message: string }> {
  try {
    const res = await fetch(`/api/auth/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ UserName: email, Password: password }),
    });
    if (res.ok) {
      return { ok: true, message: "" };
    }
    return { ok: false, message: "Invalid email or password. Please try again." };
  } catch (err) {
    console.error(err);
    return { ok: false, message: "Network error. Please try again." };
  }
}

// API Helpers

export async function fetchQuestions(): Promise<QuestionMinimal[]> {
  try {
    const res = await fetch(`${API_BASE}/question/`, {
      next: { revalidate: 10 },
    });
    if (!res.ok) throw new Error("Failed to fetch questions");
    const data = await res.json();
    return data.Questions || [];
  } catch (err) {
    console.error(err);
    return [];
  }
}

export async function fetchQuestionDetail(
  id: number,
): Promise<QuestionDetail | null> {
  try {
    const res = await fetch(`${API_BASE}/question/${id}/`, {
      cache: "no-store",
    });
    if (!res.ok) throw new Error("Failed to fetch question detail");
    const data = await res.json();
    return data;
  } catch (err) {
    console.error(err);
    return null;
  }
}

export async function submitCode(
  payload: CodeSubmissionContext,
): Promise<number> {
  const res = await fetch(`/api/judge`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });
  if (!res.ok) {
    const errText = await res.text();
    console.error("Backend error response:", errText);
    throw new Error(`${res.status} - ${errText || "Unknown Backend Error"}`);
  }
  const data = await res.json();
  if (!data.Submission_id) throw new Error("No submission ID returned from server.");
  return data.Submission_id;
}

export async function fetchVerdict(
  submissionId: number,
): Promise<SubmissionVerdict | null> {
  try {
    const res = await fetch(`/api/submissions/${submissionId}`);
    if (!res.ok) throw new Error("Verdict fetch failed");
    return await res.json();
  } catch (err) {
    console.error(err);
    return null;
  }
}

export async function fetchQuestionSubmissions(
  qId: number,
): Promise<SubmissionVerdict[]> {
  try {
    const res = await fetch(`/api/question-submissions/${qId}/`, {
      cache: "no-store",
    });
    if (!res.ok) return [];
    const data = await res.json();
    return data.AllSubmissions || [];
  } catch (err) {
    console.error(err);
    return [];
  }
}

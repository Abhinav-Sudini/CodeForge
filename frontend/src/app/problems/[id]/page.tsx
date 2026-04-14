import { fetchQuestionDetail, fetchQuestionSubmissions } from "@/lib/api";
import { notFound } from "next/navigation";
import ProblemWorkspace from "@/components/ProblemWorkspace";

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

  return <ProblemWorkspace question={question} submissions={submissions || []} />;
}

import { fetchQuestions } from "@/lib/api";
import ProblemsDashboard from "@/components/ProblemsDashboard";

export const dynamic = 'force-dynamic'; // We want the questions to be fetched fresh!

export default async function Home() {
  const questions = await fetchQuestions();

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      {/* Problems Interactive Dashboard (Hero banner is now inside to access state) */}
      <ProblemsDashboard questions={questions} />
    </div>
  );
}

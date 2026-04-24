export default function Loading() {
  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      {/* skeleton hero */}
      <div className="mb-12 rounded-3xl bg-neutral-900 border border-white/5 p-8 sm:p-12 h-[350px] relative overflow-hidden">
        <div className="absolute inset-0 -translate-x-full animate-[shimmer_2s_infinite] bg-gradient-to-r from-transparent via-white/5 to-transparent" />
        <div className="w-1/4 h-4 bg-white/10 rounded mb-6" />
        <div className="w-1/2 h-12 bg-white/10 rounded mb-4" />
        <div className="w-3/4 h-6 bg-white/10 rounded mb-8" />
        <div className="flex gap-4">
          <div className="w-32 h-16 bg-white/10 rounded-2xl" />
          <div className="w-32 h-16 bg-white/10 rounded-2xl" />
          <div className="w-32 h-16 bg-white/10 rounded-2xl" />
        </div>
      </div>

      {/* skeleton dashboard */}
      <div className="w-full animate-fade-in-up">
        {/* toolbar skeleton */}
        <div className="flex flex-col md:flex-row gap-4 mb-6">
          <div className="h-12 flex-1 bg-white/5 rounded-xl border border-white/5 overflow-hidden relative">
            <div className="absolute inset-0 -translate-x-full animate-[shimmer_2s_infinite] bg-gradient-to-r from-transparent via-white/5 to-transparent" />
          </div>
          <div className="h-12 w-32 bg-white/5 rounded-xl border border-white/5" />
          <div className="h-12 w-32 bg-white/5 rounded-xl border border-white/5" />
        </div>

        {/* categories skeleton */}
        <div className="flex gap-2 mb-8">
          <div className="h-8 w-20 bg-white/10 rounded-full" />
          <div className="h-8 w-24 bg-white/5 rounded-full" />
          <div className="h-8 w-32 bg-white/5 rounded-full" />
        </div>

        {/* table skeleton */}
        <div className="bg-neutral-900/50 border border-white/5 rounded-2xl overflow-hidden">
          <table className="w-full text-left">
            <thead>
              <tr className="border-b border-white/10">
                <th className="px-6 py-4"><div className="h-4 w-10 bg-white/10 rounded" /></th>
                <th className="px-6 py-4"><div className="h-4 w-32 bg-white/10 rounded" /></th>
                <th className="px-6 py-4"><div className="h-4 w-20 bg-white/10 rounded" /></th>
                <th className="px-6 py-4"></th>
              </tr>
            </thead>
            <tbody className="divide-y divide-white/5">
              {[...Array(5)].map((_, i) => (
                <tr key={i}>
                  <td className="px-6 py-5"><div className="h-4 w-12 bg-white/5 rounded overflow-hidden relative"><div className="absolute inset-0 -translate-x-full animate-[shimmer_2s_infinite] bg-gradient-to-r from-transparent via-white/10 to-transparent" /></div></td>
                  <td className="px-6 py-5"><div className="h-5 w-64 bg-white/5 rounded overflow-hidden relative"><div className="absolute inset-0 -translate-x-full animate-[shimmer_2s_infinite] bg-gradient-to-r from-transparent via-white/10 to-transparent" /></div></td>
                  <td className="px-6 py-5"><div className="h-6 w-24 bg-white/5 rounded-md overflow-hidden relative"><div className="absolute inset-0 -translate-x-full animate-[shimmer_2s_infinite] bg-gradient-to-r from-transparent via-white/10 to-transparent" /></div></td>
                  <td className="px-6 py-5"><div className="h-8 w-24 bg-white/5 rounded-lg ml-auto" /></td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}

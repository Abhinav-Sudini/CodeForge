export default function ProblemLoadingSkeleton() {
  return (
    <div className="flex flex-col h-screen w-full bg-neutral-950 font-sans">
      
      <nav className="flex items-center justify-between px-6 h-14 bg-[#111118] border-b border-white/5 shrink-0 select-none">
        <div className="w-24 h-6 bg-white/5 rounded-md animate-pulse"></div>
        <div className="w-64 h-6 bg-white/5 rounded-md animate-pulse"></div>
        <div className="w-[100px]" />
      </nav>

      
      <div className="flex flex-1 w-full overflow-hidden">
        
        <div style={{ width: "45%" }} className="flex flex-col h-full bg-[#0a0a0f]">
          
          <div className="flex bg-neutral-950 border-b border-white/10 shrink-0 px-6 py-3 gap-6">
            <div className="w-20 h-4 bg-white/5 rounded animate-pulse"></div>
            <div className="w-24 h-4 bg-white/5 rounded animate-pulse"></div>
          </div>

          <div className="p-8 flex flex-col gap-6">
            <div className="w-16 h-6 bg-indigo-500/10 rounded-full animate-pulse"></div>
            <div className="w-3/4 h-10 bg-white/5 rounded-md animate-pulse"></div>
            
            <div className="flex gap-4 mb-4">
              <div className="w-32 h-8 bg-white/5 rounded-lg animate-pulse"></div>
              <div className="w-40 h-8 bg-white/5 rounded-lg animate-pulse"></div>
            </div>

            <div className="w-full h-32 bg-white/5 rounded-lg animate-pulse mt-4"></div>
            <div className="w-full h-64 bg-white/5 rounded-lg animate-pulse mt-4"></div>
          </div>
        </div>

        
        <div className="w-1.5 bg-neutral-900 border-x border-white/5" />

        
        <div style={{ width: "55%" }} className="flex flex-col h-full bg-[#1e1e1e]">
          <div className="flex items-center justify-between px-4 h-11 bg-[#111118] border-b border-white/5 shrink-0">
            <div className="w-24 h-6 bg-white/5 rounded-md animate-pulse"></div>
            <div className="flex gap-2">
              <div className="w-20 h-6 bg-white/5 rounded-md animate-pulse"></div>
              <div className="w-20 h-6 bg-indigo-500/20 rounded-md animate-pulse"></div>
            </div>
          </div>
          <div className="flex-1 w-full flex p-4">
             <div className="w-12 h-full flex flex-col gap-2">
                {[1,2,3,4,5,6,7,8,9,10].map(i => <div key={i} className="w-4 h-3 bg-white/5 rounded animate-pulse"></div>)}
             </div>
             <div className="flex-1 h-full flex flex-col gap-4 pt-2">
                <div className="w-1/3 h-4 bg-white/5 rounded animate-pulse"></div>
                <div className="w-1/4 h-4 bg-white/5 rounded animate-pulse"></div>
                <div className="w-1/2 h-4 bg-white/5 rounded animate-pulse"></div>
             </div>
          </div>
          <div className="h-[44px] bg-[#111118] border-t border-white/10 flex items-center px-5">
             <div className="w-24 h-4 bg-white/5 rounded animate-pulse"></div>
          </div>
        </div>
      </div>
    </div>
  );
}

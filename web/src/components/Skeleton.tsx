export function SkeletonGrid() {
  return (
    <div className="space-y-8 animate-fade-in">
      {[1, 2].map((group) => (
        <section key={group}>
          <div className="flex items-center gap-3 mb-4">
            <div className="h-4 w-28 bg-elevated rounded animate-pulse" />
            <div className="flex-1 h-px bg-line" />
          </div>
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-3">
            {Array.from({ length: group === 1 ? 4 : 2 }).map((_, i) => (
              <SkeletonCard key={i} delay={i * 100} />
            ))}
          </div>
        </section>
      ))}
    </div>
  );
}

function SkeletonCard({ delay }: { delay: number }) {
  return (
    <div
      className="bg-card border border-line rounded-xl overflow-hidden animate-card-enter"
      style={{ animationDelay: `${delay}ms` }}
    >
      <div className="h-[2px] bg-line" />
      <div className="p-5 space-y-3">
        <div className="flex items-start gap-3.5">
          <div className="w-10 h-10 rounded-lg bg-elevated animate-pulse" />
          <div className="flex-1 space-y-2">
            <div className="h-4 w-3/4 bg-elevated rounded animate-pulse" />
            <div className="h-3 w-1/2 bg-elevated rounded animate-pulse" />
          </div>
        </div>
        <div className="h-3 w-2/3 bg-elevated rounded animate-pulse" />
        <div className="flex gap-1.5">
          <div className="h-5 w-16 bg-elevated rounded animate-pulse" />
          <div className="h-5 w-20 bg-elevated rounded animate-pulse" />
        </div>
      </div>
    </div>
  );
}

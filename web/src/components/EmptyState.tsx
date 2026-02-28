export function EmptyState({ searching }: { searching: boolean }) {
  return (
    <div className="flex flex-col items-center justify-center py-24 animate-fade-in">
      <div className="text-5xl mb-6 opacity-40">
        {searching ? "🔍" : "🧭"}
      </div>
      <h2 className="font-display font-semibold text-lg text-tx2 mb-2">
        {searching ? "No routes match" : "No routes discovered"}
      </h2>
      <p className="text-sm text-tx3 max-w-md text-center leading-relaxed">
        {searching
          ? "Try adjusting your search or namespace filter."
          : "RouteBoard is watching for Ingress and HTTPRoute resources. Create one to see it appear here in real-time."}
      </p>
      {!searching && (
        <div className="mt-6 px-4 py-3 bg-card border border-line rounded-lg font-mono text-xs text-tx3">
          <span className="text-tx2">$</span>{" "}
          kubectl create ingress my-app --rule="app.example.com/*=my-svc:80"
        </div>
      )}
    </div>
  );
}

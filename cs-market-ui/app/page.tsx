import { SkinGrid } from '@/components/skin-grid';
import { FilterSidebar } from '@/components/filter-sidebar';

export default function Home() {
  return (
    <div className="flex flex-col md:flex-row gap-6">
      <aside className="w-full md:w-64 shrink-0">
        <FilterSidebar />
      </aside>
      <main className="flex-1">
        <SkinGrid />
      </main>
    </div>
  );
}
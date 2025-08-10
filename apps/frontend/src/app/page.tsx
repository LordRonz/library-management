'use client';

import { BookTable } from "@/components/BookTable";

export default function Home() {
  return (
    <div className="py-8">
      <h1 className="text-4xl font-extrabold tracking-tight lg:text-5xl mb-4">Library Dashboard</h1>
      <p className="text-lg text-muted-foreground mb-8">Manage your book collection with ease.</p>
      <BookTable />
    </div>
  );
}
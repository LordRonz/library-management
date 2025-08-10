'use client';

import React, { useEffect } from 'react';
import { useParams } from 'next/navigation';
import { useBooks } from '@/hooks/useBooks';
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from '@/components/ui/card';
import Link from 'next/link';
import { Button } from '@/components/ui/button';
import { Loader2, ArrowLeft } from 'lucide-react';

const BookDetailsPage: React.FC = () => {
  const { id } = useParams();
  const { selectedBook, fetchBook, loading, error } = useBooks();

  useEffect(() => {
    if (id) {
      fetchBook(id as string);
    }
  }, [id, fetchBook]);

  if (loading) return (
    <div className="flex items-center justify-center h-64">
      <Loader2 className="h-8 w-8 animate-spin text-primary" />
      <p className="ml-2 text-muted-foreground">Loading book details...</p>
    </div>
  );
  
  if (error) return (
    <div className="text-center py-8 text-destructive">
      <p>Error: {error}</p>
      <p>Failed to load book details. Please try again later.</p>
    </div>
  );
  if (!selectedBook) return (
    <div className="text-center py-8 text-muted-foreground">
      <p className="text-lg">Book not found.</p>
      <Button variant="link" asChild>
        <Link href="/">
          <ArrowLeft className="mr-2 h-4 w-4" /> Back to Library
        </Link>
      </Button>
    </div>
  );

  return (
    <div className="py-8">
      <Button variant="outline" asChild className="mb-6">
        <Link href="/">
          <ArrowLeft className="mr-2 h-4 w-4" /> Back to Library
        </Link>
      </Button>
      <Card className="max-w-2xl mx-auto">
        <CardHeader>
          <CardTitle className="text-3xl font-bold">{selectedBook.title}</CardTitle>
          <CardDescription className="text-lg text-muted-foreground">
            {selectedBook.author} - {selectedBook.year}
          </CardDescription>
        </CardHeader>
        <CardContent className="grid gap-4">
          <div>
            <p className="text-sm font-medium text-muted-foreground">Genre:</p>
            <p className="text-base">{selectedBook.genre || 'N/A'}</p>
          </div>
          <div>
            <p className="text-sm font-medium text-muted-foreground">ISBN:</p>
            <p className="text-base">{selectedBook.isbn || 'N/A'}</p>
          </div>
          <div>
            <p className="text-sm font-medium text-muted-foreground">Description:</p>
            <p className="text-base">{selectedBook.description || 'N/A'}</p>
          </div>
          <div className="text-sm text-muted-foreground mt-4">
            <p>Created At: {new Date(selectedBook.created_at!).toLocaleString()}</p>
            <p>Updated At: {new Date(selectedBook.updated_at!).toLocaleString()}</p>
          </div>
        </CardContent>
      </Card>
    </div>
  );
};

export default BookDetailsPage;
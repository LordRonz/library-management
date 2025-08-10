'use client';

import React, { useEffect, useState } from 'react';
import { useBooks } from '@/hooks/useBooks';
import { Book } from '@/types/book';
import {
  Table,
  TableHeader,
  TableBody,
  TableRow,
  TableHead,
  TableCell,
} from '@/components/ui/table';
import { Button } from '@/components/ui/button';
import { BookForm } from './BookForm';
import { DeleteConfirmationDialog } from './DeleteConfirmationDialog';
import Link from 'next/link';
import { PlusCircle, Loader2, BookOpen, Edit, Trash2 } from 'lucide-react';

import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';

export const BookTable: React.FC = () => {
  const { books, fetchBooks, loading, error } = useBooks();
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const [selectedBook, setSelectedBook] = useState<Book | null>(null);

  useEffect(() => {
    fetchBooks();
  }, [fetchBooks]);

  const handleAdd = () => {
    setSelectedBook(null);
    setIsFormOpen(true);
  };

  const handleEdit = (book: Book) => {
    setSelectedBook(book);
    setIsFormOpen(true);
  };

  const handleDelete = (book: Book) => {
    setSelectedBook(book);
    setIsDeleteDialogOpen(true);
  };

  if (loading)
    return (
      <div className="flex items-center justify-center h-64">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
        <p className="ml-2 text-muted-foreground">Loading books...</p>
      </div>
    );

  if (error)
    return (
      <div className="text-center py-8 text-destructive">
        <p>Error: {error}</p>
        <p>Failed to load books. Please try again later.</p>
      </div>
    );

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between">
        <div>
          <CardTitle>Books</CardTitle>
          <CardDescription>
            Manage your books and view their details.
          </CardDescription>
        </div>
        <Button onClick={handleAdd}>
          <PlusCircle className="mr-2 h-4 w-4" /> Add Book
        </Button>
      </CardHeader>
      <CardContent>
        {books.length === 0 ? (
          <div className="text-center py-16 text-muted-foreground">
            <BookOpen className="mx-auto h-12 w-12 mb-4" />
            <p className="text-lg font-semibold">No books found</p>
            <p className="text-sm">
              Start by adding a new book to your library.
            </p>
          </div>
        ) : (
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Title</TableHead>
                <TableHead>Author</TableHead>
                <TableHead>Year</TableHead>
                <TableHead className="text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {books.map((book) => (
                <TableRow key={book.id}>
                  <TableCell className="font-medium">{book.title}</TableCell>
                  <TableCell>{book.author}</TableCell>
                  <TableCell>{book.year}</TableCell>
                  <TableCell className="flex justify-end gap-2">
                    <Button variant="outline" size="icon" asChild>
                      <Link href={`/books/${book.id}`}>
                        <BookOpen className="h-4 w-4" />
                      </Link>
                    </Button>
                    <Button
                      variant="outline"
                      size="icon"
                      onClick={() => handleEdit(book)}
                    >
                      <Edit className="h-4 w-4" />
                    </Button>
                    <Button
                      variant="destructive"
                      size="icon"
                      onClick={() => handleDelete(book)}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        )}
      </CardContent>
      <CardFooter className="flex justify-between">
        <div className="text-xs text-muted-foreground">
          Showing <strong>{books.length}</strong>{' '}
          {books.length === 1 ? 'book' : 'books'}
        </div>
      </CardFooter>

      <BookForm
        isOpen={isFormOpen}
        onOpenChange={setIsFormOpen}
        book={selectedBook}
      />

      {selectedBook && (
        <DeleteConfirmationDialog
          isOpen={isDeleteDialogOpen}
          onOpenChange={setIsDeleteDialogOpen}
          book={selectedBook}
        />
      )}
    </Card>
  );
};

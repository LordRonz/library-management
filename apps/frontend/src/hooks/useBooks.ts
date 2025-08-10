'use client';

import { useCallback } from 'react';
import { useBookContext } from '@/contexts/BookContext';
import { bookApi } from '@/lib/api';
import { CreateBookRequest, UpdateBookRequest } from '@/types/book';
import { toast } from 'sonner';

export const useBooks = () => {
  const { state, actions } = useBookContext();
  const { 
    setLoading, 
    setError, 
    setBooks, 
    addBook: addBookAction, 
    updateBook: updateBookAction, 
    deleteBook: deleteBookAction, 
    setSelectedBook 
  } = actions;

  const fetchBooks = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const books = await bookApi.getBooks();
      setBooks(books);
    } catch (error) {
      const errorMessage =
        error instanceof Error ? error.message : 'Failed to fetch books';
      setError(errorMessage);
      toast.error('Error', {
        description: errorMessage,
      });
    }
  }, [setLoading, setError, setBooks]);

  const fetchBook = useCallback(
    async (id: string) => {
      try {
        setLoading(true);
        const book = await bookApi.getBook(id);
        setSelectedBook(book);
        return book;
      } catch (error) {
        const errorMessage =
          error instanceof Error ? error.message : 'Failed to fetch book';
        setError(errorMessage);
        toast.error('Error', {
          description: errorMessage,
        });
        throw error;
      } finally {
        setLoading(false);
      }
    },
    [setLoading, setError, setSelectedBook],
  );

  const createBook = useCallback(
    async (bookData: CreateBookRequest) => {
      try {
        setLoading(true);
        const newBook = await bookApi.createBook(bookData);
        addBookAction(newBook);
        toast.success('Success', {
          description: 'Book created successfully',
        });
        return newBook;
      } catch (error) {
        const errorMessage =
          error instanceof Error ? error.message : 'Failed to create book';
        setError(errorMessage);
        toast.error('Error', {
          description: errorMessage,
        });
        throw error;
      } finally {
        setLoading(false);
      }
    },
    [setLoading, setError, addBookAction],
  );

  const updateBook = useCallback(
    async (bookData: UpdateBookRequest) => {
      try {
        setLoading(true);
        const updatedBook = await bookApi.updateBook(bookData);
        updateBookAction(updatedBook);
        toast.success('Success', {
          description: 'Book updated successfully',
        });
        return updatedBook;
      } catch (error) {
        const errorMessage =
          error instanceof Error ? error.message : 'Failed to update book';
        setError(errorMessage);
        toast.error('Error', {
          description: errorMessage,
        });
        throw error;
      } finally {
        setLoading(false);
      }
    },
    [setLoading, setError, updateBookAction],
  );

  const deleteBook = useCallback(
    async (id: string) => {
      try {
        setLoading(true);
        await bookApi.deleteBook(id);
        deleteBookAction(id);
        toast.success('Success', {
          description: 'Book deleted successfully',
        });
      } catch (error) {
        const errorMessage =
          error instanceof Error ? error.message : 'Failed to delete book';
        setError(errorMessage);
        toast.error('Error', {
          description: errorMessage,
        });
        throw error;
      } finally {
        setLoading(false);
      }
    },
    [setLoading, setError, deleteBookAction],
  );

  return {
    books: state.books,
    loading: state.loading,
    error: state.error,
    selectedBook: state.selectedBook,
    fetchBooks,
    fetchBook,
    createBook,
    updateBook,
    deleteBook,
    setSelectedBook: actions.setSelectedBook,
  };
};

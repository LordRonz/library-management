import { Book, CreateBookRequest, UpdateBookRequest } from '@/types/book';

const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080/api';

class ApiError extends Error {
  constructor(public status: number, message: string) {
    super(message);
    this.name = 'ApiError';
  }
}

const handleResponse = async (response: Response) => {
  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}));
    throw new ApiError(
      response.status,
      errorData.message || 'An error occurred',
    );
  }
  return response.json();
};

export const bookApi = {
  // Get all books
  getBooks: async (): Promise<Book[]> => {
    const response = await fetch(`${API_BASE_URL}/books`);
    return handleResponse(response);
  },

  // Get single book
  getBook: async (id: string): Promise<Book> => {
    const response = await fetch(`${API_BASE_URL}/books/${id}`);
    return handleResponse(response);
  },

  // Create book
  createBook: async (book: CreateBookRequest): Promise<Book> => {
    const response = await fetch(`${API_BASE_URL}/books`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(book),
    });
    return handleResponse(response);
  },

  // Update book
  updateBook: async (book: UpdateBookRequest): Promise<Book> => {
    const response = await fetch(`${API_BASE_URL}/books/${book.id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(book),
    });
    return handleResponse(response);
  },

  // Delete book
  deleteBook: async (id: string): Promise<void> => {
    const response = await fetch(`${API_BASE_URL}/books/${id}`, {
      method: 'DELETE',
    });
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      throw new ApiError(
        response.status,
        errorData.message || 'Failed to delete book',
      );
    }
  },
};

'use client';

import React, {
  createContext,
  useContext,
  useReducer,
  useCallback,
} from 'react';
import { Book } from '@/types/book';

interface BookState {
  books: Book[];
  loading: boolean;
  error: string | null;
  selectedBook: Book | null;
}

type BookAction =
  | { type: 'SET_LOADING'; payload: boolean }
  | { type: 'SET_ERROR'; payload: string | null }
  | { type: 'SET_BOOKS'; payload: Book[] }
  | { type: 'ADD_BOOK'; payload: Book }
  | { type: 'UPDATE_BOOK'; payload: Book }
  | { type: 'DELETE_BOOK'; payload: string }
  | { type: 'SET_SELECTED_BOOK'; payload: Book | null };

const initialState: BookState = {
  books: [],
  loading: false,
  error: null,
  selectedBook: null,
};

const bookReducer = (state: BookState, action: BookAction): BookState => {
  switch (action.type) {
    case 'SET_LOADING':
      return { ...state, loading: action.payload };
    case 'SET_ERROR':
      return { ...state, error: action.payload, loading: false };
    case 'SET_BOOKS':
      return { ...state, books: action.payload, loading: false, error: null };
    case 'ADD_BOOK':
      return { ...state, books: [...state.books, action.payload] };
    case 'UPDATE_BOOK':
      return {
        ...state,
        books: state.books.map((book) =>
          book.id === action.payload.id ? action.payload : book,
        ),
      };
    case 'DELETE_BOOK':
      return {
        ...state,
        books: state.books.filter((book) => book.id !== action.payload),
      };
    case 'SET_SELECTED_BOOK':
      return { ...state, selectedBook: action.payload };
    default:
      return state;
  }
};

interface BookContextType {
  state: BookState;
  dispatch: React.Dispatch<BookAction>;
  actions: {
    setLoading: (loading: boolean) => void;
    setError: (error: string | null) => void;
    setBooks: (books: Book[]) => void;
    addBook: (book: Book) => void;
    updateBook: (book: Book) => void;
    deleteBook: (id: string) => void;
    setSelectedBook: (book: Book | null) => void;
  };
}

const BookContext = createContext<BookContextType | undefined>(undefined);

export const BookProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [state, dispatch] = useReducer(bookReducer, initialState);

  const actions = {
    setLoading: useCallback(
      (loading: boolean) => dispatch({ type: 'SET_LOADING', payload: loading }),
      [],
    ),
    setError: useCallback(
      (error: string | null) => dispatch({ type: 'SET_ERROR', payload: error }),
      [],
    ),
    setBooks: useCallback(
      (books: Book[]) => dispatch({ type: 'SET_BOOKS', payload: books }),
      [],
    ),
    addBook: useCallback(
      (book: Book) => dispatch({ type: 'ADD_BOOK', payload: book }),
      [],
    ),
    updateBook: useCallback(
      (book: Book) => dispatch({ type: 'UPDATE_BOOK', payload: book }),
      [],
    ),
    deleteBook: useCallback(
      (id: string) => dispatch({ type: 'DELETE_BOOK', payload: id }),
      [],
    ),
    setSelectedBook: useCallback(
      (book: Book | null) =>
        dispatch({ type: 'SET_SELECTED_BOOK', payload: book }),
      [],
    ),
  };

  return (
    <BookContext.Provider value={{ state, dispatch, actions }}>
      {children}
    </BookContext.Provider>
  );
};

export const useBookContext = () => {
  const context = useContext(BookContext);
  if (!context) {
    throw new Error('useBookContext must be used within a BookProvider');
  }
  return context;
};
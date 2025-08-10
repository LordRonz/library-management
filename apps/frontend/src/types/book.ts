export interface Book {
  id: string;
  title: string;
  author: string;
  year: number;
  description?: string;
  isbn?: string;
  genre?: string;
  created_at?: string;
  updated_at?: string;
}

export interface CreateBookRequest {
  title: string;
  author: string;
  year: number;
  description?: string;
  isbn?: string;
  genre?: string;
}

export interface UpdateBookRequest extends CreateBookRequest {
  id: string;
}

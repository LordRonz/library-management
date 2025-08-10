import { z } from 'zod';

export const bookSchema = z.object({
  title: z.string().min(1, 'Title is required').max(255, 'Title is too long'),
  author: z
    .string()
    .min(1, 'Author is required')
    .max(255, 'Author name is too long'),
  year: z
    .number()
    .min(1000, 'Invalid year')
    .max(new Date().getFullYear(), 'Year cannot be in the future'),
  description: z.string().max(1000, 'Description is too long').optional(),
  isbn: z.string().max(20, 'ISBN is too long').optional(),
  genre: z.string().max(100, 'Genre is too long').optional(),
});

export type BookFormData = z.infer<typeof bookSchema>;

import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import { BookProvider } from "@/contexts/BookContext";
import { Toaster } from "@/components/ui/sonner";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

import { Github } from 'lucide-react';
import { Button } from '@/components/ui/button';
import Link from 'next/link';

export const metadata: Metadata = {
  title: "Library Management",
  description: "A simple library management system.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className="dark">
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased min-h-screen flex flex-col`}
      >
        <BookProvider>
          <header className="border-b">
            <div className="container mx-auto flex items-center justify-between p-4">
              <h1 className="text-lg font-semibold">Library Management</h1>
              <Button variant="outline" size="icon" asChild>
                <Link href="https://github.com/lordronz/library-management" target="_blank">
                  <Github className="h-4 w-4" />
                </Link>
              </Button>
            </div>
          </header>
          <main className="container mx-auto p-4 flex-grow flex items-center justify-center">
            <div className="w-full max-w-4xl">{children}</div>
          </main>
          <footer className="border-t py-4">
            <div className="container mx-auto text-center text-sm text-muted-foreground">
              <p>&copy; {new Date().getFullYear()} By-Food Assessment. All rights reserved.</p>
            </div>
          </footer>
          <Toaster />
        </BookProvider>
      </body>
    </html>
  );
}

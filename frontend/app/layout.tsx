import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import Navbar from "@/components/layout/Navbar";
import { Toaster } from "react-hot-toast";

const inter = Inter({
  subsets: ["latin"],
  variable: "--font-sans",
});

export const metadata: Metadata = {
  title: "FilmFolk - Discover, Review, Connect",
  description: "Your social platform for movie lovers. Discover new films, share reviews, and connect with fellow cinephiles.",
  keywords: "movies, reviews, films, cinema, social, community",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className="dark">
      <body className={`${inter.variable} font-sans antialiased`}>
        <div className="relative min-h-screen bg-background">
          {/* Background gradient */}
          <div className="fixed inset-0 -z-10 overflow-hidden">
            <div className="absolute top-0 left-1/4 w-96 h-96 bg-primary/20 rounded-full blur-3xl animate-pulse" />
            <div className="absolute bottom-0 right-1/4 w-96 h-96 bg-purple-500/20 rounded-full blur-3xl animate-pulse delay-1000" />
          </div>

          <Navbar />
          <main className="pt-16">{children}</main>
          <Toaster
            position="top-right"
            toastOptions={{
              duration: 4000,
            }}
          />
        </div>
      </body>
    </html>
  );
}

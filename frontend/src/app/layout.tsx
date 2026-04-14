import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import Navbar from "@/components/Navbar";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "CodeForge | Competitive Programming Platform",
  description: "Next-gen competitive programming platform to enhance your coding skills.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className="dark" suppressHydrationWarning>
      <body className={`${inter.className} min-h-screen bg-neutral-950 text-neutral-50 antialiased selection:bg-indigo-500/30`} suppressHydrationWarning>
        <Navbar />
        <main className="flex-1 w-full h-full">
          {children}
        </main>
      </body>
    </html>
  );
}

import type { Metadata } from "next";
import { Inter, JetBrains_Mono } from "next/font/google";
import "./globals.css";
import Header from "@/components/Header";
import Footer from "@/components/Footer";

const inter = Inter({
  subsets: ["latin"],
  variable: "--font-inter",
});

const jetbrainsMono = JetBrains_Mono({
  subsets: ["latin"],
  variable: "--font-mono",
});

export const metadata: Metadata = {
  title: "Hatch — Self-hostable HTTP request inspector + mocker",
  description: "Hatch is a self-hostable HTTP request inspector and mocker. One Go binary, SQLite, no SaaS, no per-request fee. Open source, MIT.",
  openGraph: {
    title: "Hatch — Self-hostable HTTP request inspector + mocker",
    description: "Hatch is a self-hostable HTTP request inspector and mocker. One Go binary, SQLite, no SaaS, no per-request fee. Open source, MIT.",
    url: "https://hatch.surf",
    siteName: "Hatch",
    images: [
      {
        url: "/brand/og/default.png",
        width: 1200,
        height: 630,
      },
    ],
    type: "website",
  },
  twitter: {
    card: "summary_large_image",
    title: "Hatch — Self-hostable HTTP request inspector + mocker",
    description: "Hatch is a self-hostable HTTP request inspector and mocker. One Go binary, SQLite, no SaaS, no per-request fee. Open source, MIT.",
    images: ["/brand/og/default.png"],
  },
  icons: {
    icon: "/brand/favicon/favicon-48.png",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className={`${inter.variable} ${jetbrainsMono.variable}`}>
      <body className="min-h-screen bg-zinc-950 text-white font-sans antialiased">
        <a href="#main-content" className="sr-only focus:not-sr-only focus:absolute focus:top-4 focus:left-4 focus:z-50 focus:px-4 focus:py-2 focus:bg-blue-500 focus:text-white focus:rounded-lg">
          Skip to content
        </a>
        <Header />
        <main id="main-content" role="main">
          {children}
        </main>
        <Footer />
      </body>
    </html>
  );
}

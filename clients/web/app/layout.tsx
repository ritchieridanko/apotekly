import type { Metadata } from "next";
import { JetBrains_Mono, Ubuntu } from "next/font/google";

import "./globals.css";
import Providers from "./providers";

const jetBrainsMono = JetBrains_Mono({
  fallback: ["Arial", "Helvetica", "sans-serif"],
  style: ["italic", "normal"],
  subsets: ["latin"],
  variable: "--font-jet-brains-mono",
  weight: ["100", "200", "300", "400", "500", "600", "700", "800"],
});

const ubuntu = Ubuntu({
  fallback: ["Arial", "Helvetica", "sans-serif"],
  style: ["normal", "italic"],
  subsets: ["latin"],
  variable: "--font-ubuntu",
  weight: ["300", "400", "500", "700"],
});

export const metadata: Metadata = {
  title: "Apotek.ly",
  description: "Healthcare E-Commerce",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html
      lang="en"
      className={`h-full antialiased ${jetBrainsMono.variable} ${ubuntu.variable}`}
    >
      <body className="min-h-full flex flex-col">
        <Providers>{children}</Providers>
      </body>
    </html>
  );
}

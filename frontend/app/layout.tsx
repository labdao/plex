import "@/styles/globals.css";

import { cn } from "@lib/utils";
import { Metadata } from "next";
import { Inter, Space_Grotesk } from "next/font/google";

import UserLoader from "@/app/components/auth/UserLoader";
import Footer from "@/components/global/Footer";
import Header from "@/components/global/Header";
import { Providers } from "@/lib/providers";

export const fontSans = Inter({
  subsets: ["latin"],
  variable: "--font-sans",
});

export const fontPrimary = Space_Grotesk({
  subsets: ["latin"],
  variable: "--font-primary",
});

export default function RootLayout(props: React.PropsWithChildren) {
  return (
    <html lang="en">
      <body className={cn("min-h-screen bg-background font-primary antialiased", fontSans.variable, fontPrimary.variable)}>
        <Providers>
          <div className="flex flex-col min-h-screen w-full bg-gray-100">
            <Header />
            <div className="grow">
              <UserLoader>{props.children}</UserLoader>
            </div>
            <Footer />
          </div>
        </Providers>
      </body>
    </html>
  );
}

export const metadata: Metadata = {
  title: "Lab Exchange",
  icons: [
    {
      rel: "shortcut icon",
      url: "/icons/favicon.ico",
    },
    {
      rel: "icon",
      type: "image/png",
      sizes: "32x32",
      url: "/icons/favicon-32x32.png",
    },
    {
      rel: "icon",
      type: "image/png",
      sizes: "16x16",
      url: "/icons/favicon-16x16.png",
    },
    {
      rel: "apple-touch-icon",
      sizes: "180x180",
      url: "/icons/apple-touch-icon.png",
    },
  ],
  manifest: "/site.webmanifest",
  other: {
    "mask-icon": "/icons/safari-pinned-tab.svg",
    "msapplication-TileColor": "#6bdbad",
    "msapplication-config": "/browserconfig.xml",
  },
};

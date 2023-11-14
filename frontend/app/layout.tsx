import "@/styles/globals.css";

import { Metadata } from "next";
import { Space_Grotesk } from "next/font/google";
import localFont from "next/font/local";

import UserLoader from "@/components/auth/UserLoader";
import Footer from "@/components/global/Footer";
import Header from "@/components/global/Header";
import { Providers } from "@/lib/providers";
import { cn } from "@/lib/utils";

const fontRaster = localFont({
  src: [
    {
      path: "../fonts/FKRasterRomanCompact-Blended.woff2",
      weight: "400",
      style: "normal",
    },
  ],
  variable: "--font-raster",
});

const fontMono = localFont({
  src: [
    {
      path: "../fonts/PPFraktionMono-Regular.woff2",
      weight: "400",
      style: "normal",
    },
    {
      path: "../fonts/PPFraktionMono-Bold.woff2",
      weight: "700",
      style: "normal",
    },
    {
      path: "../fonts/PPFraktionMono-RegularItalic.woff2",
      weight: "400",
      style: "italic",
    },
    {
      path: "../fonts/PPFraktionMono-BoldItalic.woff2",
      weight: "700",
      style: "italic",
    },
  ],
  variable: "--font-mono",
});

const fontBody = localFont({
  src: [
    {
      path: "../fonts/PPNeueMontreal-Regular.woff2",
      weight: "400",
      style: "normal",
    },
    {
      path: "../fonts/PPNeueMontreal-Bold.woff2",
      weight: "700",
      style: "normal",
    },
    {
      path: "../fonts/PPNeueMontreal-Italic.woff2",
      weight: "400",
      style: "italic",
    },
    {
      path: "../fonts/PPNeueMontreal-BoldItalic.woff2",
      weight: "700",
      style: "italic",
    },
  ],
  variable: "--font-body",
});

const fontHeading = Space_Grotesk({
  variable: "--font-heading",
  weight: ["500"],
  style: ["normal"],
  subsets: ["latin"],
});

export default function RootLayout(props: React.PropsWithChildren) {
  return (
    <html lang="en">
      <body
        className={cn(
          "min-h-screen bg-background font-body antialiased",
          fontRaster.variable,
          fontMono.variable,
          fontBody.variable,
          fontHeading.variable
        )}
      >
        <Providers>
          <div className="flex flex-col w-full min-h-screen bg-gray-100">
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

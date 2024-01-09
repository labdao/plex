"use client";

import { SiDiscord, SiGithub, SiTwitter } from "@icons-pack/react-simple-icons";
import Link from "next/link";

import { Button } from "../ui/button";
import Logo from "./Logo";

export default function Footer() {
  const footerSections = [
    {
      title: "app",
      links: [
        {
          title: "experiments",
          href: "/experiments",
        },
        {
          title: "data",
          href: "/data",
        },
        {
          title: "tasks",
          href: "/tasks",
        },
      ],
    },
    {
      title: "community",
      links: [
        {
          title: "discord",
          href: "https://discord.gg/labdao",
        },
        {
          title: "docs",
          href: "https://docs.labdao.xyz/",
        },
      ],
    },
    {
      title: "convexity",
      links: [
        // {
        //   title: "careers",
        //   href: "http://convexity.bio",
        // },
        {
          title: "blog",
          href: "https://medium.com/@labdao",
        },
      ],
    },
    {
      title: "help",
      links: [
        {
          title: "contact",
          href: "mailto:info@convexity.bio",
        },
        {
          title: "discord",
          href: "https://discord.gg/labdao",
        },
      ],
    },
  ];

  return (
    <div className="container py-10">
      <div className="flex flex-col items-center justify-between gap-8 p-6 mx-auto text-center bg-brand-pattern md:text-left md:flex-row from-pink-50 via-yellow-50 to-yellow-50 rounded-xl">
        <div>
          <div className="mb-2 text-3xl font-heading">Operated by Convexity Labs</div>
          <p className="font-mono opacity-75">Advancing safe and decentralised computational life sciences research</p>
        </div>
        {/*
        <Button variant="secondary" asChild>
          <a href="https://www.convexity.bio" target="_blank" rel="noopener">
            Learn More
          </a>
        </Button>
        */}
      </div>
      <div className="flex flex-wrap mt-8">
        <div className="flex flex-row items-center justify-between order-last w-full gap-4 lg:items-start lg:w-1/4 lg:justify-normal lg:flex-col lg:order-first">
          <div className="flex gap-4">
            <div>
              <Logo className="w-8 h-8 opacity-25 text-muted-foreground" />
            </div>
            <Link href="https://twitter.com/lab_dao" target="_blank" rel="noopener" className="text-muted-foreground hover:text-foreground">
              <SiTwitter size={32} />
            </Link>
            <Link href="https://github.com/labdao" target="_blank" rel="noopener" className="text-muted-foreground hover:text-foreground">
              <SiGithub size={32} />
            </Link>
            <Link href="https://discord.gg/labdao" target="_blank" rel="noopener" className="text-muted-foreground hover:text-foreground">
              <SiDiscord size={32} />
            </Link>
          </div>
          <div className="opacity-25 text-muted-foreground">&copy; 2023 Openlab&nbsp;Association</div>
        </div>
        <div className="flex flex-wrap w-full mb-8 lg:w-3/4">
          {footerSections.map((section) => (
            <div key={section.title} className="w-full pl-0 pr-4 sm:w-1/2 lg:px-8 md:w-1/4">
              <div className="font-mono text-xl font-semibold">{section.title}</div>
              <ul className="mt-2">
                {section.links.map((link) => (
                  <li key={link.title} className="mb-2">
                    <Link
                      href={link.href}
                      className="text-muted-foreground hover:text-foreground"
                      target={link.href.startsWith("http") ? "_blank" : undefined}
                    >
                      {link.title}
                    </Link>
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}

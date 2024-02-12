"use client";

import { getAccessToken } from "@privy-io/react-auth";
import { ColumnDef } from "@tanstack/react-table";
import dayjs from "dayjs";
import backendUrl from "lib/backendUrl";
import { Code2Icon, CopyCheck, CopyIcon } from "lucide-react";
import React, { useEffect, useState } from "react";

import { Breadcrumbs } from "@/components/global/Breadcrumbs";
import { Button } from "@/components/ui/button";
import { DataTable } from "@/components/ui/data-table";
import { DataTableColumnHeader } from "@/components/ui/data-table-column-header";
import { ScrollArea, ScrollBar } from "@/components/ui/scroll-area";

export default function ListApiKeys() {
    interface ApiKey {
        Key: string;
        Scope: string;
        CreatedAt: string;
        ExpiresAt: string;
    }

    const handleApiKeyClick = async (apiKey: string) => {
        try {
            await navigator.clipboard.writeText(apiKey);
            setCopiedKey(apiKey);
            setTimeout(() => {
                setCopiedKey(null);
            }, 2000);
        } catch (err) {
            console.error("Error copying to clipboard:", err);
        }
    }

    const columns: ColumnDef<ApiKey>[] = [
        {
            accessorKey: "Key",
            header: ({ column }) => <DataTableColumnHeader column={column} title="Key" />,
            sortingFn: "alphanumeric",
            cell: ({ row }) => {
                const apiKey: string = row.getValue("Key");
                const trimmedApiKey = apiKey.length > 10 ? `${apiKey.substring(0, 10)}...${apiKey.slice(-4)}` : apiKey;
                return (
                    <div
                        style={{ cursor: "pointer", display: "inline-flex", alignItems: "center", gap: "8px"}}
                        onClick={() => handleApiKeyClick(apiKey)}
                        className="copy-container"
                    >
                        {trimmedApiKey}
                        {copiedKey === apiKey ? (
                            <CopyCheck size={16} />
                        ) : (
                            <span className="copy-icon-hover">
                                <CopyIcon size={16} className="copy-icon"/>
                            </span>
                        )}
                    </div>
                );
            },
        },
        {
            accessorKey: "Scope",
            header: "Scope",
        },
        {
            accessorKey: "CreatedAt",
            header: "Created At",
            cell : ({ row }) => {
                return dayjs(row.getValue("CreatedAt")).format("YYYY-MM-DD HH:mm:ss")
            },
        },
        {
            accessorKey: "ExpiresAt",
            header: "Expires At",
            cell : ({ row }) => {
                return dayjs(row.getValue("ExpiresAt")).format("YYYY-MM-DD HH:mm:ss")
            },
        },
    ];

    const [apiKeys, setApiKeys] = useState<ApiKey[]>([]);
    const [loading, setLoading] = useState(true);
    const [copiedKey, setCopiedKey] = useState<string | null>(null);

    const handleGenerateApiKey = async () => {
        try {
            const authToken = await getAccessToken();
            const response = await fetch(`${backendUrl()}/api-keys`, {
                method: "POST",
                headers: {
                    Authorization: `Bearer ${authToken}`,
                },
            });

            if (!response.ok) {
                console.error("Error generating API key:", response);
                return;
            }

            const newApiKey = await response.json();
            setApiKeys([newApiKey, ...apiKeys]);
        } catch (error) {
            console.error("Error generating API key:", error);
        }
    }

    useEffect(() => {
        const fetchApiKeys = async () => {
            setLoading(true);
            try {
                const authToken = await getAccessToken();
                const response = await fetch(`${backendUrl()}/api-keys`, {
                    headers: {
                        Authorization: `Bearer ${authToken}`,
                    },
                });

                if (!response.ok) {
                    console.error("Error fetching API keys:", response);
                    return;
                }

                const data = await response.json();
                setApiKeys(data);
            } catch (error) {
                console.error("Error fetching API keys:", error);
            } finally {
                setLoading(false);
            }
        }

        fetchApiKeys();
    }, []);

    return (
        <div className="relative flex flex-col h-screen max-w-full grow">
        <Breadcrumbs
            items={[{ name: "My API Keys", href: "/api-keys" }]}
            actions={
                <Button size="sm" onClick={handleGenerateApiKey}>
                    <Code2Icon />
                    Generate API Key
                </Button>
            }
        />
        <ScrollArea className="bg-white grow w-[calc(100vw-12rem)]">
            <DataTable columns={columns} data={apiKeys} loading={loading} />
            <ScrollBar orientation="horizontal" />
            <ScrollBar orientation="vertical" />
        </ScrollArea>
        {/* If you have pagination or other components, include them here */}
    </div>
    )
}

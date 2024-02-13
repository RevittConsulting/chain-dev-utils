"use client";

import { useState, useEffect, useRef } from "react";
import ChainCard from "@/components/chain-card";
import { Chain } from "../types/chain";
import { getChains } from "@/api/chain";
import { ScrollArea, ScrollBar } from "@/components/ui/scroll-area";

export default function Home() {
  const [chains, setChains] = useState<Chain[]>([]);

  useEffect(() => {
    const fetchChains = async () => {
      const c = await getChains();
      console.log(c);
      setChains(c);
    };

    fetchChains();
  }, []);

  return (
    <div className="py-20 mx-20">
      <ScrollArea className="w-full rounded-md border">
        <div className="flex w-max space-x-4 p-4">
          {chains.map((chain, index) => (
            <ChainCard key={index} chain={chain} />
          ))}
        </div>
        <ScrollBar orientation="horizontal" />
      </ScrollArea>
    </div>
  );
}
"use client";

import "./skin.scss";

import { PluginCommands } from "molstar/lib/mol-plugin/commands";
import { PluginConfig } from "molstar/lib/mol-plugin/config";
import { createPluginUI } from "molstar/lib/mol-plugin-ui";
import { PluginUIContext } from "molstar/lib/mol-plugin-ui/context";
import { renderReact18 } from "molstar/lib/mol-plugin-ui/react18";
import { DefaultPluginUISpec, PluginUISpec } from "molstar/lib/mol-plugin-ui/spec";
import { Color } from "molstar/lib/mol-util/color";
import { createRef, useEffect } from "react";

import { cn } from "@/lib/utils";

declare global {
  interface Window {
    molstar?: PluginUIContext;
  }
}

interface MolstarProps {
  url?: string;
  showControls?: boolean;
  isExpanded?: boolean;
  className?: string;
}

const Molstar = ({ url, isExpanded, showControls, className }: MolstarProps) => {
  const parent = createRef<HTMLDivElement>();

  // In debug mode of react's strict mode, this code will
  // be called twice in a row, which might result in unexpected behavior.
  useEffect(() => {
    const MySpec: PluginUISpec = {
      ...DefaultPluginUISpec(),

      // See config options here: https://github.com/molstar/molstar/blob/master/src/mol-plugin/config.ts#L25
      config: [
        [PluginConfig.Viewport.ShowExpand, false],
        [PluginConfig.Viewport.ShowControls, false],
      ],

      // See layout options here: https://github.com/molstar/molstar/blob/master/src/mol-plugin/layout.ts#L23
      layout: {
        initial: {
          isExpanded: isExpanded,
          showControls: showControls,
        },
      },
    };
    async function init() {
      try {
        console.log("Initializing Molstar...");
        window.molstar = await createPluginUI({
          target: parent.current as HTMLDivElement,
          spec: MySpec,
          render: renderReact18,
        });

        const renderer = window.molstar.canvas3d!.props.renderer;
        PluginCommands.Canvas3D.SetSettings(window.molstar, {
          settings: { renderer: { ...renderer, backgroundColor: 0xf9fafb as Color } },
        });
      } catch (error) {
        console.error("Error initializing Molstar:", error);
      }
    }

    async function update() {
      try {
        console.log("Updating Molstar...");
        if (window.molstar) {
          window.molstar.clear();
          if (url) {
            const data = await window.molstar.builders.data.download({ url: url }, { state: { isGhost: true } });
            const trajectory = await window.molstar.builders.structure.parseTrajectory(data, "pdb");
            await window.molstar.builders.structure.hierarchy.applyPreset(trajectory, "default");
          }
        }
      } catch (error) {
        console.error("Error updating Molstar:", error);
      }
    }
    if (!window.molstar) {
      init();
      update();
    } else {
      update();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [url]);

  useEffect(() => {
    return () => {
      window.molstar?.dispose();
      window.molstar = undefined;
    };
  }, []);

  return <div ref={parent} className={cn(className, "relative z-10")} />;
};

export default Molstar;

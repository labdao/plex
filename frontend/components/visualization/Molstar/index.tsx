"use client";

import "./skin.scss";

import { PluginCommands } from "molstar/lib/mol-plugin/commands";
import { PluginConfig } from "molstar/lib/mol-plugin/config";
import { createPluginUI } from "molstar/lib/mol-plugin-ui";
import { PluginUIContext } from "molstar/lib/mol-plugin-ui/context";
import { renderReact18 } from "molstar/lib/mol-plugin-ui/react18";
import { DefaultPluginUISpec, PluginUISpec } from "molstar/lib/mol-plugin-ui/spec";
import { ColorNames } from "molstar/lib/mol-util/color/names";
import { createRef, useEffect } from "react";
import { Color } from "molstar/lib/mol-util/color";

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

const Molstar = ({ url, showControls, isExpanded, className }: MolstarProps) => {
  const parent = createRef<HTMLDivElement>();

  // In debug mode of react's strict mode, this code will
  // be called twice in a row, which might result in unexpected behavior.
  useEffect(() => {
    const MySpec: PluginUISpec = {
      ...DefaultPluginUISpec(),

      // See config options here: https://github.com/molstar/molstar/blob/master/src/mol-plugin/config.ts#L25
      config: [[PluginConfig.Viewport.ShowExpand, false]],

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
        window.molstar = await createPluginUI({
          target: parent.current as HTMLDivElement,
          spec: MySpec,
          render: renderReact18,
        });

        const renderer = window.molstar.canvas3d!.props.renderer;
        PluginCommands.Canvas3D.SetSettings(window.molstar, {
          settings: { renderer: { ...renderer, backgroundColor: 0xf9fafb as Color } },
        });
        if (url) {
          const data = await window.molstar.builders.data.download({ url: url }, { state: { isGhost: true } });
          const trajectory = await window.molstar.builders.structure.parseTrajectory(data, "pdb");
          await window.molstar.builders.structure.hierarchy.applyPreset(trajectory, "default");
        } else {
          window.molstar.clear();
        }
      } catch (error) {
        console.error("Error initializing Molstar:", error);
      }
    }
    init();

    return () => {
      window.molstar?.dispose();
      window.molstar = undefined;
    };
  }, [isExpanded, parent, showControls, url]);

  return <div ref={parent} style={{ position: "relative" }} className={className || ""} />;
};

export default Molstar;

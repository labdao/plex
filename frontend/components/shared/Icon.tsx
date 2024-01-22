import { LucideProps } from "lucide-react";
import dynamicIconImports from "lucide-react/dynamicIconImports";
import dynamic from "next/dynamic";
import React, { lazy, Suspense } from "react";

const fallback = <div className="w-[16px] h-[16px]" />;

interface IconProps extends Omit<LucideProps, "ref"> {
  name: keyof typeof dynamicIconImports;
}

const Icon = ({ name, ...props }: IconProps) => {
  const LucideIcon = dynamic(dynamicIconImports[name], { loading: () => <p>Loading...</p> });

  return (
    <Suspense fallback={fallback}>
      <LucideIcon strokeWidth={1.5} size={16} absoluteStrokeWidth {...props} />
    </Suspense>
  );
};

export default Icon;

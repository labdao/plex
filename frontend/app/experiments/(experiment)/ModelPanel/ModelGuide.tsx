import { HelpCircleIcon } from "lucide-react";

import { ToolDetail } from "@/lib/redux";

const renderDescriptionParagraphs = (text: string) => {
  const paragraphs = text.split("\n");
  const hasNumberedSteps = paragraphs.some((paragraph) => paragraph.match(/^\d+\. /));

  if (hasNumberedSteps) {
    const steps = paragraphs.filter((paragraph) => paragraph.match(/^\d+\. /));
    const nonStepParagraphs = paragraphs.filter((paragraph) => !paragraph.match(/^\d+\. /));

    return (
      <>
        {nonStepParagraphs.map((paragraph, index) => (
          <p key={index} className="mt-2">
            {paragraph}
          </p>
        ))}
        <ol className="mt-2 ml-5 list-decimal list-outside">
          {steps.map((step, index) => (
            <li key={index} className="pb-4 mt-1 mb-4 border-b last:border-0 last:m-0 last:p-0 marker:text-foreground">
              {step.replace(/^\d+\. /, "")}
            </li>
          ))}
        </ol>
      </>
    );
  } else {
    return paragraphs.map((paragraph, index) => (
      <p key={index} className="mt-2">
        {paragraph}
      </p>
    ));
  }
};

interface ModelGuideProps {
  tool: ToolDetail;
}

export default function ModelGuide({ tool }: ModelGuideProps) {
  return (
    <div>
      <div className="text-left uppercase font-heading">How to Write Prompts</div>
      <div className="pt-0">
        <div className="space-y-2 text-muted-foreground">{renderDescriptionParagraphs(tool.ToolJson.guide)}</div>
      </div>
    </div>
  );
}

import { HelpCircleIcon } from "lucide-react";

import { ToolDetail } from "@/lib/redux";
import { ReactElement, JSXElementConstructor, ReactNode, PromiseLikeOfReactNode, JSX } from "react";

const renderDescriptionParagraphs = (text: string) => {
  const paragraphs = text.split("\n");

  const renderParagraph = (paragraph: string) => {
    const modifiedParagraph = paragraph.replace(/\*(.*?)\*/g, (match, p1) => `<strong>${p1}</strong>`);
    return { __html: modifiedParagraph };
  };

  const elements = [];
  let listItems: JSX.Element[] = [];

  paragraphs.forEach((paragraph, index) => {
    if (paragraph.match(/^\d+\. /)) {
      listItems.push(
        <li key={`step-${index}`} className="pb-4 mt-1 mb-4 border-b last:border-0 last:m-0 last:p-0 marker:text-foreground"
            dangerouslySetInnerHTML={renderParagraph(paragraph.replace(/^\d+\. /, ""))}>
        </li>
      );
    } else {
      if (listItems.length > 0) {
        elements.push(<ol key={`ol-${elements.length}`} className="mt-2 ml-5 list-decimal list-outside">{listItems}</ol>);
        listItems = [];
      }
      elements.push(
        <p key={`nonstep-${index}`} className="mt-2" dangerouslySetInnerHTML={renderParagraph(paragraph)}>
        </p>
      );
    }
  });

  if (listItems.length > 0) {
    elements.push(<ol key={`ol-${elements.length}`} className="mt-2 ml-5 list-decimal list-outside">{listItems}</ol>);
  }

  return <>{elements}</>;
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

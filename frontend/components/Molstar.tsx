import React, { useEffect, useRef } from 'react';

type MolstarComponentProps = {
  moleculeUrl: string;
  customDataFormat: string; // e.g., "cif"
};

declare global {
  namespace JSX {
    interface IntrinsicElements {
      'pdbe-molstar': React.DetailedHTMLProps<React.HTMLAttributes<HTMLElement>, HTMLElement> & {
        'custom-data-url': string;
        'custom-data-format': string;
        style?: React.CSSProperties;
      };
    }
  }
}

const MolstarComponent: React.FC<MolstarComponentProps> = ({ moleculeUrl, customDataFormat }) => {
  const molstarRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const loadScript = (src: string) => {
      const script = document.createElement('script');
      script.src = src;
      script.async = false;
      document.head.appendChild(script);
    };

    // loadScript('https://cdn.jsdelivr.net/npm/babel-polyfill/dist/polyfill.min.js');
    // loadScript('https://cdn.jsdelivr.net/npm/@webcomponents/webcomponentsjs/webcomponents-lite.js');
    // loadScript('https://cdn.jsdelivr.net/npm/@webcomponents/webcomponentsjs/custom-elements-es5-adapter.js');
    loadScript('https://www.ebi.ac.uk/pdbe/pdb-component-library/js/pdbe-molstar-component-3.1.3.js');

    const link = document.createElement('link');
    link.rel = 'stylesheet';
    link.href = 'https://www.ebi.ac.uk/pdbe/pdb-component-library/css/pdbe-molstar-3.1.3.css';
    document.head.appendChild(link);

    if (molstarRef.current) {
      molstarRef.current.setAttribute('custom-data-url', moleculeUrl);
      molstarRef.current.setAttribute('custom-data-format', customDataFormat);
    }
  }, [moleculeUrl, customDataFormat]);

  return (
    <div ref={molstarRef} style={{ width: '50px', height: '50px' }}>
        <pdbe-molstar
            custom-data-url={moleculeUrl}
            custom-data-format={customDataFormat}
            hide-controls="true"
            style={{ width: '25%', height: '25%' }}
            // style={{ width: '100%', height: '100%' }}
        ></pdbe-molstar>
    </div>
  );
};

export default MolstarComponent;

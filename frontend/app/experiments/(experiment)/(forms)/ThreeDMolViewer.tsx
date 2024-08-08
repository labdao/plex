import React, { useState, useEffect, useRef } from 'react';
import axios from 'axios';
// This import will be handled dynamically below, so no need for an import statement at the top.

interface AtomSpec {
    chain: string;
    resi: number;
}

interface ThreeDMolViewerProps {
    onSubmit: (data: { pdb: File | null; binderLength: number; hotspots: string }) => void;
}

const ThreeDMolViewer: React.FC<ThreeDMolViewerProps> = ({ onSubmit }) => {
    const [viewer, setViewer] = useState<any>(null);
    const [selectedResidues, setSelectedResidues] = useState<Record<string, boolean>>({});
    const [binderLength, setBinderLength] = useState(90);
    const fileInput = useRef<HTMLInputElement>(null);

    useEffect(() => {
        (async () => {
            const module = await import('3dmol/build/3Dmol.js');
            const $3Dmol = module.default ? module.default : module;
            const element = document.getElementById('container-01') as HTMLElement;
            const config = { backgroundColor: 'white' };
            const viewer = $3Dmol.createViewer(element, config);
            setViewer(viewer);
        })();
    }, []);

    const handleFileUpload = (event: React.ChangeEvent<HTMLInputElement>) => {
        const file = event.target.files?.[0];
        if (file && viewer) {
            const reader = new FileReader();
            reader.onload = function(e) {
                const result = e.target?.result as string;
                viewer.addModel(result, "pdb");
                updateStyles();
                viewer.getModel().setClickable({}, true, (atom: AtomSpec) => {
                    const residueKey = `${atom.chain}${atom.resi}`;
                    setSelectedResidues(prev => ({
                        ...prev,
                        [residueKey]: !prev[residueKey],
                    }));
                });
                viewer.render();
            };
            reader.readAsText(file);
        }
    };

    const updateStyles = () => {
        if (!viewer) return;
        viewer.setStyle({}, { cartoon: { color: 'grey', opacity: 0.2 } });
        Object.keys(selectedResidues).forEach(residueKey => {
            const [chain, resi] = residueKey.split(':');
            if (selectedResidues[residueKey]) {
                viewer.setStyle({ chain, resi }, { cartoon: { color: 'red', opacity: 1.0 } });
            }
        });
        viewer.render();
    };

    const formatSelectedResidues = () => {
        return Object.keys(selectedResidues).filter(key => selectedResidues[key]).join(', ');
    };

    const handleBinderLengthChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setBinderLength(Number(event.target.value));
    };

    const handleSubmit = async () => {
        const pdbFile = fileInput.current?.files?.[0] ?? null;
        const hotspots = formatSelectedResidues();
        onSubmit({ pdb: pdbFile, binderLength, hotspots });
    };

    return (
        <div>
            <div id="container-01" style={{ width: '800px', height: '600px' }}></div>
            <input type="file" ref={fileInput} onChange={handleFileUpload} />
            <input type="text" value={formatSelectedResidues()} readOnly />
            <div>
                <label>Binder Length:</label>
                <input
                    type="range"
                    min="60"
                    max="120"
                    value={binderLength}
                    onChange={handleBinderLengthChange}
                />
                <span style={{ marginLeft: '10px' }}>{binderLength}</span>
                <button onClick={handleSubmit}>Submit</button>
            </div>
        </div>
    );
};

export default ThreeDMolViewer;
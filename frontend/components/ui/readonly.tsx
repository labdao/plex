import React from 'react';
import { cn } from "@/lib/utils";

// Styles for the read-only overlay
const overlayStyles: React.CSSProperties = {
    position: 'absolute',
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
    backgroundColor: 'rgba(255, 255, 255, 0.2)',
    cursor: 'default',
    zIndex: 10
};

interface ReadOnlyWrapperProps {
    children: React.ReactNode;
    readOnly?: boolean;
}

const ReadOnlyWrapper: React.FC<ReadOnlyWrapperProps> = ({ children, readOnly = false }) => {
    return (
        <div style={{ position: 'relative' }} className={cn(readOnly && "relative")}>
            {children}
            {readOnly && <div style={overlayStyles} aria-hidden="true"></div>}
        </div>
    );
};
export default ReadOnlyWrapper;

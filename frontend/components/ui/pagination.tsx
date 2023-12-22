import React from 'react';

interface PaginationProps {
    currentPage: number;
    totalPages: number;
    onPageChange: (page: number) => void;
}

export const Pagination: React.FC<PaginationProps> = ({ currentPage, totalPages, onPageChange }) => {
    const renderPageNumbers = () => {
        const pageNumbers: (number | string)[] = [];

        if (totalPages <= 10) {
            for (let i = 1; i <= totalPages; i++) {
                pageNumbers.push(i);
            }
        } else {
            pageNumbers.push(1);

            let rangeStart = Math.max(2, currentPage - 4);
            let rangeEnd = Math.min(totalPages - 1, currentPage + 4);

            if (currentPage <= 6) {
                rangeStart = 2;
                rangeEnd = 10;
            }

            if (currentPage >= totalPages - 5) {
                rangeStart = totalPages - 9;
                rangeEnd = totalPages - 1;
            }

            if (rangeStart > 2) {
                pageNumbers.push('...');
            }

            for (let i = rangeStart; i <= rangeEnd; i++) {
                pageNumbers.push(i);
            }

            if (rangeEnd < totalPages - 1) {
                pageNumbers.push('...');
            }

            pageNumbers.push(totalPages);
        }

        return pageNumbers;
    };

    return (
        <div className="pagination-controls" style={{ marginTop: '20px', textAlign: 'center' }}>
            <button 
                onClick={() => onPageChange(1)} 
                disabled={currentPage === 1}
                style={{ marginRight: '10px' }}    
            >
                First
            </button>
            {renderPageNumbers().map((page, index) => (
                typeof page === 'number' ? (
                    <button 
                        key={page} 
                        onClick={() => onPageChange(page)} 
                        disabled={currentPage === page}
                        style={{ margin: '0 5px' }}
                    >
                        {page}
                    </button>
                ) : (
                    <span key={index} style={{ margin: '0 5px' }}>...</span>
                )
            ))}
            <button 
                onClick={() => onPageChange(totalPages)} 
                disabled={currentPage === totalPages}
                style={{ marginLeft: '10px' }}
            >
                Last
            </button>
        </div>
    );
};
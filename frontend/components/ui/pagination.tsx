import React from 'react';

interface PaginationProps {
    currentPage: number;
    totalPages: number;
    onPageChange: (page: number) => void;
}

export const Pagination: React.FC<PaginationProps> = ({ currentPage, totalPages, onPageChange }) => {
    const renderPageNumbers = () => {
        const pageNumbers: (number | string)[] = [];

        pageNumbers.push(1);

        let rangeStart: number, rangeEnd: number;

        if (totalPages <= 10) {
            rangeStart = 2;
            rangeEnd = totalPages - 1;
        } else {
            if (currentPage <= 6) {
                rangeStart = 2;
                rangeEnd = 10;
            } else if (currentPage >= totalPages - 5) {
                rangeStart = totalPages - 9;
                rangeEnd = totalPages - 1;
            } else {
                rangeStart = currentPage - 4;
                rangeEnd = currentPage + 4;
            }
        }

        if (rangeStart > 2) {
            pageNumbers.push('...');
        }

        for (let i = rangeStart; i <= rangeEnd; i++) {
            if (i !== 1 && i !== totalPages) {
                pageNumbers.push(i);
            }
        }

        if (rangeEnd < totalPages - 1) {
            pageNumbers.push('...');
        }

        if (totalPages !== 1) {
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
                        key={index} 
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
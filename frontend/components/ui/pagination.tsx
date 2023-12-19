import React from 'react';

interface PaginationProps {
    currentPage: number;
    totalPages: number;
    onPageChange: (page: number) => void;
}

export const Pagination: React.FC<PaginationProps> = ({ currentPage, totalPages, onPageChange }) => {
    const renderPageNumbers = () => {
      const pageNumbers: (number | string)[] = [];
      for (let i = 1; i <= Math.min(10, totalPages); i++) {
        pageNumbers.push(i);
      }
  
      if (totalPages > 10) {
        if (currentPage > 6) {
          pageNumbers.splice(1, currentPage - 5, '...');
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
          page === '...' ? 
            <span key={index} style={{ margin: '0 5px' }}>...</span> : 
            <button 
              key={page} 
              onClick={() => onPageChange(typeof page === 'number' ? page : currentPage + 5)} 
              disabled={currentPage === page}
              style={{ margin: '0 5px' }}
            >
              {page}
            </button>
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
  
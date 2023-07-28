import React from 'react';

const OpenInColab = ({ link }) => (
    <a target="_blank" rel="noopener noreferrer" href={link}>
        <img src="https://colab.research.google.com/assets/colab-badge.svg" alt="Open In Colab" />
    </a>
);

export default OpenInColab;
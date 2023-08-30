'use client'

import React, { useState, useEffect } from 'react';

export default function InitJob() {
    interface Tool {
        CID: string;
        ToolJSON: string;
    }

    interface DataFile {
        ID: number;
        CID: string;
        WalletAddress: string;
        Filename: string;
        Timestamp: Date;
        Public: boolean;
        Visible: boolean; 
    }

    const [tools, setTools] = useState([]);
    const [dataFiles, setDataFiles] = useState([]);

    useEffect(() => {
        fetch('http://localhost:8080/get-tools')
            .then(response => {
                if (!response.ok) {
                    throw new Error(`HTTP error ${response.status}`);
                }
                return response.json();
            })
            .then(data => {
                setTools(data);
            })
            .catch(error => {
                console.error('There was an error!', error);
            });

        fetch('http://localhost:8080/get-datafiles')
            .then(response => {
                if (!response.ok) {
                    throw new Error(`HTTP error ${response.status}`);
                }
                return response.json();
            })
            .then(data => {
                setDataFiles(data);
            })
            .catch(error => {
                console.error('There was an error!', error);
            });
    }, []);

    return (
        <div>
            <h1>Initialize a Job</h1>
            <p>Choose a tool and data files to initialize a job.</p>
            <form>
                <label>
                    Select a tool:
                    <select>
                        {tools.map((tool, index) => {
                            const toolData = JSON.parse(tool.ToolJSON);
                            return (
                                <option key={index} value={tool.CID}>{toolData.name}</option>
                            );
                        })}
                    </select>
                </label>
                <label>
                    Select data files:
                    <select multiple>
                        {dataFiles.map((file, index) => (
                            // Replace 'fileProperty' with the actual property you want to display
                            <option key={index} value={file.CID}>{file.Filename}</option>
                        ))}
                    </select>
                </label>
                <input type="submit" value="Submit" />
            </form>
        </div>
    )
}
'use client'

import React, { useState, useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';

import backendUrl from 'lib/backendUrl';


export default function AddGraph() {
    const dispatch = useDispatch();

    interface Tool {
      Name: string
      CID: string
      ToolJson: { key: string, value: any }
    }

    interface DataFile {
        ID: number;
        CID: string;
        WalletAddress: string;
        Filename: string;
        Timestamp: Date;
    }

    const handleSelectToolChange = (value: string) => {
      const json_data = JSON.parse(value)
      console.log(json_data)
      console.log('Selected tool:', value)
    }

    const [tools, setTools] = useState<Tool[]>([]);
    const [selectedToolIndex, setSelectedToolIndex] = useState('')
    const [dataFiles, setDataFiles] = useState<DataFile[]>([]);
    const [selectedDataFiles, setSelectedDataFiles] = useState<string[]>([]);

    useEffect(() => {
        fetch(`${backendUrl()}/get-tools`)
            .then(response => response.json())
            .then(data => setTools(data))
            .catch(error => console.error('Error fetching tools:', error));

        fetch(`${backendUrl()}/get-datafiles`)
            .then(response => response.json())
            .then(data => setDataFiles(data))
            .catch(error => console.error('Error fetching data files:', error));
    }, []);

    const handleSubmit = (event: any) => {
        event.preventDefault()
        console.log(selectedToolIndex)
        console.log(tools[parseInt(selectedToolIndex)])

        const data = {
            tool: selectedTool,
            inputs: selectedDataFiles,
            // scatteringMethod: "dotProduct"
        };


        fetch(`${backendUrl()}/init-job`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
        })
        .then(response => response.json())
        .then(data => console.log('Job initialized:', data))
        .catch((error) => console.error('Error initializing job:', error));
    };

    return (
        <div>
            <h1>Initialize a Job</h1>
            <p>Choose a tool and data files to initialize a job.</p>
            <form onSubmit={handleSubmit}>
                <label>
                    Select a tool:
                    <select onChange={e => setSelectedToolIndex(e.target.value)}>
                        {tools.map((tool, index) => {
                            return (
                                <option key={index} value={index}>{tool.Name}</option>
                            );
                        })}
                    </select>
                </label>
                <label>
                    Select data files:
                    <select multiple onChange={e => setSelectedDataFiles(Array.from(e.target.selectedOptions, option => option.value))}>
                        {dataFiles.map((file, index) => (
                            <option key={index} value={file.Filename}>{file.Filename}</option>
                        ))}
                    </select>
                </label>
                <input type="submit" value="Submit" />
            </form>
        </div>
    )
}

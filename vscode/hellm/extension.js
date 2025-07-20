const vscode = require('vscode');
const { spawn } = require('child_process');
const path = require('path');
const fs = require('fs');

/**
 * Formats a HeLLM file using the hellm parse command
 * @param {vscode.TextDocument} document - The document to format
 * @returns {Promise<string>} - The formatted content
 */
function formatDocument(document) {
    return new Promise((resolve, reject) => {
        const config = vscode.workspace.getConfiguration('hellm');
        const hellmPath = config.get('hellmPath', 'hellm');
        
        // Use the document's file path directly
        const filePath = document.uri.fsPath;
        
        console.log(`Formatting file: ${filePath}`);
        
        const hellmProcess = spawn(hellmPath, ['parse', filePath], {
            cwd: vscode.workspace.workspaceFolders?.[0]?.uri?.fsPath || process.cwd()
        });

        let stdout = '';
        let stderr = '';

        hellmProcess.stdout.on('data', (data) => {
            stdout += data.toString();
        });

        hellmProcess.stderr.on('data', (data) => {
            stderr += data.toString();
        });

        hellmProcess.on('close', (code) => {
            if (code === 0) {
                console.log('File formatted successfully');
                resolve(stdout);
            } else {
                const error = `hellm parse failed with code ${code}: ${stderr}`;
                console.error(error);
                vscode.window.showErrorMessage(`HeLLM format error: ${stderr || 'Unknown error'}`);
                reject(new Error(error));
            }
        });

        hellmProcess.on('error', (err) => {
            const error = `Failed to execute hellm: ${err.message}`;
            console.error(error);
            vscode.window.showErrorMessage(`HeLLM format error: ${err.message}`);
            reject(err);
        });
    });
}

/**
 * @param {vscode.ExtensionContext} context
 */
function activate(context) {
    console.log('HeLLM extension is now active');

    // Register document formatting provider for built-in Format Document command
    const documentFormattingProvider = vscode.languages.registerDocumentFormattingEditProvider(
        [
            { language: 'hellm' },
            { scheme: 'file', language: 'hellm' },
            { pattern: '**/*.hl' }
        ],
        {
            async provideDocumentFormattingEdits(document) {
                try {
                    console.log('Format Document requested for:', document.uri.fsPath);
                    console.log('Document language ID:', document.languageId);
                    
                    // Save the document first if it has unsaved changes
                    if (document.isDirty) {
                        console.log('Document has unsaved changes, saving first...');
                        await document.save();
                    }

                    const formatted = await formatDocument(document);
                    const fullRange = new vscode.Range(
                        document.positionAt(0),
                        document.positionAt(document.getText().length)
                    );
                    
                    console.log('Returning formatting edits');
                    return [vscode.TextEdit.replace(fullRange, formatted)];
                } catch (error) {
                    console.error('Format error:', error);
                    vscode.window.showErrorMessage(`Formatting failed: ${error.message}`);
                    return [];
                }
            }
        }
    );

    console.log('Document formatting provider registered');

    context.subscriptions.push(documentFormattingProvider);
}

function deactivate() {}

module.exports = {
    activate,
    deactivate
}; 
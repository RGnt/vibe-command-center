import tokyoNightTheme from './tokyo-night.json';
import initEditor from 'monaco-mermaid';

let initialized = false;

export const setupMonaco = (monaco: any) => {
    if (initialized) return;
    
    // Define Tokyo Night theme
    monaco.editor.defineTheme('tokyo-night', tokyoNightTheme);
    
    // Initialize mermaid syntax highlighting
    initEditor(monaco);
    
    initialized = true;
};

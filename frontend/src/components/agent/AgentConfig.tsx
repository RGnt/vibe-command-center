import React, { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { dashboardWikisQueryOptions, dashboardTodosQueryOptions } from '../../utils/queries';
import { PatchDiff } from '@pierre/diffs/react';
import { PanelGroup, Panel, PanelResizeHandle } from 'react-resizable-panels';
import { ChevronDown, ChevronRight, ListCollapse, Expand } from 'lucide-react';

interface Diff {
    file: string;
    diff: string;
}

const AgentConfig = () => {
    const { data: wikis } = useQuery(dashboardWikisQueryOptions);
    const { data: todos } = useQuery(dashboardTodosQueryOptions);

    const [selectedWikiId, setSelectedWikiId] = useState<string>('');
    const [selectedTodoId, setSelectedTodoId] = useState<string>('');
    const [additionalPrompt, setAdditionalPrompt] = useState<string>('');
    const [isLoading, setIsLoading] = useState(false);
    const [response, setResponse] = useState<string | null>(null);
    const [diffs, setDiffs] = useState<Diff[]>([]);
    const [error, setError] = useState<string | null>(null);
    const [collapsedFiles, setCollapsedFiles] = useState<Set<string>>(new Set());

    const handleRunAgent = async () => {
        if (!selectedWikiId || !selectedTodoId) {
            setError("Please select both a Wiki and a Task.");
            return;
        }

        const selectedWiki = wikis?.find(w => w.id.toString() === selectedWikiId);
        const selectedTodo = todos?.find(t => t.id.toString() === selectedTodoId);

        if (!selectedWiki || !selectedTodo) return;

        setIsLoading(true);
        setError(null);
        setResponse(null);
        setDiffs([]);
        setCollapsedFiles(new Set()); // Expand all on new run

        try {
            const res = await fetch('/api/v1/agents/run', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    wiki_prompt: selectedWiki.content || "You are a helpful assistant.",
                    task_description: `Task: ${selectedTodo.title}\nDescription: ${selectedTodo.description || ''}`,
                    additional_prompt: additionalPrompt
                })
            });

            if (!res.ok) {
                const errData = await res.json();
                throw new Error(errData.detail || "Failed to run agent.");
            }

            const data = await res.json();
            setResponse(data.agent_response);
            if (data.changes && Array.isArray(data.changes)) {
                setDiffs(data.changes);
            }
        } catch (err: any) {
            setError(err.message);
        } finally {
            setIsLoading(false);
        }
    };

    const toggleFileCollapse = (filename: string) => {
        setCollapsedFiles(prev => {
            const next = new Set(prev);
            if (next.has(filename)) {
                next.delete(filename);
            } else {
                next.add(filename);
            }
            return next;
        });
    };

    const toggleAll = () => {
        if (collapsedFiles.size === diffs.length && diffs.length > 0) {
            // Expand all
            setCollapsedFiles(new Set());
        } else {
            // Collapse all
            setCollapsedFiles(new Set(diffs.map(d => d.file)));
        }
    };

    const isAllCollapsed = diffs.length > 0 && collapsedFiles.size === diffs.length;

    return (
        <div className="h-full w-full overflow-hidden flex flex-col">
            <div className="p-6 pb-2 shrink-0">
                <h1 className="text-3xl font-bold tracking-tight text-text-primary">Agent Configuration</h1>
            </div>

            <div className="flex-1 overflow-hidden p-6 pt-2">
                <PanelGroup direction="horizontal" className="h-full rounded-xl border border-border-base bg-bg-base/50">
                    
                    {/* Left Panel: Context & Task + Agent Response */}
                    <Panel defaultSize={40} minSize={25} className="flex flex-col h-full overflow-y-auto p-6 gap-6 custom-scrollbar">
                        <div className="bg-bg-card border border-border-base rounded-xl p-6 shadow-sm flex flex-col gap-4">
                            <h2 className="text-xl font-semibold text-text-primary">Context & Task</h2>
                            
                            <div className="flex flex-col gap-2">
                                <label className="text-sm font-medium text-text-muted">Wiki Prompt (Instructions)</label>
                                <select 
                                    className="p-3 rounded-lg border border-border-base bg-bg-base text-text-primary focus:outline-none focus:ring-2 focus:ring-primary/50"
                                    value={selectedWikiId}
                                    onChange={(e) => setSelectedWikiId(e.target.value)}
                                >
                                    <option value="">-- Select a Wiki --</option>
                                    {wikis?.map(wiki => (
                                        <option key={wiki.id} value={wiki.id}>{wiki.title}</option>
                                    ))}
                                </select>
                            </div>

                            <div className="flex flex-col gap-2">
                                <label className="text-sm font-medium text-text-muted">Task Description</label>
                                <select 
                                    className="p-3 rounded-lg border border-border-base bg-bg-base text-text-primary focus:outline-none focus:ring-2 focus:ring-primary/50"
                                    value={selectedTodoId}
                                    onChange={(e) => setSelectedTodoId(e.target.value)}
                                >
                                    <option value="">-- Select a Task --</option>
                                    {todos?.map(todo => (
                                        <option key={todo.id} value={todo.id}>{todo.title}</option>
                                    ))}
                                </select>
                            </div>

                            <div className="flex flex-col gap-2">
                                <label className="text-sm font-medium text-text-muted">Additional Prompt (Optional)</label>
                                <textarea 
                                    className="p-3 rounded-lg border border-border-base bg-bg-base text-text-primary focus:outline-none focus:ring-2 focus:ring-primary/50 min-h-[120px] resize-y"
                                    placeholder="Add any specific instructions for this run..."
                                    value={additionalPrompt}
                                    onChange={(e) => setAdditionalPrompt(e.target.value)}
                                />
                            </div>

                            <div className="mt-2 flex justify-end">
                                <button 
                                    onClick={handleRunAgent}
                                    disabled={isLoading}
                                    className="px-6 py-2.5 bg-primary text-primary-foreground font-semibold rounded-lg hover:bg-primary-hover disabled:opacity-50 disabled:cursor-not-allowed transition-all shadow-sm"
                                >
                                    {isLoading ? (
                                        <div className="flex items-center gap-2">
                                            <span className="animate-spin inline-block w-4 h-4 border-2 border-current border-t-transparent rounded-full" />
                                            Running Agent...
                                        </div>
                                    ) : 'Run Agent'}
                                </button>
                            </div>

                            {error && (
                                <div className="p-4 mt-2 bg-red-500/10 border border-red-500/20 rounded-lg text-red-500">
                                    {error}
                                </div>
                            )}
                        </div>

                        {response && (
                            <div className="bg-bg-card border border-border-base rounded-xl p-6 shadow-sm flex flex-col gap-4">
                                <h2 className="text-xl font-semibold text-text-primary">Agent Response</h2>
                                <div className="p-4 bg-bg-base rounded-lg border border-border-base whitespace-pre-wrap font-mono text-sm text-text-primary overflow-y-auto">
                                    {response}
                                </div>
                            </div>
                        )}
                    </Panel>

                    {/* Draggable Divider */}
                    <PanelResizeHandle className="w-2 bg-border-base hover:bg-primary/50 active:bg-primary transition-colors flex items-center justify-center cursor-col-resize relative group">
                        <div className="h-8 w-1 bg-border-strong rounded-full group-hover:bg-primary-foreground transition-colors" />
                    </PanelResizeHandle>

                    {/* Right Panel: Code Changes */}
                    <Panel defaultSize={60} minSize={25} className="flex flex-col h-full bg-bg-base">
                        <div className="p-6 h-full flex flex-col overflow-hidden">
                            <div className="flex items-center justify-between mb-4 shrink-0">
                                <h2 className="text-2xl font-bold text-text-primary flex items-center gap-2">
                                    Code Changes
                                    <span className="text-sm font-normal text-text-muted bg-bg-card px-2 py-0.5 rounded-full">
                                        {diffs.length} files
                                    </span>
                                </h2>
                                
                                {diffs.length > 0 && (
                                    <button 
                                        onClick={toggleAll}
                                        className="flex items-center gap-2 px-3 py-1.5 text-sm font-medium text-text-muted hover:text-text-primary bg-bg-card hover:bg-bg-hover border border-border-base rounded-lg transition-colors"
                                    >
                                        {isAllCollapsed ? (
                                            <><Expand className="w-4 h-4" /> Expand All</>
                                        ) : (
                                            <><ListCollapse className="w-4 h-4" /> Collapse All</>
                                        )}
                                    </button>
                                )}
                            </div>

                            <div className="flex-1 overflow-y-auto pr-2 custom-scrollbar flex flex-col gap-4">
                                {diffs.length === 0 ? (
                                    <div className="h-full flex items-center justify-center text-text-muted">
                                        No code changes to display. Run the agent to generate changes.
                                    </div>
                                ) : (
                                    diffs.map((diffItem, idx) => {
                                        const isCollapsed = collapsedFiles.has(diffItem.file);
                                        return (
                                            <div key={idx} className="bg-bg-card border border-border-base rounded-xl overflow-hidden shadow-sm shrink-0">
                                                <button 
                                                    onClick={() => toggleFileCollapse(diffItem.file)}
                                                    className="w-full px-4 py-3 bg-bg-base border-b border-border-base flex items-center gap-3 hover:bg-bg-hover transition-colors text-left"
                                                >
                                                    <span className="text-text-muted">
                                                        {isCollapsed ? <ChevronRight className="w-5 h-5" /> : <ChevronDown className="w-5 h-5" />}
                                                    </span>
                                                    <span className="font-mono text-sm font-semibold text-text-primary truncate">
                                                        {diffItem.file}
                                                    </span>
                                                </button>
                                                
                                                {!isCollapsed && (
                                                    <div className="overflow-x-auto">
                                                        <PatchDiff patch={diffItem.diff} />
                                                    </div>
                                                )}
                                            </div>
                                        );
                                    })
                                )}
                            </div>
                        </div>
                    </Panel>

                </PanelGroup>
            </div>
            
            <style>{`
                .custom-scrollbar::-webkit-scrollbar {
                    width: 8px;
                    height: 8px;
                }
                .custom-scrollbar::-webkit-scrollbar-track {
                    background: transparent;
                }
                .custom-scrollbar::-webkit-scrollbar-thumb {
                    background-color: rgba(156, 163, 175, 0.3);
                    border-radius: 20px;
                }
                .custom-scrollbar::-webkit-scrollbar-thumb:hover {
                    background-color: rgba(156, 163, 175, 0.5);
                }
            `}</style>
        </div>
    );
};

export default AgentConfig;

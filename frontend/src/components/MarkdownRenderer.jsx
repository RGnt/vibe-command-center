import React, { useEffect, useRef } from 'react';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import remarkMath from 'remark-math';
import rehypeKatex from 'rehype-katex';
import mermaid from 'mermaid';

mermaid.initialize({
    startOnLoad: true,
    theme: 'dark',
    securityLevel: 'loose',
    fontFamily: 'Inter, sans-serif'
});

const Mermaid = ({ chart }) => {
    const containerRef = useRef(null);

    useEffect(() => {
        if (containerRef.current) {
            // Re-render when chart changes
            containerRef.current.removeAttribute('data-processed');
            containerRef.current.innerHTML = chart;
            mermaid.run({
                nodes: [containerRef.current]
            }).catch(e => console.error('Mermaid rendering failed', e));
        }
    }, [chart]);

    return (
        <div ref={containerRef} className="mermaid flex justify-center my-6 bg-bg-base/50 p-4 rounded-lg border border-border">
            {chart}
        </div>
    );
};

const DiagramEmbed = ({ id }) => {
    const [diagram, setDiagram] = React.useState(null);
    const [loading, setLoading] = React.useState(true);

    React.useEffect(() => {
        fetch(`/api/diagrams/${id}`)
            .then(res => res.ok ? res.json() : null)
            .then(data => {
                setDiagram(data);
                setLoading(false);
            })
            .catch(() => setLoading(false));
    }, [id]);

    if (loading) return <div className="p-4 text-center text-text-muted">Loading diagram {id}...</div>;
    if (!diagram) return <div className="p-4 text-center text-red-500">Diagram not found</div>;

    return (
        <div className="my-6 border border-border rounded-lg overflow-hidden bg-bg-base/30">
            <div className="p-2 bg-bg-subtle border-b border-border text-sm font-semibold text-text-base">
                {diagram.name}
            </div>
            <div className="p-4">
                <Mermaid chart={diagram.code} />
            </div>
            {diagram.explanation && (
                <div className="p-4 border-t border-border bg-bg-base text-text-muted text-sm italic">
                    {diagram.explanation}
                </div>
            )}
        </div>
    );
};

const MarkdownRenderer = ({ content }) => {
    // Pre-process shortcodes: {{diagram:123}} -> [Diagram 123](#diagram:123)
    const processedContent = (content || '').replace(/\{\{diagram:(\d+)\}\}/g, '[Diagram $1](#diagram:$1)');

    return (
        <div className="prose prose-sm dark:prose-invert max-w-none 
                        prose-headings:text-text-base prose-p:text-text-muted 
                        prose-a:text-primary prose-a:no-underline hover:prose-a:underline
                        prose-strong:text-text-base prose-strong:font-semibold
                        prose-code:text-accent prose-code:bg-bg-base prose-code:px-1 prose-code:rounded
                        prose-pre:bg-bg-base prose-pre:border prose-pre:border-border
                        prose-th:text-text-base prose-th:bg-bg-base prose-td:text-text-muted
                        prose-blockquote:border-primary prose-blockquote:text-text-muted prose-blockquote:bg-bg-base/50 prose-blockquote:px-4 prose-blockquote:py-1 prose-blockquote:rounded-r
                        prose-li:text-text-muted
                        prose-img:rounded-lg prose-img:shadow-sm">
            <ReactMarkdown
                remarkPlugins={[remarkGfm, remarkMath]}
                rehypePlugins={[rehypeKatex]}
                components={{
                    code({ node, inline, className, children, ...props }) {
                        const match = /language-(\w+)/.exec(className || '');
                        const isMermaid = match && match[1] === 'mermaid';

                        if (!inline && isMermaid) {
                            return <Mermaid chart={String(children).replace(/\n$/, '')} />;
                        }

                        return !inline ? (
                            <div className="overflow-x-auto my-2">
                                <code className={className} {...props}>
                                    {children}
                                </code>
                            </div>
                        ) : (
                            <code className={className} {...props}>
                                {children}
                            </code>
                        );
                    },
                    a: ({ node, href, children, ...props }) => {
                        if (href?.startsWith('#diagram:')) {
                            const id = href.replace('#diagram:', '');
                            return <DiagramEmbed id={id} />;
                        }
                        if (href?.startsWith('#wiki:')) {
                            return (
                                <a
                                    href={href}
                                    onClick={(e) => {
                                        e.preventDefault();
                                        const slug = href.replace('#wiki:', '');
                                        const event = new CustomEvent('wikiNavigate', { detail: { slug } });
                                        window.dispatchEvent(event);
                                    }}
                                    className="text-primary hover:underline font-medium cursor-pointer"
                                    {...props}
                                >
                                    {children}
                                </a>
                            );
                        }
                        return (
                            <a href={href} target="_blank" rel="noopener noreferrer" {...props}>
                                {children}
                            </a>
                        );
                    }
                }}
            >
                {processedContent}
            </ReactMarkdown>
        </div>
    );
};

export default MarkdownRenderer;

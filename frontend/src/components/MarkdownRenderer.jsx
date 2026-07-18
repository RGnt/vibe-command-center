import React from 'react';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import remarkMath from 'remark-math';
import rehypeKatex from 'rehype-katex';

const MarkdownRenderer = ({ content }) => {
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
                    a: ({ node, href, children, ...props }) => {
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
                {content || ''}
            </ReactMarkdown>
        </div>
    );
};

export default MarkdownRenderer;

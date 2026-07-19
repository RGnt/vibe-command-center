import { queryOptions } from '@tanstack/react-query';

// Projects
export const projectsQueryOptions = queryOptions({
    queryKey: ['projects'],
    queryFn: async () => {
        const res = await fetch('/api/projects');
        if (res.status === 401) return [];
        if (!res.ok) throw new Error('Failed to fetch projects');
        return res.json();
    }
});

// Dashboard specific global todos
export const dashboardTodosQueryOptions = queryOptions({
    queryKey: ['dashboard-todos'],
    queryFn: async () => {
        const res = await fetch('/api/todos');
        if (res.status === 401) return [];
        if (!res.ok) throw new Error('Failed to fetch dashboard todos');
        const data = await res.json();
        // Return only open tasks
        const openTasks = (data || []).filter(t => !t.completed && t.stage !== 'Done');
        return openTasks.slice(0, 10);
    }
});

// Dashboard specific global wikis
export const dashboardWikisQueryOptions = queryOptions({
    queryKey: ['dashboard-wikis'],
    queryFn: async () => {
        const res = await fetch('/api/wikis');
        if (res.status === 401) return [];
        if (!res.ok) throw new Error('Failed to fetch dashboard wikis');
        const data = await res.json();
        return (data || []).slice(0, 5);
    }
});

// Specific Project
export const projectQueryOptions = (projectId) => queryOptions({
    queryKey: ['projects', projectId],
    queryFn: async () => {
        const res = await fetch(`/api/projects/${projectId}`);
        if (!res.ok) throw new Error('Failed to fetch project');
        return res.json();
    }
});

// Specific Project Todos
export const projectTodosQueryOptions = (projectId) => queryOptions({
    queryKey: ['projects', projectId, 'todos'],
    queryFn: async () => {
        const res = await fetch(`/api/todos?project_id=${projectId}`);
        if (!res.ok) throw new Error('Failed to fetch project todos');
        return res.json();
    }
});

// Global Wikis (All)
export const globalWikisQueryOptions = queryOptions({
    queryKey: ['wikis'],
    queryFn: async () => {
        const res = await fetch('/api/wikis');
        if (res.status === 401) return [];
        if (!res.ok) throw new Error('Failed to fetch wikis');
        return res.json();
    }
});

// Specific Project Wikis
export const projectWikisQueryOptions = (projectId) => queryOptions({
    queryKey: ['projects', projectId, 'wikis'],
    queryFn: async () => {
        const res = await fetch(`/api/wikis?project_id=${projectId}`);
        if (!res.ok) throw new Error('Failed to fetch project wikis');
        return res.json();
    }
});

// Workflows
export const workflowsQueryOptions = queryOptions({
    queryKey: ['workflows'],
    queryFn: async () => {
        const res = await fetch('/api/workflows');
        if (res.status === 401) return [];
        if (!res.ok) throw new Error('Failed to fetch workflows');
        return res.json();
    }
});

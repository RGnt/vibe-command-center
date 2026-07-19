import React from 'react';
import { createFileRoute } from '@tanstack/react-router';
import { useSuspenseQuery } from '@tanstack/react-query';
import KanbanBoard from '../../../components/kanban/KanbanBoard';
import { projectQueryOptions, projectTodosQueryOptions, workflowsQueryOptions } from '../../../utils/queries';
import { useMutation, useQueryClient } from '@tanstack/react-query';

export const Route = createFileRoute('/projects/$projectId/')({
  loader: async ({ params, context }) => {
    await Promise.all([
      context.queryClient.ensureQueryData(projectQueryOptions(params.projectId)),
      context.queryClient.ensureQueryData(projectTodosQueryOptions(params.projectId)),
      context.queryClient.ensureQueryData(workflowsQueryOptions),
    ]);
  },
  component: ProjectKanbanRoute,
});

function ProjectKanbanRoute() {
  const { projectId } = Route.useParams();
  const queryClient = useQueryClient();
  
  const { data: project } = useSuspenseQuery(projectQueryOptions(projectId));
  const { data: todos } = useSuspenseQuery(projectTodosQueryOptions(projectId));
  const { data: workflows } = useSuspenseQuery(workflowsQueryOptions);

  const createTodoMutation = useMutation({
    mutationFn: async (todo) => {
      const response = await fetch('/api/todos', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ ...todo, project_id: parseInt(projectId) }),
      });
      if (!response.ok) throw new Error('Failed to create todo');
      return response.json();
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['projects', projectId, 'todos'] });
      queryClient.invalidateQueries({ queryKey: ['dashboard-todos'] });
    },
  });

  const updateTodoMutation = useMutation({
    mutationFn: async ({ id, updates }) => {
      const response = await fetch(`/api/todos/${id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(updates),
      });
      if (!response.ok) throw new Error('Failed to update todo');
      return response.json();
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['projects', projectId, 'todos'] });
      queryClient.invalidateQueries({ queryKey: ['dashboard-todos'] });
    },
  });

  const deleteTodoMutation = useMutation({
    mutationFn: async (id) => {
      const response = await fetch(`/api/todos/${id}`, { method: 'DELETE' });
      if (!response.ok) throw new Error('Failed to delete todo');
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['projects', projectId, 'todos'] });
      queryClient.invalidateQueries({ queryKey: ['dashboard-todos'] });
    },
  });

  const createSubtaskMutation = useMutation({
    mutationFn: async (todo) => {
      const response = await fetch('/api/todos', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ ...todo, project_id: parseInt(projectId) }),
      });
      if (!response.ok) throw new Error('Failed to create subtask');
      return response.json();
    },
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['projects', projectId, 'todos'] }),
  });

  return (
    <div className="flex-1 p-6 h-full overflow-hidden flex flex-col">
        <h1 className="text-2xl font-bold mb-4">{project.name} Board</h1>
        <KanbanBoard 
            activeProjectWorkflowId={project.workflow_id}
            todos={todos}
            workflows={workflows}
            onTodoCreate={(todo) => createTodoMutation.mutate(todo)}
            onTodoUpdate={(id, updates) => updateTodoMutation.mutate({ id, updates })}
            onTodoDelete={(id) => deleteTodoMutation.mutate(id)}
            onAddSubtask={(todo) => createSubtaskMutation.mutate(todo)}
        />
    </div>
  );
}

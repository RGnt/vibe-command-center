import React from 'react';
import { createFileRoute, useNavigate } from '@tanstack/react-router';
import { useSuspenseQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import Dashboard from '../components/dashboard/Dashboard';
import { projectsQueryOptions, dashboardTodosQueryOptions, dashboardWikisQueryOptions } from '../utils/queries';

export const Route = createFileRoute('/')({
  loader: async ({ context }) => {
    // Ensure all required dashboard data is fetched before render
    await Promise.all([
      context.queryClient.ensureQueryData(projectsQueryOptions),
      context.queryClient.ensureQueryData(dashboardTodosQueryOptions),
      context.queryClient.ensureQueryData(dashboardWikisQueryOptions),
    ]);
  },
  component: DashboardRoute,
  pendingComponent: () => (
    <div className="flex-1 p-8 bg-bg-base/50 flex justify-center items-center h-full">
        <div className="text-text-muted animate-pulse">Loading dashboard...</div>
    </div>
  )
});

function DashboardRoute() {
  const queryClient = useQueryClient();

  // Data is guaranteed to be available here without isLoading checks
  const { data: projects } = useSuspenseQuery(projectsQueryOptions);
  const { data: todos } = useSuspenseQuery(dashboardTodosQueryOptions);
  const { data: wikis } = useSuspenseQuery(dashboardWikisQueryOptions);

  const deleteProjectMutation = useMutation({
    mutationFn: async (id) => {
      const res = await fetch(`/api/projects/${id}`, { method: 'DELETE' });
      if (!res.ok) throw new Error('Failed to delete');
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: projectsQueryOptions.queryKey });
    }
  });

  const navigate = useNavigate();
  const handleSelectProject = (project) => {
      navigate({ to: `/projects/${project.id}` });
  };

  const handleDeleteProject = (id) => {
      deleteProjectMutation.mutate(id);
  };

  return (
    <Dashboard 
        projects={projects}
        todos={todos}
        wikis={wikis}
        fetchProjects={() => queryClient.invalidateQueries({ queryKey: projectsQueryOptions.queryKey })}
        onSelectProject={handleSelectProject}
        onDeleteProject={handleDeleteProject}
    />
  );
}

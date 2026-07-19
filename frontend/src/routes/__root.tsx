import React, { useState } from 'react';
import { createRootRoute, Outlet, useNavigate } from '@tanstack/react-router';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import Sidebar from '../components/layout/Sidebar';
import ProjectModal from '../components/layout/ProjectModal';
import { useTheme } from '../contexts/ThemeContext';
import { useAuth } from '../contexts/AuthContext';
import Login from '../components/auth/Login';

import { projectsQueryOptions } from '../utils/queries';

export const Route = createRootRoute({
  loader: async ({ context }) => {
    if (context.auth?.user) {
        await context.queryClient.ensureQueryData(projectsQueryOptions);
    }
  },
  component: RootComponent,
});

function RootComponent() {
  const { user } = useAuth();
  const { theme, toggleTheme } = useTheme();
  const [isProjectModalOpen, setIsProjectModalOpen] = useState(false);
  const queryClient = useQueryClient();
  const navigate = useNavigate();

  const createProjectMutation = useMutation({
    mutationFn: async (projectData) => {
        const res = await fetch('/api/projects', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(projectData)
        });
        if (!res.ok) throw new Error('Failed to create project');
        return res.json();
    },
    onSuccess: (newProject) => {
        queryClient.invalidateQueries({ queryKey: projectsQueryOptions.queryKey });
        setIsProjectModalOpen(false);
        navigate({ to: `/projects/${newProject.id}` });
    }
  });

  const handleImportProject = async (e) => {
    const file = e.target.files[0];
    if (!file) return;

    const reader = new FileReader();
    reader.onload = async (e) => {
        try {
            const content = JSON.parse(e.target.result);
            const res = await fetch('/api/projects/import', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(content)
            });
            if (res.ok) {
                queryClient.invalidateQueries({ queryKey: projectsQueryOptions.queryKey });
            } else {
                alert("Failed to import project");
            }
        } catch (err) {
            console.error("Import error", err);
            alert("Invalid JSON file");
        }
    };
    reader.readAsText(file);
    e.target.value = ''; // Reset input
  };

  if (!user) {
    return (
      <div className="min-h-screen bg-bg-base transition-colors duration-300">
        <Login />
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-bg-base transition-colors duration-300">
      <div className="w-full h-screen overflow-hidden flex bg-gradient-to-br from-bg-base via-bg-base to-bg-panel text-text-base">
        <Sidebar 
            isDark={theme === 'dark'}
            onToggleTheme={toggleTheme}
            onNewProject={() => setIsProjectModalOpen(true)}
            onImportProject={handleImportProject}
        />
        <main className="flex-1 flex flex-col h-full overflow-hidden">
            <Outlet />
        </main>
      </div>
      
      <ProjectModal 
          isOpen={isProjectModalOpen}
          onClose={() => setIsProjectModalOpen(false)}
          onSave={(data) => createProjectMutation.mutate(data)}
      />
    </div>
  );
}

import React from 'react';
import { createFileRoute } from '@tanstack/react-router';
import { useSuspenseQuery } from '@tanstack/react-query';
import WikiView from '../../../components/WikiView';
import { projectQueryOptions, projectWikisQueryOptions } from '../../../utils/queries';

export const Route = createFileRoute('/projects/$projectId/wiki')({
  loader: async ({ params, context }) => {
    await Promise.all([
      context.queryClient.ensureQueryData(projectQueryOptions(params.projectId)),
      context.queryClient.ensureQueryData(projectWikisQueryOptions(params.projectId)),
    ]);
  },
  component: ProjectWikiRoute,
});

function ProjectWikiRoute() {
  const { projectId } = Route.useParams();
  const { data: project } = useSuspenseQuery(projectQueryOptions(projectId));

  return (
    <div className="flex-1 p-6 h-full overflow-hidden">
        <WikiView project={project} isGlobal={false} />
    </div>
  );
}

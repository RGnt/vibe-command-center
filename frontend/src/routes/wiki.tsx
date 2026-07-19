import React from 'react';
import { createFileRoute } from '@tanstack/react-router';
import { useSuspenseQuery } from '@tanstack/react-query';
import WikiView from '../components/WikiView';
import { globalWikisQueryOptions } from '../utils/queries';

export const Route = createFileRoute('/wiki')({
  loader: async ({ context }) => {
    await context.queryClient.ensureQueryData(globalWikisQueryOptions);
  },
  component: GlobalWikiRoute,
});

function GlobalWikiRoute() {
  return (
    <div className="flex-1 p-6 h-full overflow-hidden">
        <WikiView isGlobal={true} />
    </div>
  );
}

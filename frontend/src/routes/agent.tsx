import { createFileRoute } from '@tanstack/react-router';
import AgentConfig from '../components/agent/AgentConfig';
import { dashboardTodosQueryOptions, dashboardWikisQueryOptions } from '../utils/queries';

export const Route = createFileRoute('/agent')({
  loader: async ({ context }) => {
    // Ensure data is available for the dropdowns
    await Promise.all([
      context.queryClient.ensureQueryData(dashboardTodosQueryOptions),
      context.queryClient.ensureQueryData(dashboardWikisQueryOptions),
    ]);
  },
  component: AgentConfig
});

/**
 * This files contains the titles and descriptions for the different sections of the site
 */
export default {
  index: {
    title: '{ appName }: Discover & launch great Kubernetes-ready apps',
    description:
      '{ appName } is a platform for discovering & launching great Kubernetes-ready' +
        'apps. Browse the catalog and deploy your applications in your Kubernetes cluster'
  },
  charts: {
    title: 'Kubernetes-ready applications catalog | { appName }',
    description:
      'Browse the { appName } catalog of Kubernetes-ready apps. Deploy all apps you need ' +
        'in your infrastructure or the cloud with a command using Helm Charts'
  },
  repoCharts: {
    title: '{ repo } repository of Kubernetes-ready applications | { appName }',
    description:
      'Browse the { appName } catalog of the { repo } repository of Kubernetes-ready apps. ' +
        'Deploy all apps you need in your infrastructure or the cloud with a command using ' +
        'Helm Charts'
  },
  chartDetails: {
    title: '{ name } for Kubernetes | { appName }',
    description: 'Deploy the latest { name } in Kubernetes. { description }'
  },
  chartDetailsWithVersion: {
    title: '{ name } { version } for Kubernetes | { appName }',
    description:
      'Deploy the { name } { version } in Kubernetes. { description }'
  },
  deployments: {
    title: 'Manage Deployments in Kubernetes | { appName }',
    description: 'Browse, edit or create deployments in Kubernetes.'
  },
  deploymentDetails: {
    title: '{ name } in Kubernetes | { appName }',
    description: 'Deployment { name } in Kubernetes. { description }'
  },
  repositories: {
    title: 'Repositories of Kubernetes-ready applications | { appName }',
    description: 'Manage repositories of helm charts'
  },
};

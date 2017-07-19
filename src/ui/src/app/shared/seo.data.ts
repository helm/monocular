/**
 * This files contains the titles and descriptions for the different sections of the site
 */
export default {
  index: {
    title: 'Discover & launch great Kubernetes-ready apps',
    description:
      '{ appName } is a platform for discovering & launching great Kubernetes-ready' +
        'apps. Browse the catalog and deploy your applications in your Kubernetes cluster'
  },
  search: {
    title: 'Results for "{ search }" Kubernetes-ready applications',
    description:
      'Results for "{ search }" in the { appName } catalog of Kubernetes-ready applications. ' +
        'Deploy all apps you need in your infrastructure or the cloud'
  },
  charts: {
    title: 'Kubernetes-ready applications catalog',
    description:
      'Browse the { appName } catalog of Kubernetes-ready apps. Deploy all apps you need ' +
        'in your infrastructure or the cloud with a command using Helm Charts'
  },
  repoCharts: {
    title: '{ repo } repository of Kubernetes-ready applications',
    description:
      'Browse the { appName } catalog of the { repo } repository of Kubernetes-ready apps. ' +
        'Deploy all apps you need in your infrastructure or the cloud with a command using ' +
        'Helm Charts'
  },
  details: {
    title: '{ name } for Kubernetes',
    description: 'Deploy the latest { name } in Kubernetes. { description }'
  },
  detailsWithVersion: {
    title: '{ name } { version } for Kubernetes',
    description:
      'Deploy the { name } { version } in Kubernetes. { description }'
  }
};

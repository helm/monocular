# ChartSvc

ChartSvc is a service for KubeApps that reads chart metadata from the database
and presents it in a RESTful API. It should be used with the
[chart-repo-sync](https://github.com/kubernetes-helm/monocular/tree/master/src/api/cmd/chart-repo-sync) to populate chart
metadata in the database.

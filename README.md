README
--

Suse Stycle DevOps CI/CD Demo.

This is the frontend microservice for displaying the projects by calling RESTful APi with backend https://github.com/bluezd/suse-demo-list-projects

## Architecture

![image](https://user-images.githubusercontent.com/977107/137465403-69686aa4-950c-4fe6-b6a1-c1901907d5ad.png)

## Deploy

  - helm chart
    - https://github.com/bluezd/suse-demo-frontend-chart
    - https://github.com/bluezd/suse-demo-list-projects-chart
  - change the `suse-demo-frontend-chart` `templates/deployment.yaml` `LIST_PROJECTS_ENDPOINT` environment variable accordingly.
  - helm install

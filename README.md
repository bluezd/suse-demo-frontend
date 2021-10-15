README
--

## docker build

```
TAG=`date '+%Y-%m-%d-%H-%M-%S'`
docker build -t bluezd/suse-demo-frontend:$TAG .
```

## docker run

```
docker run -d -p 32080:8000 -e LIST_PROJECTS_ENDPOINT='http://172.17.0.3:8001/projects' bluezd/suse-demo-frontend:$TAG
```

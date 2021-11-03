pipeline {  
//声明执行构建任务的 agent，指明该agnet以POD的形式运行在名称为“k8s”的kubernetes集群中，agent的标签是“jenkins-salve-pod”  
    agent {
        kubernetes {
            cloud 'kubernetes'
            label 'slave-for-backend'
            defaultContainer 'jnlp'
            yaml """
apiVersion: v1
kind: Pod
metadata:
  name: jenkins-salve
  labels:
    app: jenkins-slave
spec:
  securityContext:
    runAsUser: 0
    fsGroup: 0
  hostAliases:  
  - ip: 192.168.0.68
    hostnames:  
    - cloud.agilor.co
  - ip: 192.168.0.235
    hostnames:
    - gitlab.cloud.agilor.co
    - jenkins.cloud.agilor.co
  containers:
  - name: jnlp
    image: jenkins/inbound-agent:4.3-4
    imagePullPolicy: IfNotPresent
    lifecycle:
      postStart:
        exec:
          command: ['git', 'config', '--global', 'http.sslVerify', 'false']
  - name: docker
    image: docker:19.03
    imagePullPolicy: IfNotPresent
    command: ["/bin/sh","-c"]
    args: ["while true; do sleep 1; done;"]
    volumeMounts:
    - name: jenkinsdocker
      mountPath: /var/run/docker.sock
  - name: tool-image
    image: cnych/helm-kubectl-curl-git-jq-yq
    imagePullPolicy: IfNotPresent
    command: ["/bin/sh","-c"]
    args: ["while true; do sleep 1; done;"]
  volumes:
  - name: jenkinsdocker
    hostPath:
      path: /var/run/docker.sock
      type: Socket
  serviceAccount: jenkins-admin
    """
        }
    }  

    environment{
        PROJECT_NAME = "rancher-demo" 
        APP_NAME = "suse-demo-frontend"
        DOCKER_REGISTRY = "192.168.0.42"
        IMAGE_TAG="v1.0.$BUILD_ID"
        IMAGE_NAME_WITH_TAG =  "${DOCKER_REGISTRY}/${PROJECT_NAME}/${APP_NAME}:${IMAGE_TAG}"
        DEPLOY_GITREPO_TOKEN = credentials('gitlab-token')
    }
  
    stages {
        stage('Build Image') {
            steps {
                container('docker') {
                    withDockerRegistry(credentialsId: 'harbor-admin', url: 'http://192.168.0.42') {

                        sh 'docker build -t ${IMAGE_NAME_WITH_TAG} -f Dockerfile .'
                        sh 'docker push ${IMAGE_NAME_WITH_TAG}'
                    }
                }
            }
        }
        stage('GitOps Deploy') {
            steps {
              container('tool-image') {
                sh """
                  git config --global user.name gitlab
                  git config --global user.email gitlab@agilor
                  git config --global http.sslVerify false
                  git clone --branch main --depth 1 https://root:$env.DEPLOY_GITREPO_TOKEN@gitlab.cloud.agilor.co/root/suse-demo-frontend-chart.git repo
                  # After cloning
                  cd repo
                  ls -l
                  echo old value:
                  cat values.yaml | yq r - 'image.tag'
                  echo replacing with new value:
                  echo $env.IMAGE_TAG
                  yq w --inplace values.yaml 'image.tag' "$env.IMAGE_TAG"
                  git diff
                  if ! git diff-index --quiet HEAD --; then
                    git add .
                    git commit -m "helm values updated by tekton pipeline with tag $env.IMAGE_TAG"
                    git push origin main
                  else
                      echo "no changes, git repository is up to date"
                  fi
                """
              }
            }
        }
    }
}

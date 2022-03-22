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
        PROJECT_NAME = "neuvector-demo" 
        APP_NAME = "suse-demo-frontend"
        DOCKER_REGISTRY = "172.16.99.145:8080"
        IMAGE_TAG="v1.0.$BUILD_ID"
        IMAGE_NAME_WITH_TAG =  "${DOCKER_REGISTRY}/${PROJECT_NAME}/${APP_NAME}:${IMAGE_TAG}"
    }
  
    stages {
        stage('Build Image') {
            steps {
                container('docker') {
					sh 'docker build -t ${IMAGE_NAME_WITH_TAG} -f Dockerfile .'
                }
            }
        }

		stage('Image Vulnerability Scan By NeuVector') {
			steps {
				neuvector nameOfVulnerabilityToExemptFour: '', nameOfVulnerabilityToExemptOne: '', nameOfVulnerabilityToExemptThree: '', nameOfVulnerabilityToExemptTwo: '', nameOfVulnerabilityToFailFour: '', nameOfVulnerabilityToFailOne: '', nameOfVulnerabilityToFailThree: '', nameOfVulnerabilityToFailTwo: '', numberOfHighSeverityToFail: '', numberOfMediumSeverityToFail: '1', registrySelection: 'Local', repository: "${DOCKER_REGISTRY}/${PROJECT_NAME}/${APP_NAME}", scanLayers: true, scanTimeout: 10, tag: "v1.0.${env.BUILD_ID}"
			}
		}

        stage('Push Image') {
            steps {
                container('docker') {
                    withDockerRegistry(credentialsId: 'harbor-admin', url: 'http://172.16.99.145:8080') {
                        sh 'docker push ${IMAGE_NAME_WITH_TAG}'
                    }
                }
            }
        }

        stage('Generate Deploy Yamls') {
            steps {
				sh '''
				  sed -i "s/IMAGE_TAG/${IMAGE_TAG}/g" Template.yaml
				'''
            }
        }

        stage('Deploy to Kubernetes') {
            steps {
                container('jnlp') {
					kubernetesDeploy configs: 'Template.yaml', kubeConfig: [path: ''], kubeconfigId: 'k8s', secretName: '', ssh: [sshCredentialsId: '*', sshServer: ''], textCredentials: [certificateAuthorityData: '', clientCertificateData: '', clientKeyData: '', serverUrl: 'https://']
                }
            }
        }
    }
}

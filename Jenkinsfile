pipeline {
  agent any
  stages {
    stage('Initialize') {
      steps {
        echo 'this is a pipeline'
      }
    }

    stage('Build ') {
      steps {
        sh '''cd ../broker-service && env GOOS=linux CGO_ENABLED=0 go build -o brokerApp ./cmd/api
docker-compose up --build -d'''
      }
    }

  }
}
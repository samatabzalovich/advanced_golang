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
        sh '''@echo "Building broker binary..."
cd ../broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api
@echo "Done!"'''
      }
    }

  }
}
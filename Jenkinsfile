pipeline {
    agent none 
    stages {
        stage('Build') { 
            agent {
                image 'golang:1.15-alpine'
            }
            steps {
                sh '''
                go mod tidy
                go build -o  main.go
                sh ''' 
            }
        }
    }
}
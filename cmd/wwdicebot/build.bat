SET GOOS=linux
del wwdicebot
go build .
docker build -t klaital/wwdicebot:latest .
docker push klaital/wwdicebot:latest
kubectl --kubeconfig  C:\Users\kenka\.kube\config apply -f k8s.yaml

eval $(minikube docker-env)
docker build -t goods:0.1 -f services/goods/Dockerfile .
docker build -t orders:0.1 -f services/orders/Dockerfile .
docker build -t notify:0.1 -f services/notify/Dockerfile .
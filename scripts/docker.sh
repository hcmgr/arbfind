sudo docker build -t arb-backend .
sudo docker rm arb-backend
sudo docker run --name arb-backend -p 8080:8080 arb-backend
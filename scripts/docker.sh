sudo docker build -t arb-backend .
sudo docker stop arb-backend
sudo docker rm arb-backend
sudo docker run -d --name arb-backend -p 10001:10001 arb-backend
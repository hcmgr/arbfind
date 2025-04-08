### Dockers steps

```
## build image
docker pull golang:1.24.1
sudo docker build -t arb-backend .

## list images
docker images

## run
sudo docker run -p 8080:8080 arb-backend


```
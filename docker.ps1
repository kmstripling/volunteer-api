
docker network create volunteer-network

docker run --name mariadbtest -e MYSQL_ROOT_PASSWORD=mypass -p 3306:3306 -d docker.io/library/mariadb:latest


docker stop volunteerservice
docker rm volunteerservice
docker run --name volunteerservice -p 3000:3000 -d docker.io/library/code:latest

docker network create volunteer-network
docker network connect volunteer-network mariadbtest
docker network connect volunteer-network volunteerservice

docker network inspect volunteer-network
default:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o release/linux/amd64/lock9
	sudo docker build -t mingchoi/lock9 .
	sudo docker stop lock9
	sudo docker rm lock9
	sudo docker run \
	--name lock9 \
	--env LOCK9_API_SECRET=${LOCK9_API_SECRET} \
	--env LOCK9_DB_SECRET=${LOCK9_DB_SECRET} \
	--env LOCK9_DB_NAME=${LOCK9_DB_NAME} \
	--env LOCK9_JAPANGROUP_FORMURL=${LOCK9_JAPANGROUP_FORMURL} \
	--link mariadb:db \
	-d \
	mingchoi/lock9

build_linux_amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -o release/linux/amd64/lock9

build_image:
	sudo docker build -t mingchoi/lock9 .

run_container:
	sudo docker run \
	--name lock9 \
	--env LOCK9_API_SECRET=${LOCK9_API_SECRET} \
	--env LOCK9_DB_SECRET=${LOCK9_DB_SECRET} \
	--env LOCK9_DB_NAME=${LOCK9_DB_NAME} \
	--env LOCK9_JAPANGROUP_FORMURL=${LOCK9_JAPANGROUP_FORMURL} \
	--link mariadb:db \
	-d \
	lock9

rerun_container:
	sudo docker stop lock9
	sudo docker rm lock9
	sudo docker run \
	--name lock9 \
	--env LOCK9_API_SECRET=${LOCK9_API_SECRET} \
	--env LOCK9_DB_SECRET=${LOCK9_DB_SECRET} \
	--env LOCK9_DB_NAME=${LOCK9_DB_NAME} \
	--env LOCK9_JAPANGROUP_FORMURL=${LOCK9_JAPANGROUP_FORMURL} \
	--link mariadb:db \
	-d \
	mingchoi/lock9

upload_image:
	sudo docker push mingchoi/lock9
	
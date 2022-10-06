image:
	docker build -t omma-kebab-server:latest -f Dockerfile .

container:
	docker run -p 8081:8081 --env-file ./local.env --link some-mariadb:db \
	--name omma-kebab-server omma-kebab-server:latest

run_container:
	docker run -p 8081:8081 --env-file ./local.env --link omma-kebab-mariadb:db \
	--name omma-kebab-server omma-kebab-server:1.0.0
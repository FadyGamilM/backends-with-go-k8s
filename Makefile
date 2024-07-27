build-counter:
	docker build -f counter/Dockerfile -t fadygamil/counter .

build-poller:
	docker build -f poller/Dockerfile -t fadygamil/poller .

build-server:
	docker build -f server/Dockerfile -t fadygamil/server .

run-server:
	docker run --name server -p 8081:8081 fadygamil/server

run-counter:
	docker run --name counter fadygamil/counter

run-poller:
	docker run --name poller fadygamil/poller

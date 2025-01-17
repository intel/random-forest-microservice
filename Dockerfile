FROM debian

ENV DEBIAN_FRONTEND=noninteractive
SHELL ["/bin/bash", "-o", "pipefail", "-c"]
WORKDIR /app
# Update packages
RUN apt-get update && apt-get upgrade --no-install-recommends -y && apt-get install pip make wget -y --no-install-recommends
# Pull in data
COPY ./random_forest ./random_forest
# Install dependencies
RUN pip install --break-system-packages --upgrade pip && pip install --break-system-packages -r /app/random_forest/requirements.txt 
# Install Go
COPY ./Makefile ./Makefile
ENV PATH=$PATH:/usr/local/go/bin
RUN make install_deps
COPY ./oddforest-microservice ./oddforest-microservice
# Build API Server
WORKDIR /app/oddforest-microservice/src
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ../../bin/oddforest_server.run main.go
# Run API Server
CMD ["/bin/bash","-c", "/app/bin/oddforest_server.run"]
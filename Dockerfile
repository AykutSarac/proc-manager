FROM golang:1.17 as development
WORKDIR /app
COPY . /app 
RUN ["go", "install"]
ENTRYPOINT ["manager"]
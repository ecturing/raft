FROM golang:1.22.2-bullseye
EXPOSE 8080 1234
RUN make clean
RUN make build
COPY ./build/raft .
RUN chmod +x raft
CMD ["./raft"]
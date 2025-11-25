FROM gcr.io/distroless/static-debian12:nonroot
COPY service-example-go /
ENTRYPOINT ["/service-example-go"]

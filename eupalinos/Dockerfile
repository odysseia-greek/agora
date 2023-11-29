FROM alpine

ARG project_name
ARG TARGETOS
ARG TARGETARCH

ENV project_name=${project_name}

RUN mkdir /app
WORKDIR /

COPY dist/bin/${TARGETOS}/${TARGETARCH}/${project_name} /app/${project_name}

ENTRYPOINT [ "sh", "-c", "/app/${project_name}" ]
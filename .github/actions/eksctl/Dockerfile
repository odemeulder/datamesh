FROM alpine:latest

LABEL "com.github.actions.name"="eksctl"
LABEL "com.github.actions.description"="eksctl is a simple CLI tool for creating clusters on EKS"
LABEL "com.github.actions.icon"="server"
LABEL "com.github.actions.color"="blue"

COPY entrypoint.sh /

# install curl and eksctl
RUN apk add --update --no-cache --virtual .build-deps curl
RUN curl -s --location "https://github.com/weaveworks/eksctl/releases/download/0.63.0/eksctl_$(uname -s)_amd64.tar.gz" | tar xz -C /tmp 
RUN mv /tmp/eksctl /usr/local/bin
RUN apk del .build-deps

# Set AWS credentials
RUN mkdir ~/.aws
RUN echo "[default]" >> ~/.aws/credentials
RUN echo "aws_access_key_id = $AWS_ACCESS_KEY_ID" >>~/.aws/credentials
RUN echo "aws_secret_access_key = $AWS_SECRET_ACCESS_KEY" >>~/.aws/credentials

ENTRYPOINT ["/entrypoint.sh"]

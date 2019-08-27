FROM golang:1.9

# setup go
ENV GOBIN $GOPATH/bin
ENV PATH $GOBIN:/usr/local/go/bin:$PATH

COPY build $GOBIN

# set docs
RUN mkdir /docs
COPY api/index.html /docs/index.html
ENV GRAPH_DOCS_DIR="/docs/*"
# set env
ENV COMMAND "serve"
RUN destrib-graph --version
CMD destrib-graph $COMMAND

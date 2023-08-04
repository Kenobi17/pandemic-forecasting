FROM python:3.9

RUN curl https://dl.google.com/go/go1.20.linux-amd64.tar.gz | tar -C /usr/local -xzf -
ENV PATH $PATH:/usr/local/go/bin
RUN pip install numpy matplotlib statsmodels

WORKDIR /app
COPY . .
RUN go mod download

RUN go build -o /app/main .

CMD ["/app/main"]

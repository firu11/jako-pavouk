# ----- frontend builder -----
FROM node:22-alpine AS frontend

WORKDIR /app

COPY frontend/package*.json ./
RUN npm ci

COPY frontend ./
RUN npm run build


# ----- backend builder -----
FROM golang:1.26-alpine AS backend

WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/. .

RUN GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o /out/server .


# ----- final image -----
FROM alpine:3.22

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata curl

COPY --from=backend /out/server /app/server
COPY --from=backend /app/overovaci_email.html /app/overovaci_email.html
COPY --from=backend /app/pavoucekDoEmailu.png /app/pavoucekDoEmailu.png
COPY --from=frontend /app/dist /app/public

EXPOSE 8080

CMD ["./server"]

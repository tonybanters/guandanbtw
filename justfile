default:
    @just --list

dev:
    just --parallel server client

server:
    cd server && air

client:
    cd client && npm run dev

build:
    cd server && go build -o bin/guandanbtw .
    cd client && npm run build

install:
    cd client && npm install

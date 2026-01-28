default:
    @just --list

dev:
    #!/usr/bin/env bash
    trap 'kill 0' EXIT
    just server &
    just client &
    wait

server:
    cd server && PORT=8081 air

client:
    cd client && npm run dev

build:
    cd server && go build -o bin/guandanbtw .
    cd client && npm run build

install:
    cd client && npm install

deploy:
    #!/usr/bin/env bash
    cd /www/sites/guandanbtw
    sudo git fetch
    sudo git merge origin/master
    just build
    sudo systemctl restart guandanbtw

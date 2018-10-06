GOOS=linux GOARCH=amd64 go get github.com/abiosoft/caddyplug/caddyplug
alias caddyplug='GOOS=linux GOARCH=amd64 caddyplug'

PLUGINS=cors,realip,expires,cache,jwt

cp /dev/null plugins.go

printf "package caddyhttp\n\n" >> plugins.go

for plugin in $(echo $PLUGINS | tr "," " "); do
    printf "import _ \"$(caddyplug package $plugin)\"\n" >> plugins.go
done

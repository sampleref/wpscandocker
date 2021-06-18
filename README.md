# wpscandocker

Simple docker running REST api to perform wpscan 

# Build docker
`docker build -t wpscandocker:1.0 --network=host .`

# Run docker
`docker run -it -p 8080:8080 --name=wpscan wpscandocker:1.0`    
`docker stop wpscan`    
`docker rm wpscan`    

# Test URL's
`curl --header "Content-Type: application/json"  --request POST  --data '{"Url":"http://usablewp.com","Action":"check"}'  http://127.0.0.1:8080/urlcheck`

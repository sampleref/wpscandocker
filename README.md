# wpscandocker

Simple docker running REST api to perform wpscan 

# Build docker
`docker build -t wpscandocker:1.0 --network=host .`

# Run docker
`docker run -it --rm -p 8080:8080 --name=wpscan wpscandocker:1.0`    
`docker stop wpscan`    
`docker rm wpscan`

# Run docker for debugging
`docker run -it --rm -v ../wpscandocker:/appsrc -p 8080:8080 --name=wpscan wpscandocker:1.0 bash`

# Test with cURL
`curl -v --header "Content-Type: application/json"  --request GET  http://127.0.0.1:8080/updatedb`

`curl -v --header "Content-Type: application/json"  --request GET  http://127.0.0.1:8080/getallreports`

`curl -v --header "Content-Type: application/json"  --request POST  --data '{"Url":"https://eracorp.io"}'  http://127.0.0.1:8080/checkurl`

`curl -v --header "Content-Type: application/json"  --request POST  --data '{"Id":"DhScAjNt"}'  http://127.0.0.1:8080/getreportbyid`

`curl -v --header "Content-Type: application/json"  --request POST  --data '{"Id":"DhScAjNt"}'  http://127.0.0.1:8080/deletereportbyid`

_Synced test cURL:_   
`id=$(curl --header "Content-Type: application/json"  --request POST  --data '{"Url":"https://eracorp.io"}'  http://127.0.0.1:8080/checkurl | jq -r '.Id' 2>/dev/null)`   
`curl -v --header "Content-Type: application/json"  --request POST  --data '{"Id":"'"$id"'"}'  http://127.0.0.1:8080/getreportbyid`   
`curl -v --header "Content-Type: application/json"  --request POST  --data '{"Id":"'"$id"'"}'  http://127.0.0.1:8080/deletereportbyid`

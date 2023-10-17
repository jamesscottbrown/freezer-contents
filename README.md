
Starting serveR:  


# API

## Getting state

curl http://localhost:8080/state


## Removing item

curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"Container":"bowl-1"}' \
  http://localhost:8080/remove


  // or /bowl-1/remove


## Adding item

curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"Name": "Soup!", "Date": "2023-04-01", "Freezer": "cellar freezer",  "Containers": ["bowl-2", "bowl-3"] }' \
  http://localhost:8080/add | jq


curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"Name": "Soup!", "Date": "2023-04-01", "Freezer": "upstairs freezer",  "Containers": ["bowl-2", "bowl-3"] }' \
  http://localhost:8080/add | jq


## Moving

curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"Container": "1", "NewFreezer": "cellar freezer" }' \
  http://localhost:8080/move | jq

Tests:

- move one of several items; check container moved but parent item remains
- move only container of item; check parent item removed
- move container, check added to existing item in new fridge
- move container, check new item created


curl http://localhost:9090/state

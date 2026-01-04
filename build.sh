cd ui/
npm install
npm run build

cd ..
# see https://stackoverflow.com/a/72518051
export CGO_ENABLED=0
go build .

go test  -coverprofile=coverage.out
# $ go install github.com/mattn/goveralls@latest
goveralls -coverprofile=coverage.out -reponame=go-webdriver -repotoken=${GOVERALLS_GO_H3DIST_TOKEN} -service=local
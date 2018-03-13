targetImageName = myapp:5.0

all: package

package: build createImage

build:
	go build -o myapp main.go
createImage:
	docker build -t "$(targetImageName)" .

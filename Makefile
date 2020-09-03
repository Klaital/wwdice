.PHONY: test clean all wwdicebot-push

# Update this env var for cross-compiling
GOOS := linux

# Buildbot will override this with the commit hash
VERSION := 1.0.0

clean:
	rm -f wwdicebot wwdicebot.exe wwdicecli wwdicecli.exe characterroller characterroller.exe

test:
	go test pkg/...

wwdicebot: pkg/dice/*.go cmd/wwdicebot/*.go
	go build cmd/wwdicebot

wwdicecli: pkg/dice/*.go cmd/wwdicecli/*.go
	go build cmd/wwdicecli

characterroller: pkg/dice/*.go pkg/characters/*.go cmd/characterroller/*.go
	go build cmd/characterroller

all: wwdicebot wwdicecli characterroller

wwdicebot-push: wwdicebot wwdicebot/Dockerfile
    docker build -t klaital/wwdicebot:$(VERSION) -f cmd/wwdicebot/Dockerfile .
    docker push klaital/wwdicebot:$(VERSION)

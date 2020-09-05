.PHONY: test clean all wwdicebot-push

# Update this env var for cross-compiling
GOOS := linux

# Buildbot will override this with the commit hash
VERSION := 1.0.0

all: wwdicebot wwdicecli characterroller

clean:
	rm -f wwdicebot wwdicebot.exe wwdicecli wwdicecli.exe characterroller characterroller.exe

test:
	go test ./pkg/...

wwdicebot: pkg/dice/*.go cmd/wwdicebot/*.go
	go build -o wwdicebot cmd/wwdicebot/main.go

wwdicecli: pkg/dice/*.go cmd/wwdicecli/*.go
	go build -o wwdicecli cmd/wwdicecli/main.go

characterroller: pkg/dice/*.go pkg/characters/*.go cmd/characterroller/*.go
	go build -o characterroller cmd/characterroller/main.go


wwdicebot-push: wwdicebot cmd/wwdicebot/Dockerfile
	docker build -t klaital/wwdicebot:$(VERSION) -f cmd/wwdicebot/Dockerfile .
	docker push klaital/wwdicebot:$(VERSION)

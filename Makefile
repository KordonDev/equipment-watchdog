PACKAGE_NAME          := github.com/kordondev/equipment-watchdog
GOLANG_CROSS_VERSION  ?= v1.21


.PHONY: release-backend-dry-run
release-backend-dry-run:
	goreleaser release --snapshot --clean

.PHONY: release-backend
release-backend:
	goreleaser release --snapshot --clean

gorelease-dry:
	goreleaser --snapshot --skip-publish --rm-dist
gorelease:
	goreleaser

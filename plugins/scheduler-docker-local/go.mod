module github.com/dokku/dokku/plugins/scheduler-docker-local

go 1.15

require (
	github.com/codeskyblue/go-sh v0.0.0-20190412065543-76bd3d59ff27
	github.com/dokku/dokku/plugins/common v0.0.0-00010101000000-000000000000
	github.com/dokku/dokku/plugins/config v0.0.0-00010101000000-000000000000
	github.com/onsi/gomega v1.10.4 // indirect
	github.com/spf13/pflag v1.0.5
)

replace github.com/dokku/dokku/plugins/common => ../common
replace github.com/dokku/dokku/plugins/config => ../config

module sekolah-platform/platform/message-broker

go 1.25.12

require (
	github.com/nats-io/nats-server/v2 v2.10.22
	github.com/nats-io/nats.go v1.37.0
	sekolah-platform/platform/context v0.0.0
	sekolah-platform/platform/events v0.0.0
)

require (
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/minio/highwayhash v1.0.3 // indirect
	github.com/nats-io/jwt/v2 v2.5.8 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.54.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/time v0.15.0 // indirect
	sekolah-platform/platform/shared v0.0.0 // indirect
)

replace sekolah-platform/platform/context => ../context

replace sekolah-platform/platform/events => ../events

replace sekolah-platform/platform/shared => ../shared

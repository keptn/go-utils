module github.com/keptn/go-utils

go 1.13

require (
	github.com/Azure/go-autorest/autorest v0.9.2
	github.com/Masterminds/goutils v1.1.0 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/Masterminds/sprig v2.22.0+incompatible // indirect
	github.com/cloudevents/sdk-go v0.9.2
	github.com/cyphar/filepath-securejoin v0.2.2 // indirect
	github.com/go-openapi/errors v0.19.2
	github.com/go-openapi/strfmt v0.19.3
	github.com/go-openapi/swag v0.19.5
	github.com/go-openapi/validate v0.19.3
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/google/uuid v1.1.1
	github.com/gophercloud/gophercloud v0.5.0 // indirect
	github.com/gorilla/websocket v1.4.1
	github.com/huandu/xstrings v1.2.0 // indirect
	github.com/imdario/mergo v0.3.8 // indirect
	github.com/mitchellh/copystructure v1.0.0 // indirect
	gopkg.in/src-d/go-git.v4 v4.13.1
	gopkg.in/yaml.v2 v2.2.4
	k8s.io/api v0.0.0-20191003000013-35e20aa79eb8
	k8s.io/apimachinery v0.0.0-20190913080033-27d36303b655
	k8s.io/client-go v0.0.0-20191003000419-f68efa97b39e
	k8s.io/helm v2.14.3+incompatible
	k8s.io/utils v0.0.0-20191010214722-8d271d903fe4 // indirect
)

replace (
	k8s.io/api => k8s.io/api v0.0.0-20191003000013-35e20aa79eb8
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190913080033-27d36303b655
	k8s.io/client-go => k8s.io/client-go v0.0.0-20191003000419-f68efa97b39e
)

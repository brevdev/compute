generate-shadeform-client:
	docker run --rm -v "${PWD}:/local" openapitools/openapi-generator-cli:v7.8.0 generate \
		--additional-properties disallowAdditionalPropertiesIfNotPresent=false \
		-i https://raw.githubusercontent.com/shadeform/docs/refs/heads/main/openapi.yaml \
		-g go \
		--git-user-id shadeform \
		--git-repo-id cloud-shadeform \
		-o /local/gen/shadeform
	sudo chown -R $(shell id -u):$(shell id -g) gen/shadeform
	gofmt -s -w gen/shadeform
	rm -rf gen/shadeform/go.mod gen/shadeform/go.sum


validate-shadeform-openapi:
	docker run --rm -v ${PWD}:/local -v "/Users/ronald/shade/docs:/docs" openapitools/openapi-generator-cli:v7.8.0 validate -i /docs/openapi.yaml

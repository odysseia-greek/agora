

generate_docs:
	docker run --rm \
		-v $(CURDIR)/docs:/out \
		-v $(CURDIR)/proto:/protos \
		localproto:latest --doc_opt=markdown,docs.md
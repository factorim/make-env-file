SERVICE := Make Env File
  
.PHONY: help build

help:
	@printf "\
	$(SERVICE)\n\
	\n\
	$(bold)SYNOPSIS$(sgr0)\n\
    make [COMMANDS]\n\
	\n\
	$(bold)COMMANDS$(sgr0)\n\
		$(bold)help$(sgr0)\n\
		  Display help\n\n\
		$(bold)build$(sgr0)\n\
		  Create a make-env-file executable\n\n\
	"
	
build:
	docker build -t factorim/make-env-file .
	docker create -ti --name temp-image factorim/make-env-file bash
	docker cp temp-image:/go/bin/make-env-file .
	docker rm -f temp-image
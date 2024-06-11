PROTOC=protoc

# export PATH := $(TESTDIR):$(PATH)
# Macro to execute a command passed as argument.
# DO NOT DELETE EMPTY LINE at the end of the macro, it's required to separate commands.
define exec-command
$(1)

endef

# DO NOT DELETE EMPTY LINE at the end of the macro, it's required to separate commands.
define print_caption
  @echo "🏗️ "
  @echo "🏗️ " $1
  @echo "🏗️ "

endef

# Macro to compile Protobuf $(2) into directory $(1). $(3) can provide additional flags.
# DO NOT DELETE EMPTY LINE at the end of the macro, it's required to separate commands.
# Arguments:
#  $(1) - output directory
#  $(2) - path to the .proto file
define proto_compile
  $(call print_caption, "Processing $(2) --> $(1)")

  $(PROTOC) \
    --go_out=$(strip $(1)) \
	  --go_opt=paths=source_relative \
	  --go-grpc_opt=paths=source_relative \
    --go-grpc_out=$(strip $(1)) \
    $(2)

endef

.PHONY: proto
proto: proto-server 

.PHONY: proto-model
proto-server: install-protogen-deps
	$(call proto_compile, proto-gen, models/server.proto)

.PHONY: install-protogen-deps
install-protogen-deps:
	$(GO) install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	$(GO) install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

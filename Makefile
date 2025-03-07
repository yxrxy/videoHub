MODULE := github.com/yxrrxy/videoHub
IDL_PATH := idl

.PHONY: kitex-gen-%
kitex-gen-%:
	@ kitex -module "${MODULE}" \
		${IDL_PATH}/$*.thrift
	@ go mod tidy
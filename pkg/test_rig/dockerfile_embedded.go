package test_rig

func GetDokerfileStub() string {
	return `FROM scratch
COPY ./main /main
CMD ["/main"]
`
}

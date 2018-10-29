package parser

// CheckError panics at the presence of an error.
func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

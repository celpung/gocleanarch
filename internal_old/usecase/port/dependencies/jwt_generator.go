package dependencies

type JwtGenerator interface {
	Generate(userID, email, role string) (string, error)
}

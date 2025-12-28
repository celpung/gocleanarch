package dependencies

type IDGenerator interface {
	NewID() (string, error)
}

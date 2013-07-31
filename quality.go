package quality

// InformationNeed is the Information Need the user wishes to satisfy.
type InformationNeed struct {
	Kind  string
	Query string
}

// Doc is a query result. Kind can be used by a universal search returning
// multiple types of objects.
type Doc struct {
	Kind string
	ID   int
}

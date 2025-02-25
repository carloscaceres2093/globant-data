package health

// service contains everything you need to health.
type service struct{}

// NewService creates a Health.
func NewService() *service {
	return &service{}
}

// Check checks that it is ok.
func (s *service) Check() Response {
	return Response{
		"ok",
	}
}

package health

type service struct{}

func NewService() *service {
	return &service{}
}

func (s *service) Check() Response {
	return Response{
		"ok",
	}
}

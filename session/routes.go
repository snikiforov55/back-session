package session

func (s *Service) routes() {
	s.Router.Methods("POST").Name("CreateSession").
		Path("/session").HandlerFunc(s.handleCreateSession())
	s.Router.Methods("PUT").Name("UpdateAuthInfo").
		Path("/session").HandlerFunc(s.handleUpdateAuthInfo())
	s.Router.Methods("DELETE").Name("DropSession").
		Path("/session/{id}").HandlerFunc(s.handleDropSession())
	s.Router.Methods("GET").Name("GetSessionAttributes").
		Path("/session/{id}").HandlerFunc(s.handleGetSessionAttributes())
}

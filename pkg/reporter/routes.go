package reporter

func (s *Server) routes() {
	s.router.HandleFunc("/", s.handleDashboard())
	s.router.HandleFunc("/favicon.ico", s.handleFavicon())
}

package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"sekolah-platform/services/tu-core/internal/service"
)

type Server struct {
	router *gin.Engine
	port   int

	studentService        *service.StudentService
	mutationService       *service.MutationService
	registrationService   *service.RegistrationService
	attendanceService     *service.AttendanceService
	assessmentService     *service.AssessmentService
	financeService        *service.FinanceService
	correspondenceService *service.CorrespondenceService
}

func NewServer(
	port int,
	studentService *service.StudentService,
	mutationService *service.MutationService,
	registrationService *service.RegistrationService,
	attendanceService *service.AttendanceService,
	assessmentService *service.AssessmentService,
	financeService *service.FinanceService,
	correspondenceService *service.CorrespondenceService,
) *Server {
	s := &Server{
		router: gin.Default(),
		port:   port,

		studentService:        studentService,
		mutationService:       mutationService,
		registrationService:   registrationService,
		attendanceService:     attendanceService,
		assessmentService:     assessmentService,
		financeService:        financeService,
		correspondenceService: correspondenceService,
	}

	s.setupMiddleware()
	s.setupRoutes()

	return s
}

func (s *Server) Start() error {
	return s.router.Run(fmt.Sprintf(":%d", s.port))
}

func (s *Server) setupMiddleware() {
	// 1. CORS Middleware
	s.router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, x-school-id, x-academic-period-id, x-user-id, x-user-role")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	// 2. Operational Context Extractor
	s.router.Use(func(c *gin.Context) {
		schoolID := c.GetHeader("x-school-id")
		academicPeriodID := c.GetHeader("x-academic-period-id")
		userID := c.GetHeader("x-user-id")
		role := c.GetHeader("x-user-role")

		if schoolID == "" || academicPeriodID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing operational context headers (x-school-id, x-academic-period-id)"})
			c.Abort()
			return
		}

		// Simpan ke gin.Context agar bisa diambil handler
		c.Set("school_id", schoolID)
		c.Set("academic_period_id", academicPeriodID)
		c.Set("user_id", userID)
		c.Set("role", role)
		c.Next()
	})

	// 3. Global Error Handler
	s.router.Use(func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			})
		}
	})
}

func (s *Server) setupRoutes() {
	s.router.GET("/health", s.healthCheck)

	v1 := s.router.Group("/api/v1")
	{
		students := v1.Group("/students")
		{
			students.POST("", s.createStudent)
			students.GET("", s.listStudents)
			students.GET("/:id", s.getStudent)
			students.PUT("/:id", s.updateStudent)
			students.DELETE("/:id", s.deleteStudent)
			students.POST("/:id/transfer", s.transferStudent)
			students.POST("/:id/graduate", s.graduateStudent)
		}

		mutations := v1.Group("/mutations")
		{
			mutations.POST("", s.createMutation)
			mutations.GET("", s.listMutations)
			mutations.GET("/:id", s.getMutation)
			mutations.POST("/:id/approve", s.approveMutation)
			mutations.POST("/:id/reject", s.rejectMutation)
		}

		registrations := v1.Group("/registrations")
		{
			registrations.POST("", s.createRegistration)
			registrations.GET("", s.listRegistrations)
			registrations.GET("/:id", s.getRegistration)
			registrations.POST("/:id/submit", s.submitRegistration)
			registrations.POST("/:id/verify", s.verifyRegistration)
			registrations.POST("/:id/approve", s.approveRegistration)
			registrations.POST("/:id/reject", s.rejectRegistration)
		}

		attendance := v1.Group("/attendance")
		{
			attendance.POST("/sessions", s.createSession)
			attendance.GET("/sessions", s.listSessions)
			attendance.GET("/sessions/:id", s.getSession)
			attendance.POST("/sessions/:id/submit", s.submitAttendance)
			attendance.GET("/students/:studentId", s.getStudentAttendance)
			attendance.GET("/summary", s.getAttendanceSummary)
		}

		assessments := v1.Group("/assessments")
		{
			assessments.POST("", s.createAssessment)
			assessments.GET("", s.listAssessments)
			assessments.GET("/:id", s.getAssessment)
			assessments.PUT("/:id", s.updateAssessment)
			assessments.GET("/students/:studentId", s.getStudentAssessments)
			assessments.GET("/students/:studentId/average", s.getStudentAverage)
		}

		finance := v1.Group("/finance")
		{
			invoices := finance.Group("/invoices")
			{
				invoices.POST("", s.createInvoice)
				invoices.GET("", s.listInvoices)
				invoices.GET("/:id", s.getInvoice)
				invoices.GET("/students/:studentId/outstanding", s.getStudentOutstanding)
			}
			payments := finance.Group("/payments")
			{
				payments.POST("", s.createPayment)
				payments.GET("", s.listPayments)
			}
			finance.GET("/summary", s.getPaymentSummary)
		}

		correspondence := v1.Group("/correspondence")
		{
			correspondence.POST("", s.createCorrespondence)
			correspondence.GET("", s.listCorrespondences)
			correspondence.GET("/:id", s.getCorrespondence)
			correspondence.PUT("/:id", s.updateCorrespondence)
			correspondence.POST("/:id/archive", s.archiveCorrespondence)
		}
	}
}

// --- HANDLERS (IMPLEMENTED) ---

func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "tu-core-rest"})
}

func (s *Server) createStudent(c *gin.Context) {
	schoolID, _ := c.Get("school_id")
	academicPeriodID, _ := c.Get("academic_period_id")
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var req service.CreateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ✅ SIGNATURE SESUAI ERROR COMPILER: (context, string, string, int, string, string, req)
	student, err := s.studentService.CreateStudent(
		c.Request.Context(),
		schoolID.(string),
		academicPeriodID.(string),
		1, // default semester
		userID.(string),
		role.(string),
		req,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, student)
}

func (s *Server) listStudents(c *gin.Context) {
	schoolID, _ := c.Get("school_id")
	academicPeriodID, _ := c.Get("academic_period_id")

	var filter service.StudentFilterRequest
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ✅ SIGNATURE SESUAI ERROR COMPILER: (context, string, string, filter)
	students, err := s.studentService.ListStudents(
		c.Request.Context(),
		schoolID.(string),
		academicPeriodID.(string),
		filter,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, students)
}

// --- STUBS FOR ALL OTHER ROUTES (TO PREVENT COMPILE ERRORS) ---
func (s *Server) getStudent(c *gin.Context)               { c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"}) }
func (s *Server) updateStudent(c *gin.Context)            { c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"}) }
func (s *Server) deleteStudent(c *gin.Context)            { c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"}) }
func (s *Server) transferStudent(c *gin.Context)          { c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"}) }
func (s *Server) graduateStudent(c *gin.Context)          { c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"}) }

func (s *Server) createMutation(c *gin.Context)           { c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"}) }
func (s *Server) listMutations(c *gin.Context)            { c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"}) }
func (s *Server) getMutation(c *gin.Context)              { c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"}) }
func (s *Server) approveMutation(c *gin.Context)          { c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"}) }
func (s *Server) rejectMutation(c *gin.Context)           { c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"}) }

func (s *Server) createRegistration(c *gin.Context)       { c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"}) }
func (s *Server) listRegistrations(c *gin.Context)        { c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"}) }
func (s *Server) getRegistration(c *gin.Context)          { c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"}) }
func (s *Server) submitRegistration(c *gin.Context)       { c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"}) }
func (s *Server) verifyRegistration(c *gin.Context)       { c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"}) }
func (s *Server) approveRegistration(c *gin.Context)      { c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"}) }
func (s *Server) rejectRegistration(c *gin.Context)       { c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"}) }

func (s *Server) createSession(c *gin.Context)            { c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"}) }
func (s *Server) listSessions(c *gin.Context)             { c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"}) }
func (s *Server) getSession(c *gin.Context)               { c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"}) }
func (s *Server) submitAttendance(c *gin.Context)         { c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"}) }
func (s *Server) getStudentAttendance(c *gin.Context)     { c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"}) }
func (s *Server) getAttendanceSummary(c *gin.Context)     { c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"}) }

func (s *Server) createAssessment(c *gin.Context)         { c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"}) }
func (s *Server) listAssessments(c *gin.Context)          { c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"}) }
func (s *Server) getAssessment(c *gin.Context)            { c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"}) }
func (s *Server) updateAssessment(c *gin.Context)         { c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"}) }
func (s *Server) getStudentAssessments(c *gin.Context)    { c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"}) }
func (s *Server) getStudentAverage(c *gin.Context)        { c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"}) }

func (s *Server) createInvoice(c *gin.Context)            { c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"}) }
func (s *Server) listInvoices(c *gin.Context)             { c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"}) }
func (s *Server) getInvoice(c *gin.Context)               { c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"}) }
func (s *Server) getStudentOutstanding(c *gin.Context)    { c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"}) }
func (s *Server) createPayment(c *gin.Context)            { c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"}) }
func (s *Server) listPayments(c *gin.Context)             { c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"}) }
func (s *Server) getPaymentSummary(c *gin.Context)        { c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"}) }

func (s *Server) createCorrespondence(c *gin.Context)     { c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"}) }
func (s *Server) listCorrespondences(c *gin.Context)      { c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"}) }
func (s *Server) getCorrespondence(c *gin.Context)        { c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"}) }
func (s *Server) updateCorrespondence(c *gin.Context)     { c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"}) }
func (s *Server) archiveCorrespondence(c *gin.Context)    { c.JSON(http.StatusOK, gin.H{"message": "Not implemented yet"}) }

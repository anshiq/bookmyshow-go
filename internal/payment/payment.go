package payment

// PaymentProvider defines the interface for payment providers
type PaymentProvider interface {
	Name() string
	ProcessPayment(amount float64, currency string, metadata map[string]string) (*PaymentResult, error)
	VerifyPayment(transactionID string) (*PaymentResult, error)
	RefundPayment(transactionID string, amount float64) (*PaymentResult, error)
}

// PaymentResult represents the result of a payment operation
type PaymentResult struct {
	TransactionID string
	Status        string // success, failed, pending
	Amount        float64
	Currency      string
	Error         string
}

// PaymentService manages payment operations
type PaymentService struct {
	providers map[string]PaymentProvider
}

// NewPaymentService creates a new payment service
func NewPaymentService() *PaymentService {
	return &PaymentService{
		providers: make(map[string]PaymentProvider),
	}
}

// RegisterProvider registers a payment provider
func (s *PaymentService) RegisterProvider(name string, provider PaymentProvider) {
	s.providers[name] = provider
}

// GetProvider returns a payment provider by name
func (s *PaymentService) GetProvider(name string) PaymentProvider {
	return s.providers[name]
}

// ProcessPayment processes payment using the specified provider
func (s *PaymentService) ProcessPayment(providerName string, amount float64, currency string, metadata map[string]string) (*PaymentResult, error) {
	provider, ok := s.providers[providerName]
	if !ok {
		return nil, &PaymentError{Message: "Provider not found: " + providerName}
	}
	return provider.ProcessPayment(amount, currency, metadata)
}

// PaymentError represents a payment error
type PaymentError struct {
	Message string
}

func (e *PaymentError) Error() string {
	return e.Message
}
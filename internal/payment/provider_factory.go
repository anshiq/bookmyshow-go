package payment

// PaymentProviderFactory creates payment providers
type PaymentProviderFactory struct {
	service *PaymentService
}

// NewPaymentProviderFactory creates a new factory
func NewPaymentProviderFactory() *PaymentProviderFactory {
	f := &PaymentProviderFactory{
		service: NewPaymentService(),
	}

	// Register default providers (these would be initialized with actual credentials in production)
	// For now, we provide stub providers that simulate behavior
	f.service.RegisterProvider("paytm", &PaytmProvider{merchantID: "demo", merchantKey: "demo"})
	f.service.RegisterProvider("phonepe", &PhonePeProvider{merchantID: "demo", saltKey: "demo", saltIndex: 1})

	return f
}

// GetPaymentService returns the payment service
func (f *PaymentProviderFactory) GetPaymentService() *PaymentService {
	return f.service
}

// GetProvider returns a provider by name
func (f *PaymentProviderFactory) GetProvider(name string) PaymentProvider {
	return f.service.GetProvider(name)
}
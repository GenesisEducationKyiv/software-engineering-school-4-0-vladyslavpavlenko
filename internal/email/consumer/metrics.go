package consumer

// calculateEmailSuccessRate computes the success rate of email sending operations.
func calculateEmailSuccessRate() float64 {
	totalAttempts := sentEmailsCounter.Get() + notSentEmailsCounter.Get()
	if totalAttempts == 0 {
		return 0
	}
	return float64(sentEmailsCounter.Get()) / float64(totalAttempts)
}

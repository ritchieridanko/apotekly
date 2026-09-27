package services

const defaultConfig string = `{
	"methodConfig": [{
		"retryPolicy": {
			"maxAttempts": 3,
			"initialBackoff": "0.1s",
			"maxBackoff": "1s",
			"backoffMultiplier": 2.0,
			"retryableStatusCodes": ["UNAVAILABLE", "RESOURCE_EXHAUSTED"]
		}
	}]
}`

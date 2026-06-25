package models

type (
	AccessToken struct {
		Token            string
		ExpiresInSeconds uint64
	}

	RefreshToken struct {
		Token            string
		ExpiresInSeconds uint64
	}

	AuthToken struct {
		AccessToken  *AccessToken
		RefreshToken *RefreshToken
	}
)

package tests

import (
	"errors"
	"testing"
	"time"

	"github.com/prorok210/AvitoShop/config"
	"github.com/prorok210/AvitoShop/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestJWTFunctions(t *testing.T) {
	testCases := []struct {
		name        string
		username    string
		expectedErr error
		testFunc    func() error
	}{
		{
			name:        "Generate Secret Key",
			username:    "",
			expectedErr: nil,
			testFunc: func() error {
				accKey, err := utils.GenerateSecretKey(32)
				t.Setenv("JWT_ACCESS_SECRET", accKey)
				return err
			},
		},
		{
			name:        "Generate Access Token",
			username:    "testuser",
			expectedErr: nil,
			testFunc: func() error {
				_, err := utils.GenerateAccessToken("testuser")
				return err
			},
		},
		{
			name:        "Validate Access Token",
			username:    "testuser",
			expectedErr: nil,
			testFunc: func() error {
				token, err := utils.GenerateAccessToken("testuser")
				if err != nil {
					return err
				}
				_, err = utils.ValidateToken(token)
				return err
			},
		},
		{
			name:        "Invalid Access Token",
			username:    "",
			expectedErr: errors.New("error expected"),
			testFunc: func() error {
				_, err := utils.ValidateToken("invalidToken")
				return err
			},
		},
		{
			name:        "Expired Access Token",
			username:    "testuser",
			expectedErr: errors.New("error expected"),
			testFunc: func() error {
				config.JWT_ACCESS_EXPIRATION_TIME = time.Millisecond * 100
				token, err := utils.GenerateAccessToken("testuser")
				if err != nil {
					return err
				}
				time.Sleep(time.Millisecond * 200)
				_, err = utils.ValidateToken(token)
				config.JWT_ACCESS_EXPIRATION_TIME = time.Minute * 5
				return err
			},
		},
		{
			name:        "Generate Access Token with Empty Username",
			username:    "",
			expectedErr: errors.New("username is empty"),
			testFunc: func() error {
				_, err := utils.GenerateAccessToken("")
				return err
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.testFunc()
			if tc.expectedErr != nil {
				assert.Error(t, err, "Expected error but got nil")
			} else {
				assert.NoError(t, err, "Unexpected error: %v", err)
			}
		})
	}
}

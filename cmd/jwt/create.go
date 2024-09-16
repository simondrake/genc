package jwt

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/cobra"
)

func newCreateCommand() *cobra.Command {
	var (
		signingKey string
		claims     string
	)

	// TODO: Add Signing Method flag

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "create jwt",
		Example: `
    # Create a jwt with :
      # The following reserved claims:
        # => nbf (not before time): Time before which the JWT must not be accepted for processing
        # => iat (issued at time): Time at which the JWT was issued; can be used to determine age of the JWT
        # => exp (expiration time): Time after which the JWT expires
        # => sub (subject): Subject of the JWT (the user)
      # The following custom claims:
        # => groups (A list of permission groups for the user, delimited by a semi-colon)
    $ now=$(date +%s) exp=$(date -d "-1 month" +%s); genc jwt create --signing-key "verysecret" --claims "{\"nbf\": $now, \"iat\": $now, \"exp\": $exp, \"sub\": \"imsudonow\", \"groups\": [\"admin\", \"superadmin\"]}"

    # Create a jwt with no claims
    $ genc jwt create --signing-key "verysecret"`,
		Run: func(cmd *cobra.Command, args []string) {
			var m map[string]interface{}

			if claims != "" {
				if err := json.Unmarshal([]byte(claims), &m); err != nil {
					fmt.Fprintln(os.Stderr, fmt.Errorf("error parsing claims: %w", err))
					os.Exit(1)
				}
			}

			tkn := createToken(m)

			sig, err := tkn.SignedString([]byte(signingKey))
			if err != nil {
				fmt.Fprintln(os.Stderr, fmt.Errorf("error getting signed token: %w", err))
				os.Exit(1)
			}

			fmt.Fprintln(os.Stdout, sig)
		},
	}

	createCmd.Flags().StringVar(&signingKey, "signing-key", "", "the signing key to create the jwt with")
	createCmd.Flags().StringVar(&claims, "claims", "", "claims that the jwt should be created with")

	if err := createCmd.MarkFlagRequired("signing-key"); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("internal error marking flag 'signing-key' as required: %w", err))
	}

	return createCmd
}

func createToken(claims map[string]interface{}) *jwt.Token {
	if len(claims) == 0 {
		return jwt.New(jwt.SigningMethodHS256)
	}

	mc := jwt.MapClaims{}

	for k, v := range claims {
		mc[k] = v
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, mc)
}

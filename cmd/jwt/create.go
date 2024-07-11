package jwt

import (
	"fmt"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/cobra"
)

func newCreateCommand() *cobra.Command {
	var (
		signingKey    string
		claims        map[string]string
		listDelimiter string
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
    $ now=$(date +%s) exp=$(date -d "+1 year" +%s); genc jwt create --signing-key "verysecret" --claims "nbf=$now,iat=$now,exp=$exp,sub=imsudonow,groups=admin;superadmin"

    # Create a jwt with no claims
    $ genc jwt create --signing-key "verysecret"`,
		Run: func(cmd *cobra.Command, args []string) {
			tkn := createToken(claims, listDelimiter)

			sig, err := tkn.SignedString([]byte(signingKey))
			if err != nil {
				fmt.Fprintln(os.Stderr, fmt.Errorf("error getting signed token: %w", err))
				os.Exit(1)
			}

			fmt.Fprintln(os.Stdout, sig)
		},
	}

	createCmd.Flags().StringVar(&signingKey, "signing-key", "", "the signing key to create the jwt with")
	createCmd.Flags().StringVar(&listDelimiter, "list-delimiter", ";", "the delimiter to use when specifying a claim that is a list")
	createCmd.Flags().StringToStringVar(&claims, "claims", nil, "a map of claims that the jwt should be created with")

	if err := createCmd.MarkFlagRequired("signing-key"); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("internal error marking flag 'string' as required: %w", err))
	}

	return createCmd
}

func createToken(claims map[string]string, delim string) *jwt.Token {
	if len(claims) == 0 {
		return jwt.New(jwt.SigningMethodHS256)
	}

	mc := jwt.MapClaims{}

	for k, v := range claims {
		if !strings.Contains(v, delim) {
			mc[k] = v
			continue
		}

		mc[k] = strings.Split(v, delim)
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, mc)
}

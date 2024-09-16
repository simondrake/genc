package jwt

import (
	"errors"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/cobra"
)

func newParseCommand() *cobra.Command {
	var (
		token                  string
		signingKey             string
		allowInvalidSigningKey bool
	)

	parseCmd := &cobra.Command{
		Use:   "parse",
		Short: "parse jwt",
		Example: `
    `,
		Run: func(cmd *cobra.Command, args []string) {
			claims, err := parseToken(token, signingKey, allowInvalidSigningKey)
			if err != nil {
				fmt.Fprintln(os.Stderr, fmt.Errorf("error parsing token: %w", err))
				os.Exit(1)
			}

			fmt.Fprintln(os.Stdout, *claims)
		},
	}

	parseCmd.Flags().StringVar(&token, "token", "", "the jwt token to parse")
	parseCmd.Flags().StringVar(&signingKey, "signing-key", "", "the signing key to create the jwt with")
	parseCmd.Flags().BoolVar(&allowInvalidSigningKey, "allow-invalid-signing-key", false, "whether to allow an invalid signing key")

	if err := parseCmd.MarkFlagRequired("token"); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("internal error marking flag 'token' as required: %w", err))
	}

	return parseCmd
}

func parseToken(token string, signingKey string, allowInvalidSigningKey bool) (*jwt.MapClaims, error) {
	tkn, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(signingKey), nil
	})

	switch {
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		if !allowInvalidSigningKey {
			return nil, fmt.Errorf("invalid token signature: %w", err)
		}
		// If an invalid signing key is allowed, treat it as valid
		return parseClaims(tkn)
	case tkn.Valid:
		return parseClaims(tkn)
	case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
		fmt.Fprintln(os.Stdout, "Warning: Token is expired or not yet valid")
		return parseClaims(tkn)
	case errors.Is(err, jwt.ErrTokenMalformed):
		return nil, fmt.Errorf("malformed token: %w", err)
	}

	return nil, fmt.Errorf("unexpected error: %w", err)
}

func parseClaims(tkn *jwt.Token) (*jwt.MapClaims, error) {
	claims, ok := tkn.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("token claim type is unexpected")
	}

	return &claims, nil
}

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
	if err := handleParseError(err, allowInvalidSigningKey); err != nil {
		return nil, err
	}

	claims, ok := tkn.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("token claim type is unexpected")
	}

	return &claims, nil
}

func handleParseError(err error, allowInvalidSigningKey bool) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, jwt.ErrTokenExpired) {
		fmt.Fprintln(os.Stdout, "Warning: Token is expired")
	}

	if errors.Is(err, jwt.ErrTokenNotValidYet) {
		fmt.Fprintln(os.Stdout, "Warning: Token is not yet valid")
	}

	if errors.Is(err, jwt.ErrTokenSignatureInvalid) && !allowInvalidSigningKey {
		return fmt.Errorf("invalid token signature: %w", err)
	}

	if errors.Is(err, jwt.ErrTokenMalformed) {
		return fmt.Errorf("malformed token: %w", err)
	}

	if !errors.Is(err, jwt.ErrTokenExpired) && !errors.Is(err, jwt.ErrTokenNotValidYet) && !errors.Is(err, jwt.ErrTokenSignatureInvalid) {
		return fmt.Errorf("unexpected error: %w", err)
	}

	return nil
}

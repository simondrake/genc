# Genc

Genc is a CLI tool to expose useful commands for encryption/decryption, and other security-focused operations.

It is in very early development, so it comes with no assurances. It has only been tested on Ubuntu Linux. It _may_ work with your system, but it is in no way guaranteed.

**Note:** Expect breaking changes to occur until `v0.1.0`.

# Setup/Installation

* Install - `go install github.com/simondrake/genc@latest`

# Notes

* The `pkcs7` commands currently only support PKCS1 Private Keys

# Usage

TBC

# Examples

## PKCS7

### PKCS1

```bash
# Generate public/private key pair
$ openssl genrsa -out rsa.key 2048

# @@ Optional @@
# Decode the private key
$ openssl rsa -text -in rsa.key -noout

# Extract the Public Key and save
$ openssl rsa -in rsa.key -pubout -out rsa.pub

# Create a CSR so we can sign the Certificate
$ openssl req -key rsa.key -new -out domain.csr

# @@ Optional @@
# Verify the CSR
$ openssl req -text -in domain.csr -noout -verify

# !!! TODO - Fill out the information from this command with whatever details you want !!!

# Create a Self-Signed Certificate
## Note: The -days option specifies the number of days that the certificate will be valid.
$ openssl x509 -signkey rsa.key -in domain.csr -req -days 365 -out domain.crt

# @@ Optional @@
# Verify the Certificate
$ openssl x509 -text -in domain.crt -noout

# @@ Optional @@
# Verify the Public and Private keys match (All files should share the same public key and the same hash value)
$ openssl pkey -pubout -in rsa.key | openssl sha256
$ openssl req -pubkey -in domain.csr -noout | openssl sha256
$ openssl x509 -pubkey -in domain.crt -noout | openssl sha256

# Convert Private Key to PKCS1
$ openssl pkey -in rsa.key -traditional -out pkcs1.key

# Encrypt
$ go run . pkcs7 encrypt --public-key domain.crt --string "test"
MIIBvQYJKoZIhvcNAQcDoIIBrjCCAaoCAQAxggF3MIIBcwIBADBdMEUxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQCFF/tVGl7kJuv9ogqYX57t5G9+LM2MAsGCSqGSIb3DQEBAQSCAQA5tdqDFBFJc0uc6p8Co8nqwYhnIX6s2AG+yX0Gi0FYD0KEF9UJiVMfXtkjNUh5uvpW3dA+LThb2la6d0eOAy7KbGxjh/Pujs8q3XdIMRHcUTdkTIPw8JqkvjrAHC6Sj78fT+5okWdO2Yj6p52YPtMH1soAArx4X1T1aXkwhYSWfbQ6ZROIrX7hSsCfV/q276ERw26U4wV5i6EzZt1E5yodfJXVOsbWXekfqLl5ZjjcdGb4T6muutyRTDBaFuB78XsUZfiI7cqOC1IieWad/e7/Uje9loqo+nmZ1pTqoCYfzTcHlAYoTNq80ewsBlv9RiqtIpYq2n2Q7X4wp8sg6+zHMCoGCSqGSIb3DQEHATARBgUrDgMCBwQI7EhY6z+uB4egCgQIkSNbGva47+o=

# Decrypt
$ go run . pkcs7 decrypt --private-key pkcs1.key --public-key domain.crt --string "MIIBvQYJKoZIhvcNAQcDoIIBrjCCAaoCAQAxggF3MIIBcwIBADBdMEUxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQCFF/tVGl7kJuv9ogqYX57t5G9+LM2MAsGCSqGSIb3DQEBAQSCAQA5tdqDFBFJc0uc6p8Co8nqwYhnIX6s2AG+yX0Gi0FYD0KEF9UJiVMfXtkjNUh5uvpW3dA+LThb2la6d0eOAy7KbGxjh/Pujs8q3XdIMRHcUTdkTIPw8JqkvjrAHC6Sj78fT+5okWdO2Yj6p52YPtMH1soAArx4X1T1aXkwhYSWfbQ6ZROIrX7hSsCfV/q276ERw26U4wV5i6EzZt1E5yodfJXVOsbWXekfqLl5ZjjcdGb4T6muutyRTDBaFuB78XsUZfiI7cqOC1IieWad/e7/Uje9loqo+nmZ1pTqoCYfzTcHlAYoTNq80ewsBlv9RiqtIpYq2n2Q7X4wp8sg6+zHMCoGCSqGSIb3DQEHATARBgUrDgMCBwQI7EhY6z+uB4egCgQIkSNbGva47+o="
test
```

### PKCS8

```bash
# Generate public/private key pair
$ openssl genrsa -out rsa.key 2048

# @@ Optional @@
# Decode the private key
$ openssl rsa -text -in rsa.key -noout

# Extract the Public Key and save
$ openssl rsa -in rsa.key -pubout -out rsa.pub

# Create a CSR so we can sign the Certificate
$ openssl req -key rsa.key -new -out domain.csr

# @@ Optional @@
# Verify the CSR
$ openssl req -text -in domain.csr -noout -verify

# !!! TODO - Fill out the information from this command with whatever details you want !!!

# Create a Self-Signed Certificate
## Note: The -days option specifies the number of days that the certificate will be valid.
$ openssl x509 -signkey rsa.key -in domain.csr -req -days 365 -out domain.crt

# @@ Optional @@
# Verify the Certificate
$ openssl x509 -text -in domain.crt -noout

# @@ Optional @@
# Verify the Public and Private keys match (All files should share the same public key and the same hash value)
$ openssl pkey -pubout -in rsa.key | openssl sha256
$ openssl req -pubkey -in domain.csr -noout | openssl sha256
$ openssl x509 -pubkey -in domain.crt -noout | openssl sha256

# Encrypt
$ go run . pkcs7 encrypt --public-key domain.crt --string "test"
MIIBvQYJKoZIhvcNAQcDoIIBrjCCAaoCAQAxggF3MIIBcwIBADBdMEUxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQCFF/tVGl7kJuv9ogqYX57t5G9+LM2MAsGCSqGSIb3DQEBAQSCAQA5tdqDFBFJc0uc6p8Co8nqwYhnIX6s2AG+yX0Gi0FYD0KEF9UJiVMfXtkjNUh5uvpW3dA+LThb2la6d0eOAy7KbGxjh/Pujs8q3XdIMRHcUTdkTIPw8JqkvjrAHC6Sj78fT+5okWdO2Yj6p52YPtMH1soAArx4X1T1aXkwhYSWfbQ6ZROIrX7hSsCfV/q276ERw26U4wV5i6EzZt1E5yodfJXVOsbWXekfqLl5ZjjcdGb4T6muutyRTDBaFuB78XsUZfiI7cqOC1IieWad/e7/Uje9loqo+nmZ1pTqoCYfzTcHlAYoTNq80ewsBlv9RiqtIpYq2n2Q7X4wp8sg6+zHMCoGCSqGSIb3DQEHATARBgUrDgMCBwQI7EhY6z+uB4egCgQIkSNbGva47+o=

# Decrypt
$ go run . pkcs7 decrypt --private-key rsa.key --public-key domain.crt --string "MIIBvQYJKoZIhvcNAQcDoIIBrjCCAaoCAQAxggF3MIIBcwIBADBdMEUxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQCFF/tVGl7kJuv9ogqYX57t5G9+LM2MAsGCSqGSIb3DQEBAQSCAQA5tdqDFBFJc0uc6p8Co8nqwYhnIX6s2AG+yX0Gi0FYD0KEF9UJiVMfXtkjNUh5uvpW3dA+LThb2la6d0eOAy7KbGxjh/Pujs8q3XdIMRHcUTdkTIPw8JqkvjrAHC6Sj78fT+5okWdO2Yj6p52YPtMH1soAArx4X1T1aXkwhYSWfbQ6ZROIrX7hSsCfV/q276ERw26U4wV5i6EzZt1E5yodfJXVOsbWXekfqLl5ZjjcdGb4T6muutyRTDBaFuB78XsUZfiI7cqOC1IieWad/e7/Uje9loqo+nmZ1pTqoCYfzTcHlAYoTNq80ewsBlv9RiqtIpYq2n2Q7X4wp8sg6+zHMCoGCSqGSIb3DQEHATARBgUrDgMCBwQI7EhY6z+uB4egCgQIkSNbGva47+o="
test
```

### Notes

A number of these commands can also be combined, if desired:

```bash
# Create the private key and CSR
$ openssl req -newkey rsa:2048 -keyout domain.key -out domain.csr

# Create a private key and a self-signed certificate
openssl req -newkey rsa:2048 -keyout domain.key -x509 -days 365 -out domain.crt

```

If we want our private key unencrypted, we can add the `-nodes` option

```bash
$ openssl req -newkey rsa:2048 -nodes -keyout domain.key -out domain.csr
```


# TO-DO

* [ ] AES-GCM
* [ ] RC4
* [ ] JWT?

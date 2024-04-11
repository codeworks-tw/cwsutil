# Environment Variables:

| Variable              | Type   | Value                                       |
| --------------------- | ------ | ------------------------------------------- |
| CRYPTO_KEY_HEX        | string | try "openssl rand -hex 32" to generate keys |
| CRYPTO_IV_HEX         | string | try "openssl rand -hex 16" to generate keys |
| ENV                   | string | "test"/"prod"                               |
| IS_LOCAL              | bool   | "true"/"false"/"1"/"0"                      |
| DEBUG                 | bool   | "true"/"false"/"1"/"0"                      |
| LOCALIZATION_LANGUAGE | string | "en"/"zh_tw"/"zh_cn" (default: en)          |
| S3CacheTTL            | int    | (default: 10)                               |

* *S3CacheTTL*: S3 Object local cache time to live in minutes.

# Release
* 0.1.0 - Apr. 11, 2024